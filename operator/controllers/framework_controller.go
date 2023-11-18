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
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

	// Handling the aggregated Secret
	secret := r.secretForFramework(&framework)
	foundSecret := &corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace}, foundSecret)
	if err != nil && errors.IsNotFound(err) {
		log.Log.Info("Creating a new Secret", "Namespace", secret.Namespace, "Name", secret.Name)
		err = r.Create(ctx, secret)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	} else if !reflect.DeepEqual(foundSecret.StringData, secret.StringData) {
		// Update the Secret if it already exists and the data has changed
		foundSecret.StringData = secret.StringData
		err = r.Update(ctx, foundSecret)
		if err != nil {
			return ctrl.Result{}, err
		}
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

	// Define the desired HPA object
	desiredHpa := r.hpaForFramework(&framework)

	// Check if this HPA already exists
	foundHpa := &autoscalingv2.HorizontalPodAutoscaler{}
	err = r.Get(ctx, types.NamespacedName{Name: framework.Name, Namespace: framework.Namespace}, foundHpa)
	if err != nil {
		if errors.IsNotFound(err) {
			// HPA not found - create it
			log.Log.Info("Creating a new HorizontalPodAutoscaler", "Namespace", framework.Namespace, "Name", framework.Name)
			err = r.Create(ctx, desiredHpa)
			if err != nil {
				// handle error
				return ctrl.Result{}, err
			}
			// HPA created successfully - return and requeue

			// After creating the HPA object
			if err := ctrl.SetControllerReference(&framework, desiredHpa, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		} else {
			// handle error
			return ctrl.Result{}, err
		}
	} else {
		// HPA exists - check if it needs an update
		if needsHpaUpdate(foundHpa, &framework) {
			// Update the found HPA object to match the desired state
			updateHpa(foundHpa, &framework)

			err = r.Update(ctx, foundHpa)
			if err != nil {
				// handle error
				return ctrl.Result{}, err
			}
		}
	}

	// Reconcile the Service and Ingress for each Interface
	for _, intf := range framework.Spec.Interfaces {
		// Handle Service for each interface
		err := r.reconcileService(ctx, &framework, intf)
		if err != nil {
			return ctrl.Result{}, err
		}

		// Handle Ingress for each interface (if defined)
		if len(intf.Ingress.Rules) > 0 {
			err := r.reconcileIngress(ctx, &framework, intf)
			if err != nil {
				return ctrl.Result{}, err
			}
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
