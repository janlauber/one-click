package controllers

import (
	oneclickiov1 "github.com/janlauber/one-click/api/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
