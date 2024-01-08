package controllers

import (
	"context"
	"fmt"

	oneclickiov1alpha1 "github.com/janlauber/one-click/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *RolloutReconciler) reconcilePVCs(ctx context.Context, rollout *oneclickiov1alpha1.Rollout) error {
	log := log.FromContext(ctx)

	log.Info("Reconciling PVCs", "Rollout", rollout.Name, "Namespace", rollout.Namespace, "Volume Count", len(rollout.Spec.Volumes))

	// Check if the Rollout spec's volume list is empty
	if len(rollout.Spec.Volumes) == 0 {
		log.Info("No volumes defined in Rollout, deleting all associated PVCs")
		// If no volumes are defined, delete all PVCs associated with this Rollout
		return r.deleteAllPVCsForRollout(ctx, rollout)
	}

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
		} else {
			// Existing PVC found, check for changes

			// Preventing name and storage class changes
			if foundPvc.Name != desiredPvc.Name || (foundPvc.Spec.StorageClassName != nil && *foundPvc.Spec.StorageClassName != *desiredPvc.Spec.StorageClassName) {
				log.Error(fmt.Errorf("name or storage class change not allowed"), "Invalid PVC update", "PVC.Name", foundPvc.Name)
				return fmt.Errorf("name or storage class change not allowed for PVC %s", foundPvc.Name)
			}

			// Handle size increase
			currentSize := foundPvc.Spec.Resources.Requests[corev1.ResourceStorage]
			desiredSize := resource.MustParse(volSpec.Size)
			if desiredSize.Cmp(currentSize) > 0 {

				// Check if PVC and its StorageClass allow resizing
				if foundPvc.Spec.VolumeMode == nil || *foundPvc.Spec.VolumeMode != corev1.PersistentVolumeFilesystem {
					log.Info("PVC resizing is only supported for filesystem volume mode")
					return nil // or handle the error as per your application's logic
				}

				storageClass := &storagev1.StorageClass{}
				if err := r.Get(ctx, types.NamespacedName{Name: *foundPvc.Spec.StorageClassName}, storageClass); err != nil {
					log.Error(err, "Failed to get the storage class of the PVC", "StorageClass", *foundPvc.Spec.StorageClassName)
					return err
				}

				if !allowsVolumeExpansion(storageClass) {
					log.Info("StorageClass does not allow volume expansion", "StorageClass", storageClass.Name)
					return nil // or handle the error as per your application's logic
				}

				// Update PVC size (considering Kubernetes limitations - PVCs can generally only be increased in size)
				foundPvc.Spec.Resources.Requests[corev1.ResourceStorage] = desiredSize
				err := r.Update(ctx, foundPvc)
				if err != nil {
					log.Error(err, "Failed to update PVC size", "PVC.Namespace", foundPvc.Namespace, "PVC.Name", foundPvc.Name)
					return err
				}
				log.Info("Updated PVC size", "PVC.Namespace", foundPvc.Namespace, "PVC.Name", foundPvc.Name)
			}
		}
	}

	// List all PVCs in the namespace
	pvcList := &corev1.PersistentVolumeClaimList{}
	if err := r.List(ctx, pvcList, client.InNamespace(rollout.Namespace)); err != nil {
		log.Error(err, "Failed to list PVCs", "Rollout.Namespace", rollout.Namespace)
		return err
	}

	for _, pvc := range pvcList.Items {
		if _, exists := expectedPVCs[pvc.Name]; !exists && isOwnedByRollout(&pvc, rollout) {
			// Delete the PVC if it's not in the expected list and is owned by the Rollout
			log.Info("Deleting PVC", "PVC.Namespace", pvc.Namespace, "PVC.Name", pvc.Name)
			if err := r.Delete(ctx, &pvc); err != nil {
				log.Error(err, "Failed to delete PVC", "PVC.Namespace", pvc.Namespace, "PVC.Name", pvc.Name)
				return err
			}
		}
	}

	return nil
}

func allowsVolumeExpansion(sc *storagev1.StorageClass) bool {
	return sc.AllowVolumeExpansion != nil && *sc.AllowVolumeExpansion
}

func (r *RolloutReconciler) constructPVCForRollout(rollout *oneclickiov1alpha1.Rollout, volSpec oneclickiov1alpha1.VolumeSpec) *corev1.PersistentVolumeClaim {
	labels := map[string]string{
		"rollout.one-click.dev/name": rollout.Name,
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

// deleteAllPVCsForRollout deletes all PVCs associated with a given Rollout
func (r *RolloutReconciler) deleteAllPVCsForRollout(ctx context.Context, rollout *oneclickiov1alpha1.Rollout) error {
	log := log.FromContext(ctx)

	pvcList := &corev1.PersistentVolumeClaimList{}
	if err := r.List(ctx, pvcList, client.InNamespace(rollout.Namespace)); err != nil {
		log.Error(err, "Failed to list PVCs for Rollout", "Rollout.Namespace", rollout.Namespace, "Rollout.Name", rollout.Name)
		return err
	}

	for _, pvc := range pvcList.Items {
		if isOwnedByRollout(&pvc, rollout) {
			if err := r.Delete(ctx, &pvc); err != nil {
				log.Error(err, "Failed to delete PVC", "PVC.Namespace", pvc.Namespace, "PVC.Name", pvc.Name)
				return err
			}
		}
	}

	return nil
}
