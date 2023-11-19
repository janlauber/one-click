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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	oneclickiov1alpha1 "github.com/janlauber/one-click/api/v1alpha1"
)

// RolloutReconciler reconciles a Rollout object
type RolloutReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=one-click.io,resources=frameworks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=one-click.io,resources=frameworks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=one-click.io,resources=frameworks/finalizers,verbs=update
//+kubebuilder:rbac:groups=one-click.io,resources=hpa,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=one-click.io,resources=hpa/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=one-click.io,resources=hpa/finalizers,verbs=update
//+kubebuilder:rbac:groups=one-click.io,resources=service,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=one-click.io,resources=service/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=one-click.io,resources=service/finalizers,verbs=update
//+kubebuilder:rbac:groups=one-click.io,resources=ingress,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=one-click.io,resources=ingress/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=one-click.io,resources=ingress/finalizers,verbs=update
//+kubebuilder:rbac:groups=one-click.io,resources=secret,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=one-click.io,resources=secret/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=one-click.io,resources=secret/finalizers,verbs=update
//+kubebuilder:rbac:groups=one-click.io,resources=deployment,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=one-click.io,resources=deployment/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=one-click.io,resources=deployment/finalizers,verbs=update
//+kubebuilder:rbac:groups=one-click.io,resources=serviceaccount,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=one-click.io,resources=serviceaccount/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=one-click.io,resources=serviceaccount/finalizers,verbs=update
//+kubebuilder:rbac:groups=one-click.io,resources=persistentvolumeclaim,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=one-click.io,resources=persistentvolumeclaim/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=one-click.io,resources=persistentvolumeclaim/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *RolloutReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the Rollout instance
	var rollout oneclickiov1alpha1.Rollout
	if err := r.Get(ctx, req.NamespacedName, &rollout); err != nil {
		if errors.IsNotFound(err) {
			// Object not found, could have been deleted after reconcile request, return and don't requeue
			log.Info("Rollout resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Rollout.")
		return ctrl.Result{}, err
	}

	// Reconcile ServiceAccount
	if err := r.reconcileServiceAccount(ctx, &rollout); err != nil {
		log.Error(err, "Failed to reconcile ServiceAccount.")
		return ctrl.Result{}, err
	}

	// Reconcile PVCs only if volumes are defined
	if len(rollout.Spec.Volumes) > 0 {
		if err := r.reconcilePVCs(ctx, &rollout); err != nil {
			log.Error(err, "Failed to reconcile PVCs.")
			return ctrl.Result{}, err
		}
	}

	// Reconcile Secrets
	if err := r.reconcileSecret(ctx, &rollout); err != nil {
		log.Error(err, "Failed to reconcile Secrets.")
		return ctrl.Result{}, err
	}

	// Reconcile Deployment
	if err := r.reconcileDeployment(ctx, &rollout); err != nil {
		log.Error(err, "Failed to reconcile Deployment.")
		return ctrl.Result{}, err
	}

	// Reconcile Service
	if err := r.reconcileService(ctx, &rollout); err != nil {
		log.Error(err, "Failed to reconcile Service.")
		return ctrl.Result{}, err
	}

	// Reconcile Ingress
	if err := r.reconcileIngress(ctx, &rollout); err != nil {
		log.Error(err, "Failed to reconcile Ingress.")
		return ctrl.Result{}, err
	}

	// Reconcile HPA
	if err := r.reconcileHPA(ctx, &rollout); err != nil {
		log.Error(err, "Failed to reconcile HPA.")
		return ctrl.Result{}, err
	}

	// Update status
	if err := r.updateStatus(ctx, &rollout); err != nil {
		if errors.IsConflict(err) {
			log.Info("Conflict while updating status. Retrying.")
			return ctrl.Result{Requeue: true}, nil
		}
		log.Error(err, "Failed to update status.")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RolloutReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&oneclickiov1alpha1.Rollout{}).
		Complete(r)
}