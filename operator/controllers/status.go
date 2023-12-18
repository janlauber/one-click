package controllers

import (
	"context"

	oneclickiov1alpha1 "github.com/janlauber/one-click/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *RolloutReconciler) updateStatus(ctx context.Context, f *oneclickiov1alpha1.Rollout) error {
	log := log.FromContext(ctx)

	// get the deployment replica count
	replicas, err := r.getDeploymentReplicas(ctx, f)
	if err != nil {
		log.Error(err, "Failed to get Deployment replica count")
		return err
	}

	// Get the Deployment
	// TODO: fix status at first run
	var requestSumCpu, requestSumMemory, limitSumCpu, limitSumMemory resource.Quantity

	deployment := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: f.Name, Namespace: f.Namespace}, deployment)
	if err != nil {
		log.Error(err, "Failed to get Deployment", "Rollout.Namespace", f.Namespace, "Rollout.Name", f.Name)
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
		log.Error(err, "Failed to list pods", "Rollout.Namespace", f.Namespace, "Rollout.Name", f.Name)
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

	for _, pod := range podList.Items {
		for _, container := range pod.Spec.Containers {
			// Sum requested resources
			if req, ok := container.Resources.Requests[corev1.ResourceCPU]; ok {
				requestSumCpu.Add(req)
			}
			if req, ok := container.Resources.Requests[corev1.ResourceMemory]; ok {
				requestSumMemory.Add(req)
			}

			// Sum limits
			if lim, ok := container.Resources.Limits[corev1.ResourceCPU]; ok {
				limitSumCpu.Add(lim)
			}
			if lim, ok := container.Resources.Limits[corev1.ResourceMemory]; ok {
				limitSumMemory.Add(lim)
			}
		}
	}

	// add some test data to the status
	f.Status.Deployment.Replicas = replicas
	f.Status.Deployment.PodNames = podNames
	f.Status.Deployment.Status = deploymentStatus
	// Update the Rollout status with resource information
	f.Status.Deployment.Resources = oneclickiov1alpha1.DeploymentResources{
		RequestSum: oneclickiov1alpha1.Resources{
			CPU:    requestSumCpu.AsDec().String(),
			Memory: requestSumMemory.AsDec().String(),
		},
		LimitSum: oneclickiov1alpha1.Resources{
			CPU:    limitSumCpu.AsDec().String(),
			Memory: limitSumMemory.AsDec().String(),
		},
	}

	// List the Ingresses
	ingressList := &networkingv1.IngressList{}
	listOpts = []client.ListOption{
		client.InNamespace(f.Namespace),
		// If you have a specific label selector, uncomment the following line:
		// client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(labelSelector)},
	}

	err = r.List(ctx, ingressList, listOpts...)
	if err != nil {
		log.Error(err, "Failed to list ingresses", "Rollout.Namespace", f.Namespace, "Rollout.Name", f.Name)
		return err
	}

	// Update the Rollout status with ingress information
	var ingressStatuses []oneclickiov1alpha1.IngressStatus
	for _, ingress := range ingressList.Items {
		// Filter out ingresses not related to the Rollout, if necessary

		var hosts []string
		for _, rule := range ingress.Spec.Rules {
			if rule.Host != "" {
				hosts = append(hosts, rule.Host)
			}
		}

		if len(hosts) > 0 { // Only add if there are hosts
			ingressStatus := oneclickiov1alpha1.IngressStatus{
				Hosts:  hosts,
				Status: determineIngressStatus(ingress), // Implement this function based on your logic
			}
			ingressStatuses = append(ingressStatuses, ingressStatus)
		}
	}

	f.Status.Ingresses = ingressStatuses

	// List the Services
	serviceList := &corev1.ServiceList{}
	listOpts = []client.ListOption{
		client.InNamespace(f.Namespace),
		// Use label selectors if applicable
	}

	err = r.List(ctx, serviceList, listOpts...)
	if err != nil {
		log.Error(err, "Failed to list services", "Rollout.Namespace", f.Namespace, "Rollout.Name", f.Name)
		return err
	}

	// Collect information from the Services
	var serviceStatuses []oneclickiov1alpha1.ServiceStatus
	for _, service := range serviceList.Items {
		var ports []int32
		for _, p := range service.Spec.Ports {
			ports = append(ports, p.Port)
		}

		// Basic status check - customize this as necessary
		status := "OK" // Example status, adjust based on your logic

		serviceStatus := oneclickiov1alpha1.ServiceStatus{
			Name:   service.Name,
			Ports:  ports,
			Status: status,
		}
		serviceStatuses = append(serviceStatuses, serviceStatus)
	}

	// Update the Rollout status
	f.Status.Services = serviceStatuses

	// List the PVCs
	pvcList := &corev1.PersistentVolumeClaimList{}
	listOpts = []client.ListOption{
		client.InNamespace(f.Namespace),
		// Use label selectors if applicable
	}

	err = r.List(ctx, pvcList, listOpts...)
	if err != nil {
		log.Error(err, "Failed to list PVCs", "Rollout.Namespace", f.Namespace, "Rollout.Name", f.Name)
		return err
	}

	// Collect information from the PVCs
	var volumeStatuses []oneclickiov1alpha1.VolumeStatus
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

		volumeStatus := oneclickiov1alpha1.VolumeStatus{
			Name:   pvc.Name,
			Status: status,
		}
		volumeStatuses = append(volumeStatuses, volumeStatus)
	}

	// Update the Rollout status
	f.Status.Volumes = volumeStatuses

	// Update the Rollout status
	if err := r.Status().Update(ctx, f); err != nil {
		log.Error(err, "Failed to update Rollout status")
		return err
	}

	return nil
}

func (r *RolloutReconciler) getDeploymentReplicas(ctx context.Context, f *oneclickiov1alpha1.Rollout) (int32, error) {
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

func determineIngressStatus(ingress networkingv1.Ingress) string {
	// Check if the ingress has been assigned a load balancer IP or hostname
	if len(ingress.Status.LoadBalancer.Ingress) > 0 {
		// Assuming that the presence of a load balancer IP or hostname indicates that the ingress is operational
		return "Operational"
	}

	// You can add more checks here depending on what you consider as part of the Ingress's status.
	// For example, you could check for specific conditions, annotations, or other status fields.

	// If the ingress does not meet your criteria for being operational, you can return another status
	return "Pending"
}
