package controllers

import (
	"context"

	oneclickiov1 "github.com/janlauber/one-click/api/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *FrameworkReconciler) reconcileHPA(ctx context.Context, f *oneclickiov1.Framework) error {
	log := log.FromContext(ctx)

	// Construct the desired HPA object based on the Framework specification
	desiredHpa := r.hpaForFramework(f)

	// Try to fetch the existing HPA
	foundHpa := &autoscalingv2.HorizontalPodAutoscaler{}
	err := r.Get(ctx, types.NamespacedName{Name: f.Name, Namespace: f.Namespace}, foundHpa)
	if err != nil && errors.IsNotFound(err) {
		// If the HPA is not found, create a new one
		log.Info("Creating a new HPA", "Namespace", desiredHpa.Namespace, "Name", desiredHpa.Name)
		err = r.Create(ctx, desiredHpa)
		if err != nil {
			// Handle creation error
			log.Error(err, "Failed to create HPA", "Namespace", desiredHpa.Namespace, "Name", desiredHpa.Name)
			return err
		}
	} else if err != nil {
		// Handle other errors
		log.Error(err, "Failed to get HPA", "Namespace", desiredHpa.Namespace, "Name", desiredHpa.Name)
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
				return err
			}
		}
	}

	return nil
}

func (r *FrameworkReconciler) hpaForFramework(f *oneclickiov1.Framework) *autoscalingv2.HorizontalPodAutoscaler {
	return &autoscalingv2.HorizontalPodAutoscaler{
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
}

func needsHpaUpdate(current *autoscalingv2.HorizontalPodAutoscaler, f *oneclickiov1.Framework) bool {
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

func updateHpa(hpa *autoscalingv2.HorizontalPodAutoscaler, f *oneclickiov1.Framework) {
	// Update MinReplicas
	hpa.Spec.MinReplicas = &f.Spec.HorizontalScale.MinReplicas

	// Update MaxReplicas
	hpa.Spec.MaxReplicas = f.Spec.HorizontalScale.MaxReplicas

	// Update TargetCPUUtilizationPercentage
	hpa.Spec.Metrics[0].Resource.Target.AverageUtilization = &f.Spec.HorizontalScale.TargetCPUUtilizationPercentage
}