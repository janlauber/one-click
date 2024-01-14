package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	oneclickiov1alpha1 "github.com/janlauber/one-click/api/v1alpha1"
)

func (r *RolloutReconciler) reconcileDeployment(ctx context.Context, rollout *oneclickiov1alpha1.Rollout) error {
	log := log.FromContext(ctx)

	desiredDeployment := r.deploymentForRollout(ctx, rollout)

	currentDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: rollout.Name, Namespace: rollout.Namespace}, currentDeployment)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Deployment", "Deployment.Namespace", rollout.Namespace, "Deployment.Name", rollout.Name)
		return r.Create(ctx, desiredDeployment)
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return err
	}

	// Compare the current Deployment with the Rollout spec
	if needsUpdate(currentDeployment, rollout) {
		// Update the Deployment to align it with the Rollout spec
		currentDeployment.Spec = desiredDeployment.Spec
		err = r.Update(ctx, currentDeployment)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", currentDeployment.Namespace, "Deployment.Name", currentDeployment.Name)
			return err
		}
	}

	return nil
}

func (r *RolloutReconciler) deploymentForRollout(ctx context.Context, f *oneclickiov1alpha1.Rollout) *appsv1.Deployment {
	log := log.FromContext(context.Background())
	labels := map[string]string{"rollout.one-click.dev/name": f.Name}
	replicas := int32(f.Spec.HorizontalScale.MinReplicas)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      f.Name,
			Namespace: f.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  f.Name,
						Image: fmt.Sprintf("%s/%s:%s", f.Spec.Image.Registry, f.Spec.Image.Repository, f.Spec.Image.Tag),
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse(f.Spec.Resources.Requests.CPU),
								corev1.ResourceMemory: resource.MustParse(f.Spec.Resources.Requests.Memory),
							},
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse(f.Spec.Resources.Limits.CPU),
								corev1.ResourceMemory: resource.MustParse(f.Spec.Resources.Limits.Memory),
							},
						},
						Ports: getContainerPorts(f.Spec.Interfaces),
						Env:   getEnvVars(f.Spec.Env),
					}},
					ServiceAccountName: f.Spec.ServiceAccountName,
				},
			},
		},
	}

	// if security context is defined, add it to the pod security context
	if !reflect.DeepEqual(f.Spec.SecurityContext, oneclickiov1alpha1.SecurityContextSpec{}) {
		dep.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{
			FSGroup:    &f.Spec.SecurityContext.FsGroup,
			RunAsUser:  &f.Spec.SecurityContext.RunAsUser,
			RunAsGroup: &f.Spec.SecurityContext.RunAsGroup,
		}

		capabilities := make([]corev1.Capability, len(f.Spec.SecurityContext.Capabilities.Add))
		for i, cap := range f.Spec.SecurityContext.Capabilities.Add {
			capabilities[i] = corev1.Capability(cap)
		}

		dropCapabilities := make([]corev1.Capability, len(f.Spec.SecurityContext.Capabilities.Drop))
		for i, cap := range f.Spec.SecurityContext.Capabilities.Drop {
			dropCapabilities[i] = corev1.Capability(cap)
		}

		dep.Spec.Template.Spec.Containers[0].SecurityContext = &corev1.SecurityContext{
			Privileged:               &f.Spec.SecurityContext.Privileged,
			ReadOnlyRootFilesystem:   &f.Spec.SecurityContext.ReadOnlyRootFilesystem,
			RunAsNonRoot:             &f.Spec.SecurityContext.RunAsNonRoot,
			AllowPrivilegeEscalation: &f.Spec.SecurityContext.AllowPrivilegeEscalation,
			Capabilities: &corev1.Capabilities{
				Add:  capabilities,
				Drop: dropCapabilities,
			},
		}
	}

	// if secrets are defined, add the secret f.Name + "-secrets" as envFrom
	if len(f.Spec.Secrets) > 0 {
		dep.Spec.Template.Spec.Containers[0].EnvFrom = []corev1.EnvFromSource{
			{
				SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: f.Name + "-secrets",
					},
				},
			},
		}
	}

	// Check if Username and Password are provided
	if f.Spec.Image.Username != "" && f.Spec.Image.Password != "" {
		// Logic to create or get existing secret
		secretName := f.Name + "-imagepullsecret"
		if err := r.reconcileImagePullSecret(ctx, f, secretName); err != nil {
			log.Error(err, "Failed to reconcile Image Pull Secret")
			// Handle error, possibly return it
		}

		// Attach the image pull secret to the deployment
		dep.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{
			{
				Name: secretName,
			},
		}
	} else {
		// Remove image pull secret from deployment
		dep.Spec.Template.Spec.ImagePullSecrets = nil
		// Logic to create or get existing secret
		secretName := f.Name + "-imagepullsecret"
		if err := r.reconcileImagePullSecret(ctx, f, secretName); err != nil {
			log.Error(err, "Failed to reconcile Image Pull Secret")
			// Handle error, possibly return it
		}
	}

	// Update volumes and volume mounts
	if len(f.Spec.Volumes) > 0 {
		var volumes []corev1.Volume
		var volumeMounts []corev1.VolumeMount
		for _, v := range f.Spec.Volumes {
			volumes = append(volumes, corev1.Volume{
				Name: v.Name,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: v.Name,
					},
				},
			})
			volumeMounts = append(volumeMounts, corev1.VolumeMount{
				Name:      v.Name,
				MountPath: v.MountPath,
			})
		}
		dep.Spec.Template.Spec.Volumes = volumes
		dep.Spec.Template.Spec.Containers[0].VolumeMounts = volumeMounts
	} else {
		// Handle no volumes case
		dep.Spec.Template.Spec.Volumes = nil
		dep.Spec.Template.Spec.Containers[0].VolumeMounts = nil
	}

	// Add secret checksum to pod template annotations to trigger redeployment when secrets change
	checksum := calculateSecretsChecksum(f.Spec.Secrets)
	if dep.Spec.Template.Annotations == nil {
		dep.Spec.Template.Annotations = make(map[string]string)
	}
	dep.Spec.Template.Annotations["one-click.dev/secrets-checksum"] = checksum

	ctrl.SetControllerReference(f, dep, r.Scheme)
	return dep
}

func getContainerPorts(interfaces []oneclickiov1alpha1.InterfaceSpec) []corev1.ContainerPort {
	var ports []corev1.ContainerPort
	for _, i := range interfaces {
		ports = append(ports, corev1.ContainerPort{
			ContainerPort: i.Port,
		})
	}
	return ports
}

func getEnvVars(envVars []oneclickiov1alpha1.EnvVar) []corev1.EnvVar {
	var envs []corev1.EnvVar
	for _, env := range envVars {
		envs = append(envs, corev1.EnvVar{
			Name:  env.Name,
			Value: env.Value,
		})
	}
	return envs
}

func needsUpdate(current *appsv1.Deployment, f *oneclickiov1alpha1.Rollout) bool {
	// Check replicas
	if *current.Spec.Replicas != int32(f.Spec.HorizontalScale.MinReplicas) {
		return true
	}

	// Check security context
	if !reflect.DeepEqual(current.Spec.Template.Spec.SecurityContext, &corev1.PodSecurityContext{
		FSGroup:    &f.Spec.SecurityContext.FsGroup,
		RunAsUser:  &f.Spec.SecurityContext.RunAsUser,
		RunAsGroup: &f.Spec.SecurityContext.RunAsGroup,
	}) {
		return true
	}

	// Check container security context
	if !reflect.DeepEqual(current.Spec.Template.Spec.Containers[0].SecurityContext, &corev1.SecurityContext{
		Privileged:               &f.Spec.SecurityContext.Privileged,
		ReadOnlyRootFilesystem:   &f.Spec.SecurityContext.ReadOnlyRootFilesystem,
		RunAsNonRoot:             &f.Spec.SecurityContext.RunAsNonRoot,
		AllowPrivilegeEscalation: &f.Spec.SecurityContext.AllowPrivilegeEscalation,
		Capabilities: &corev1.Capabilities{
			Add:  []corev1.Capability{},
			Drop: []corev1.Capability{},
		},
	}) {
		return true
	}

	// Check container image
	desiredImage := fmt.Sprintf("%s/%s:%s", f.Spec.Image.Registry, f.Spec.Image.Repository, f.Spec.Image.Tag)
	if len(current.Spec.Template.Spec.Containers) == 0 || current.Spec.Template.Spec.Containers[0].Image != desiredImage {
		return true
	}

	// Check secrets
	if len(f.Spec.Secrets) > 0 {
		if len(current.Spec.Template.Spec.Containers[0].EnvFrom) == 0 {
			return true
		}

		if current.Spec.Template.Spec.Containers[0].EnvFrom[0].SecretRef.LocalObjectReference.Name != f.Name+"-secrets" {
			return true
		}
	} else {
		if len(current.Spec.Template.Spec.Containers[0].EnvFrom) > 0 {
			return true
		}
	}

	// Check environment variables
	desiredEnvVars := getEnvVars(f.Spec.Env)
	if !reflect.DeepEqual(current.Spec.Template.Spec.Containers[0].Env, desiredEnvVars) {
		return true
	}

	// Check resource requests and limits
	desiredResources := createResourceRequirements(f.Spec.Resources)
	if !reflect.DeepEqual(current.Spec.Template.Spec.Containers[0].Resources, desiredResources) {
		return true
	}

	// Check ports
	if !portsMatch(current.Spec.Template.Spec.Containers[0].Ports, f.Spec.Interfaces) {
		return true
	}

	// Check volumes
	if !volumesMatch(current.Spec.Template.Spec.Volumes, f.Spec.Volumes) {
		return true
	}

	// Check service account name
	if current.Spec.Template.Spec.ServiceAccountName != f.Spec.ServiceAccountName {
		return true
	}

	// Add more checks as necessary, e.g., labels, annotations, specific configuration, etc.

	return false
}

func volumesMatch(currentVolumes []corev1.Volume, desiredVolumes []oneclickiov1alpha1.VolumeSpec) bool {
	if len(currentVolumes) != len(desiredVolumes) {
		return false
	}

	desiredVolumeMap := make(map[string]oneclickiov1alpha1.VolumeSpec)
	for _, v := range desiredVolumes {
		desiredVolumeMap[v.Name] = v
	}

	for _, currentVolume := range currentVolumes {
		volSpec, exists := desiredVolumeMap[currentVolume.Name]
		if !exists {
			// Volume is present in Deployment but not in Rollout spec
			return false
		}

		// Check PVC name
		if currentVolume.VolumeSource.PersistentVolumeClaim.ClaimName != volSpec.Name {
			return false
		}
		// Additional checks can be added here, such as PVC size, storage class, etc.
	}

	return true
}

func portsMatch(currentPorts []corev1.ContainerPort, interfaces []oneclickiov1alpha1.InterfaceSpec) bool {
	if len(currentPorts) != len(interfaces) {
		return false
	}

	for i, intf := range interfaces {
		if currentPorts[i].ContainerPort != intf.Port {
			return false
		}

		// check for name
		if currentPorts[i].Name != intf.Name {
			return false
		}
	}

	return true
}

func createResourceRequirements(resources oneclickiov1alpha1.ResourceRequirements) corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(resources.Requests.CPU),
			corev1.ResourceMemory: resource.MustParse(resources.Requests.Memory),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(resources.Limits.CPU),
			corev1.ResourceMemory: resource.MustParse(resources.Limits.Memory),
		},
	}
}

type dockerConfigEntry struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Auth     string `json:"auth"`
}

type dockerConfigJSON struct {
	Auths map[string]dockerConfigEntry `json:"auths"`
}

func (r *RolloutReconciler) reconcileImagePullSecret(ctx context.Context, rollout *oneclickiov1alpha1.Rollout, secretName string) error {
	registry := rollout.Spec.Image.Registry
	username := rollout.Spec.Image.Username
	password := rollout.Spec.Image.Password

	// Check if username and password are empty, delete the secret if it exists
	if username == "" && password == "" {
		foundSecret := &corev1.Secret{}
		err := r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: rollout.Namespace}, foundSecret)
		if err == nil {
			if metav1.IsControlledBy(foundSecret, rollout) {
				return r.Delete(ctx, foundSecret)
			}
		}
		return nil
	}

	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	dockerConfigEntryData := dockerConfigEntry{
		Username: username,
		Password: password,
		Auth:     auth,
	}
	dockerConfigJSON := dockerConfigJSON{
		Auths: map[string]dockerConfigEntry{
			registry: dockerConfigEntryData,
		},
	}

	dockerConfigJSONBytes, err := json.Marshal(dockerConfigJSON)
	if err != nil {
		return err
	}

	secretData := map[string][]byte{
		".dockerconfigjson": dockerConfigJSONBytes,
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: rollout.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(rollout, rollout.GroupVersionKind()),
			},
		},
		Type: corev1.SecretTypeDockerConfigJson,
		Data: secretData,
	}

	// Check if secret already exists, if not, create it
	foundSecret := &corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: rollout.Namespace}, foundSecret)
	if err != nil && errors.IsNotFound(err) {
		// Secret not found, create it
		return r.Create(ctx, secret)
	} else if err != nil {
		// Error occurred while checking for secret
		return err
	}

	// Secret found, update it
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		foundSecret.Data = secretData
		return r.Update(ctx, foundSecret)
	})
}
