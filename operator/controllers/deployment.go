package controllers

import (
	"fmt"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	oneclickiov1 "github.com/janlauber/one-click/api/v1"
)

func (r *FrameworkReconciler) deploymentForFramework(f *oneclickiov1.Framework) *appsv1.Deployment {
	labels := map[string]string{"app": f.Name}
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
	if current.Spec.Template.Spec.Containers[0].Image != desiredImage {
		return true
	}

	// Check environment variables
	if !reflect.DeepEqual(current.Spec.Template.Spec.Containers[0].Env, getEnvVars(f.Spec.Env)) {
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

	return false
}

func portsMatch(currentPorts []corev1.ContainerPort, interfaces []oneclickiov1.InterfaceSpec) bool {
	if len(currentPorts) != len(interfaces) {
		return false
	}

	for i, intf := range interfaces {
		if currentPorts[i].ContainerPort != intf.Port {
			return false
		}
		// Add additional port checks if necessary
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

	// Update environment variables
	deployment.Spec.Template.Spec.Containers[0].Env = getEnvVars(f.Spec.Env)

	// Update resource requests and limits
	deployment.Spec.Template.Spec.Containers[0].Resources = createResourceRequirements(f.Spec.Resources)

	// Update ports
	updateContainerPorts(&deployment.Spec.Template.Spec.Containers[0], f.Spec.Interfaces)
}

func updateContainerPorts(container *corev1.Container, interfaces []oneclickiov1.InterfaceSpec) {
	var ports []corev1.ContainerPort
	for _, intf := range interfaces {
		ports = append(ports, corev1.ContainerPort{ContainerPort: intf.Port})
		// Add additional port configuration if necessary
	}
	container.Ports = ports
}
