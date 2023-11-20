package controllers

import (
	"context"

	oneclickiov1alpha1 "github.com/janlauber/one-click/api/v1alpha1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *RolloutReconciler) reconcileHPA(ctx context.Context, f *oneclickiov1alpha1.Rollout) error {
	log := log.FromContext(ctx)

	// Construct the desired HPA object based on the Rollout specification
	desiredHpa, err := r.hpaForRollout(f)
	if err != nil {
		log.Error(err, "Failed to construct HPA", "Namespace", f.Namespace, "Name", f.Name)
		r.Recorder.Eventf(f, corev1.EventTypeWarning, "CreationFailed", "Failed to construct HPA %s", f.Name)
		return err
	}

	// Try to fetch the existing HPA
	foundHpa := &autoscalingv2.HorizontalPodAutoscaler{}
	err = r.Get(ctx, types.NamespacedName{Name: f.Name, Namespace: f.Namespace}, foundHpa)
	if err != nil && errors.IsNotFound(err) {
		// If the HPA is not found, create a new one
		log.Info("Creating a new HPA", "Namespace", desiredHpa.Namespace, "Name", desiredHpa.Name)
		err = r.Create(ctx, desiredHpa)
		if err != nil {
			// Handle creation error
			log.Error(err, "Failed to create HPA", "Namespace", desiredHpa.Namespace, "Name", desiredHpa.Name)
			r.Recorder.Eventf(f, corev1.EventTypeWarning, "CreationFailed", "Failed to create HPA %s", f.Name)
			return err
		}
		r.Recorder.Eventf(f, corev1.EventTypeNormal, "Created", "Created HPA %s", f.Name)
	} else if err != nil {
		// Handle other errors
		log.Error(err, "Failed to get HPA", "Namespace", desiredHpa.Namespace, "Name", desiredHpa.Name)
		r.Recorder.Eventf(f, corev1.EventTypeWarning, "GetFailed", "Failed to get HPA %s", f.Name)
		return err
	} else {
		// If the HPA exists, check if it needs to be updated
		if needsHpaUpdate(foundHpa, f) {
			log.Info("Updating HPA", "Namespace", foundHpa.Namespace, "Name", foundHpa.Name)
			updateHpa(foundHpa, f)
			err = r.Update(ctx, foundHpa)
			if err != nil {
				// Handle update error
				log.Error(err, "Failed to update HPA", "Namespace", foundHpa.Namespace, "Name", foundHpa.Name)
				r.Recorder.Eventf(f, corev1.EventTypeWarning, "UpdateFailed", "Failed to update HPA %s", foundHpa.Name)
				return err
			}
			r.Recorder.Eventf(f, corev1.EventTypeNormal, "Updated", "Updated HPA %s", foundHpa.Name)
		}
	}

	return nil
}

func (r *RolloutReconciler) hpaForRollout(f *oneclickiov1alpha1.Rollout) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      f.Name,
			Namespace: f.Namespace,
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       f.Name,
			},
			MinReplicas: &f.Spec.HorizontalScale.MinReplicas,
			MaxReplicas: f.Spec.HorizontalScale.MaxReplicas,
			Metrics: []autoscalingv2.MetricSpec{
				{
					Type: autoscalingv2.ResourceMetricSourceType,
					Resource: &autoscalingv2.ResourceMetricSource{
						Name:   corev1.ResourceCPU,
						Target: autoscalingv2.MetricTarget{Type: autoscalingv2.UtilizationMetricType, AverageUtilization: &f.Spec.HorizontalScale.TargetCPUUtilizationPercentage},
					},
				},
			},
		},
	}

	// Set the owner reference
	if err := controllerutil.SetControllerReference(f, hpa, r.Scheme); err != nil {
		return nil, err
	}

	return hpa, nil
}

func needsHpaUpdate(current *autoscalingv2.HorizontalPodAutoscaler, f *oneclickiov1alpha1.Rollout) bool {
	// Check MinReplicas
	if *current.Spec.MinReplicas != f.Spec.HorizontalScale.MinReplicas {
		return true
	}

	// Check MaxReplicas
	if current.Spec.MaxReplicas != f.Spec.HorizontalScale.MaxReplicas {
		return true
	}

	// Check TargetCPUUtilizationPercentage
	if *current.Spec.Metrics[0].Resource.Target.AverageUtilization != f.Spec.HorizontalScale.TargetCPUUtilizationPercentage {
		return true
	}

	return false
}

func updateHpa(hpa *autoscalingv2.HorizontalPodAutoscaler, f *oneclickiov1alpha1.Rollout) {
	// Update MinReplicas
	hpa.Spec.MinReplicas = &f.Spec.HorizontalScale.MinReplicas

	// Update MaxReplicas
	hpa.Spec.MaxReplicas = f.Spec.HorizontalScale.MaxReplicas

	// Update TargetCPUUtilizationPercentage
	hpa.Spec.Metrics[0].Resource.Target.AverageUtilization = &f.Spec.HorizontalScale.TargetCPUUtilizationPercentage
}
