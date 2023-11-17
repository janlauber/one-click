/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	oneclickiov1 "github.com/janlauber/one-click/api/v1"
)

// FrameworkReconciler reconciles a Framework object
type FrameworkReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=one-click.io,resources=frameworks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=one-click.io,resources=frameworks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=one-click.io,resources=frameworks/finalizers,verbs=update

func (r *FrameworkReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Instance to hold the fetched Framework object
	var framework oneclickiov1.Framework

	// Fetch the Framework instance using the namespaced name
	if err := r.Get(ctx, req.NamespacedName, &framework); err != nil {
		log.Log.Error(err, "unable to fetch Framework")

		// Handle the case where the Framework resource no longer exists
		// Remove all resources associated with the Framework

		// If it's not a not-found error, requeue the request
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Define the ServiceAccount you expect to exist
	expectedSa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      framework.Spec.ServiceAccountName,
			Namespace: req.Namespace,
		},
	}

	// Set Framework instance as the owner of the ServiceAccount
	if err := ctrl.SetControllerReference(&framework, expectedSa, r.Scheme); err != nil {
		// handle error
		return ctrl.Result{}, err
	}

	// Try to get the ServiceAccount
	foundSa := &corev1.ServiceAccount{}
	err := r.Get(ctx, types.NamespacedName{Name: framework.Spec.ServiceAccountName, Namespace: req.Namespace}, foundSa)
	if err != nil && errors.IsNotFound(err) {
		// ServiceAccount doesn't exist - create it
		log.Log.Info("Creating a new ServiceAccount", "Namespace", req.Namespace, "Name", framework.Spec.ServiceAccountName)
		err = r.Create(ctx, expectedSa)
		if err != nil {
			// handle error
			return ctrl.Result{}, err
		}
		// ServiceAccount created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		// handle error
		return ctrl.Result{}, err
	}

	// Check if the Deployment already exists
	deployment := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: framework.Name, Namespace: framework.Namespace}, deployment)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Deployment
		dep := r.deploymentForFramework(&framework)
		log.Log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	} else {
		// Deployment exists - check if it needs an update
		if needsUpdate(deployment, &framework) {
			// Update the Deployment to match the desired state
			updateDeployment(deployment, &framework)
			err = r.Update(ctx, deployment)
			if err != nil {
				// Handle error
				return ctrl.Result{}, err
			}
			log.Log.Info("Updated Deployment", "Namespace", deployment.Namespace, "Name", deployment.Name)
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FrameworkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&oneclickiov1.Framework{}).
		Complete(r)
}

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
