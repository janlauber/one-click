package controllers

import (
	"context"

	oneclickiov1 "github.com/janlauber/one-click/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *FrameworkReconciler) updateStatus(ctx context.Context, f *oneclickiov1.Framework) error {
	log := log.FromContext(ctx)

	// get the deployment replica count
	replicas, err := r.getDeploymentReplicas(ctx, f)
	if err != nil {
		log.Error(err, "Failed to get Deployment replica count")
		return err
	}

	// Get the Deployment
	deployment := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: f.Name, Namespace: f.Namespace}, deployment)
	if err != nil {
		log.Error(err, "Failed to get Deployment", "Framework.Namespace", f.Namespace, "Framework.Name", f.Name)
		return err
	}

	// List the Pods that are controlled by this Deployment
	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(deployment.Spec.Selector.MatchLabels)
	listOpts := []client.ListOption{
		client.InNamespace(f.Namespace),
		client.MatchingLabelsSelector{Selector: labelSelector},
	}

	err = r.List(ctx, podList, listOpts...)
	if err != nil {
		log.Error(err, "Failed to list pods", "Framework.Namespace", f.Namespace, "Framework.Name", f.Name)
		return err
	}

	// Collect the names of the Pods
	podNames := make([]string, 0)
	for _, pod := range podList.Items {
		podNames = append(podNames, pod.Name)
	}

	// Get the Deployment status if there are any pods pending or not ready
	deploymentStatus := "OK"
	for _, pod := range podList.Items {
		if pod.Status.Phase == corev1.PodPending {
			deploymentStatus = "Pending"
		}
		if pod.Status.Phase == corev1.PodRunning {
			for _, containerStatus := range pod.Status.ContainerStatuses {
				if !containerStatus.Ready {
					deploymentStatus = "NotReady"
				}
			}
		}
	}

	// add some test data to the status
	f.Status.Deployment.Replicas = replicas
	f.Status.Deployment.PodNames = podNames
	f.Status.Deployment.Status = deploymentStatus

	// List the Ingresses
	ingressList := &networkingv1.IngressList{}
	listOpts = []client.ListOption{
		client.InNamespace(f.Namespace),
		// If you have a specific label selector, uncomment the following line:
		// client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(labelSelector)},
	}

	err = r.List(ctx, ingressList, listOpts...)
	if err != nil {
		log.Error(err, "Failed to list ingresses", "Framework.Namespace", f.Namespace, "Framework.Name", f.Name)
		return err
	}

	// Collect the hostnames from the Ingress rules
	var hosts []string
	var ingressStatus string
	for _, ingress := range ingressList.Items {
		for _, rule := range ingress.Spec.Rules {
			if rule.Host != "" {
				hosts = append(hosts, rule.Host)
			}
		}
		ingressStatus = "Unknown"
		if len(ingress.Status.LoadBalancer.Ingress) > 0 {
			ingressStatus = "OK"
		}
	}

	f.Status.Ingresses = []oneclickiov1.IngressStatus{
		{
			Hosts:  hosts,
			Status: ingressStatus,
		},
	}

	// List the Services
	serviceList := &corev1.ServiceList{}
	listOpts = []client.ListOption{
		client.InNamespace(f.Namespace),
		// Use label selectors if applicable
	}

	err = r.List(ctx, serviceList, listOpts...)
	if err != nil {
		log.Error(err, "Failed to list services", "Framework.Namespace", f.Namespace, "Framework.Name", f.Name)
		return err
	}

	// Collect information from the Services
	var serviceStatuses []oneclickiov1.ServiceStatus
	for _, service := range serviceList.Items {
		var ports []int32
		for _, p := range service.Spec.Ports {
			ports = append(ports, p.Port)
		}

		// Basic status check - customize this as necessary
		status := "OK" // Example status, adjust based on your logic

		serviceStatus := oneclickiov1.ServiceStatus{
			Name:   service.Name,
			Ports:  ports,
			Status: status,
		}
		serviceStatuses = append(serviceStatuses, serviceStatus)
	}

	// Update the Framework status
	f.Status.Services = serviceStatuses

	// List the PVCs
	pvcList := &corev1.PersistentVolumeClaimList{}
	listOpts = []client.ListOption{
		client.InNamespace(f.Namespace),
		// Use label selectors if applicable
	}

	err = r.List(ctx, pvcList, listOpts...)
	if err != nil {
		log.Error(err, "Failed to list PVCs", "Framework.Namespace", f.Namespace, "Framework.Name", f.Name)
		return err
	}

	// Collect information from the PVCs
	var volumeStatuses []oneclickiov1.VolumeStatus
	for _, pvc := range pvcList.Items {
		// Basic status check - customize this as necessary
		status := "Unknown"
		if pvc.Status.Phase == corev1.ClaimBound {
			status = "Bound"
		} else if pvc.Status.Phase == corev1.ClaimPending {
			status = "Pending"
		} else if pvc.Status.Phase == corev1.ClaimLost {
			status = "Lost"
		}

		volumeStatus := oneclickiov1.VolumeStatus{
			Name:   pvc.Name,
			Status: status,
		}
		volumeStatuses = append(volumeStatuses, volumeStatus)
	}

	// Update the Framework status
	f.Status.Volumes = volumeStatuses

	// Update the Framework status
	if err := r.Status().Update(ctx, f); err != nil {
		log.Error(err, "Failed to update Framework status")
		return err
	}

	return nil
}

func (r *FrameworkReconciler) getDeploymentReplicas(ctx context.Context, f *oneclickiov1.Framework) (int32, error) {
	log := log.FromContext(ctx)

	// Get the deployment
	dep := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: f.Name, Namespace: f.Namespace}, dep)
	if err != nil {
		log.Error(err, "Failed to get Deployment")
		return 0, err
	}

	return dep.Status.Replicas, nil
}
