package controllers

import (
	"context"

	oneclickiov1alpha1 "github.com/janlauber/one-click/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// TODO: fix PVCs not being deleted when removed from Rollout spec
// TODO: crashes when volume size is increased

func (r *RolloutReconciler) reconcilePVCs(ctx context.Context, rollout *oneclickiov1alpha1.Rollout) error {
	log := log.FromContext(ctx)

	// Keep track of the PVCs that should exist according to the Rollout specification
	expectedPVCs := make(map[string]struct{})
	for _, volSpec := range rollout.Spec.Volumes {
		expectedPVCs[volSpec.Name] = struct{}{}
		desiredPvc := r.constructPVCForRollout(rollout, volSpec)

		foundPvc := &corev1.PersistentVolumeClaim{}
		err := r.Get(ctx, types.NamespacedName{Name: desiredPvc.Name, Namespace: rollout.Namespace}, foundPvc)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating a new PVC", "PVC.Namespace", rollout.Namespace, "PVC.Name", desiredPvc.Name)
			err = r.Create(ctx, desiredPvc)
			if err != nil {
				r.Recorder.Eventf(rollout, corev1.EventTypeWarning, "CreationFailed", "Failed to create PVC %s", desiredPvc.Name)
				log.Error(err, "Failed to create new PVC", "PVC.Namespace", rollout.Namespace, "PVC.Name", desiredPvc.Name)
				return err
			}
			r.Recorder.Eventf(rollout, corev1.EventTypeNormal, "Created", "Created PVC %s", desiredPvc.Name)
		} else if err != nil {
			r.Recorder.Eventf(rollout, corev1.EventTypeWarning, "GetFailed", "Failed to get PVC %s", desiredPvc.Name)
			log.Error(err, "Failed to get PVC")
			return err
		}
		// Check if PVC needs to be updated
		if err == nil {
			sizeChangeNeeded := foundPvc.Spec.Resources.Requests[corev1.ResourceStorage] != resource.MustParse(volSpec.Size)
			// Include other checks if necessary (e.g., storage class change)

			if sizeChangeNeeded {
				// Update PVC size (considering Kubernetes limitations - PVCs can generally only be increased in size)
				foundPvc.Spec.Resources.Requests[corev1.ResourceStorage] = resource.MustParse(volSpec.Size)
				if err := r.Update(ctx, foundPvc); err != nil {
					log.Error(err, "Failed to update PVC", "PVC.Namespace", foundPvc.Namespace, "PVC.Name", foundPvc.Name)
					return err
				}
			}
		}
	}

	// Clean up PVCs that should no longer exist
	pvcList := &corev1.PersistentVolumeClaimList{}
	if err := r.List(ctx, pvcList, client.InNamespace(rollout.Namespace)); err != nil {
		log.Error(err, "Failed to list PVCs")
		return err
	}

	for _, pvc := range pvcList.Items {
		if _, exists := expectedPVCs[pvc.Name]; !exists && isOwnedByRollout(&pvc, rollout) {
			// PVC is owned by Rollout but no longer in spec, delete it
			if err := r.Delete(ctx, &pvc); err != nil {
				log.Error(err, "Failed to delete PVC", "PVC.Namespace", pvc.Namespace, "PVC.Name", pvc.Name)
				return err
			}
		}
	}

	return nil
}

func (r *RolloutReconciler) constructPVCForRollout(rollout *oneclickiov1alpha1.Rollout, volSpec oneclickiov1alpha1.VolumeSpec) *corev1.PersistentVolumeClaim {
	labels := map[string]string{
		"rollout.one-click.io/name": rollout.Name,
	}

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      volSpec.Name,
			Namespace: rollout.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(volSpec.Size),
				},
			},
			StorageClassName: &volSpec.StorageClass,
		},
	}

	ctrl.SetControllerReference(rollout, pvc, r.Scheme)
	return pvc
}

func isOwnedByRollout(pvc *corev1.PersistentVolumeClaim, rollout *oneclickiov1alpha1.Rollout) bool {
	for _, ref := range pvc.GetOwnerReferences() {
		if ref.Kind == "Rollout" && ref.Name == rollout.Name {
			return true
		}
	}
	return false
}
