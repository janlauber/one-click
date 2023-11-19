package controllers

import (
	"context"
	"fmt"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	oneclickiov1 "github.com/janlauber/one-click/api/v1"
)

func (r *FrameworkReconciler) reconcileDeployment(ctx context.Context, f *oneclickiov1.Framework) error {
	log := log.FromContext(ctx)

	// Get the desired state of the Deployment from the helper function
	desiredDeployment := r.deploymentForFramework(f)

	// Check if the Deployment already exists
	foundDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: f.Name, Namespace: f.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Deployment", "Deployment.Namespace", f.Namespace, "Deployment.Name", f.Name)
		err = r.Create(ctx, desiredDeployment)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", f.Namespace, "Deployment.Name", f.Name)
			r.Recorder.Eventf(f, corev1.EventTypeWarning, "CreationFailed", "Failed to create Deployment %s", f.Name)
			return err
		}
		r.Recorder.Eventf(f, corev1.EventTypeNormal, "Created", "Created Deployment %s", f.Name)
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		r.Recorder.Eventf(f, corev1.EventTypeWarning, "GetFailed", "Failed to get Deployment %s", f.Name)
		return err
	} else {
		// Deployment exists - check if it needs an update
		if needsUpdate(foundDeployment, f) {
			log.Info("Updating Deployment", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name)
			updateDeployment(foundDeployment, f)
			err = r.Update(ctx, foundDeployment)
			if err != nil {
				log.Error(err, "Failed to update Deployment", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name)
				r.Recorder.Eventf(f, corev1.EventTypeWarning, "UpdateFailed", "Failed to update Deployment %s", foundDeployment.Name)
				return err
			}
			r.Recorder.Eventf(f, corev1.EventTypeNormal, "Updated", "Updated Deployment %s", foundDeployment.Name)
		}
	}

	return nil
}

func (r *FrameworkReconciler) deploymentForFramework(f *oneclickiov1.Framework) *appsv1.Deployment {
	labels := map[string]string{"framework.one-click.io/name": f.Name}
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

	// if volumes are defined, add them to the pod
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
	}

	ctrl.SetControllerReference(f, dep, r.Scheme)
	return dep
}

func getContainerPorts(interfaces []oneclickiov1.InterfaceSpec) []corev1.ContainerPort {
	var ports []corev1.ContainerPort
	for _, i := range interfaces {
		ports = append(ports, corev1.ContainerPort{
			ContainerPort: i.Port,
		})
	}
	return ports
}

func getEnvVars(envVars []oneclickiov1.EnvVar) []corev1.EnvVar {
	var envs []corev1.EnvVar
	for _, env := range envVars {
		envs = append(envs, corev1.EnvVar{
			Name:  env.Name,
			Value: env.Value,
		})
	}
	return envs
}

func needsUpdate(current *appsv1.Deployment, f *oneclickiov1.Framework) bool {
	// Check replicas
	if *current.Spec.Replicas != int32(f.Spec.HorizontalScale.MinReplicas) {
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

func volumesMatch(currentVolumes []corev1.Volume, desiredVolumes []oneclickiov1.VolumeSpec) bool {
	if len(currentVolumes) != len(desiredVolumes) {
		return false
	}

	desiredVolumeMap := make(map[string]oneclickiov1.VolumeSpec)
	for _, v := range desiredVolumes {
		desiredVolumeMap[v.Name] = v
	}

	for _, currentVolume := range currentVolumes {
		volSpec, exists := desiredVolumeMap[currentVolume.Name]
		if !exists {
			// Volume is present in Deployment but not in Framework spec
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

func portsMatch(currentPorts []corev1.ContainerPort, interfaces []oneclickiov1.InterfaceSpec) bool {
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

func createResourceRequirements(resources oneclickiov1.ResourceRequirements) corev1.ResourceRequirements {
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

func updateDeployment(deployment *appsv1.Deployment, f *oneclickiov1.Framework) {
	// Update replicas
	deployment.Spec.Replicas = &f.Spec.HorizontalScale.MinReplicas

	// Update container image
	deployment.Spec.Template.Spec.Containers[0].Image = fmt.Sprintf("%s/%s:%s", f.Spec.Image.Registry, f.Spec.Image.Repository, f.Spec.Image.Tag)

	// Update secrets
	if len(f.Spec.Secrets) > 0 {
		deployment.Spec.Template.Spec.Containers[0].EnvFrom = []corev1.EnvFromSource{
			{
				SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: f.Name + "-secrets",
					},
				},
			},
		}
	} else {
		deployment.Spec.Template.Spec.Containers[0].EnvFrom = nil
	}

	// Ensure the Annotations map is not nil
	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}

	// Calculate the checksum of the secrets and add it as an annotation
	secretsChecksum := calculateSecretsChecksum(f.Spec.Secrets)
	deployment.Spec.Template.Annotations["one-click.io/secrets-checksum"] = secretsChecksum

	// Update environment variables
	deployment.Spec.Template.Spec.Containers[0].Env = getEnvVars(f.Spec.Env)

	// Update resource requests and limits
	deployment.Spec.Template.Spec.Containers[0].Resources = createResourceRequirements(f.Spec.Resources)

	// Update ports
	updateContainerPorts(&deployment.Spec.Template.Spec.Containers[0], f.Spec.Interfaces)

	// Update volumes
	updateVolumes(&deployment.Spec.Template.Spec.Volumes, f.Spec.Volumes)

	// Update volume mounts
	var volumeMounts []corev1.VolumeMount
	for _, v := range f.Spec.Volumes {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      v.Name,
			MountPath: v.MountPath,
		})
	}
	deployment.Spec.Template.Spec.Containers[0].VolumeMounts = volumeMounts

	// Update service account name
	deployment.Spec.Template.Spec.ServiceAccountName = f.Spec.ServiceAccountName
}

func updateContainerPorts(container *corev1.Container, interfaces []oneclickiov1.InterfaceSpec) {
	var ports []corev1.ContainerPort
	for _, intf := range interfaces {
		ports = append(ports, corev1.ContainerPort{ContainerPort: intf.Port})
		// Add additional port configuration if necessary
	}
	container.Ports = ports
}

func updateVolumes(currentVolumes *[]corev1.Volume, volumes []oneclickiov1.VolumeSpec) {
	var newVolumes []corev1.Volume
	for _, vol := range volumes {
		newVolumes = append(newVolumes, corev1.Volume{
			Name: vol.Name,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: vol.Name,
				},
			},
		})
		// Add additional volume configuration if necessary
	}
	*currentVolumes = newVolumes
}
