package controllers

import (
	"context"

	oneclickiov1 "github.com/janlauber/one-click/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *FrameworkReconciler) reconcileServiceAccount(ctx context.Context, framework *oneclickiov1.Framework) error {
	log := log.FromContext(ctx)

	// Define the ServiceAccount you expect to exist
	expectedSa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      framework.Spec.ServiceAccountName,
			Namespace: framework.Namespace,
		},
	}

	// Set Framework instance as the owner of the ServiceAccount
	if err := ctrl.SetControllerReference(framework, expectedSa, r.Scheme); err != nil {
		log.Error(err, "Failed to set controller reference for ServiceAccount", "ServiceAccount.Namespace", framework.Namespace, "ServiceAccount.Name", framework.Spec.ServiceAccountName)
		return err
	}

	// Try to get the ServiceAccount
	foundSa := &corev1.ServiceAccount{}
	err := r.Get(ctx, types.NamespacedName{Name: framework.Spec.ServiceAccountName, Namespace: framework.Namespace}, foundSa)
	if err != nil && errors.IsNotFound(err) {
		// ServiceAccount doesn't exist - create it
		err = r.Create(ctx, expectedSa)
		if err != nil {
			// Error occurred while trying to create the ServiceAccount
			r.Recorder.Eventf(framework, corev1.EventTypeWarning, "CreationFailed", "Failed to create ServiceAccount %s", framework.Spec.ServiceAccountName)
			log.Error(err, "Failed to create new ServiceAccount", "ServiceAccount.Namespace", framework.Namespace, "ServiceAccount.Name", framework.Spec.ServiceAccountName)
			return err
		}
		r.Recorder.Eventf(framework, corev1.EventTypeNormal, "Created", "Created ServiceAccount %s", framework.Spec.ServiceAccountName)
		// ServiceAccount created successfully
	} else if err != nil {
		// Error occurred while trying to get the ServiceAccount
		r.Recorder.Eventf(framework, corev1.EventTypeWarning, "GetFailed", "Failed to get ServiceAccount %s", framework.Spec.ServiceAccountName)
		log.Error(err, "Failed to get ServiceAccount")
		return err
	}
	// ServiceAccount already exists - no action required

	return nil
}
