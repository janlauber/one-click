package controllers

import (
	"context"
	"reflect"

	oneclickiov1alpha1 "github.com/janlauber/one-click/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *RolloutReconciler) reconcileService(ctx context.Context, f *oneclickiov1alpha1.Rollout) error {
	log := log.FromContext(ctx)

	// Keep track of services that should exist according to the Rollout spec
	expectedServices := make(map[string]bool)
	for _, intf := range f.Spec.Interfaces {

		expectedServices[intf.Name+"-svc"] = true
		// Process each interface
		service := r.serviceForRollout(f, intf)

		foundService := &corev1.Service{}
		err := r.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: f.Namespace}, foundService)
		if err != nil && errors.IsNotFound(err) {
			// If the Service is not found, create a new one
			log.Info("Creating a new Service", "Namespace", service.Namespace, "Name", service.Name)
			err = r.Create(ctx, service)
			if err != nil {
				// Handle creation error
				r.Recorder.Eventf(f, corev1.EventTypeWarning, "CreationFailed", "Failed to create Service %s", f.Name)
				log.Error(err, "Failed to create Service", "Namespace", service.Namespace, "Name", service.Name)
				return err
			}
			r.Recorder.Eventf(f, corev1.EventTypeNormal, "Created", "Created Service %s", f.Name)
		} else if err != nil {
			// Handle other errors
			r.Recorder.Eventf(f, corev1.EventTypeWarning, "GetFailed", "Failed to get Service %s", f.Name)
			log.Error(err, "Failed to get Service", "Namespace", service.Namespace, "Name", service.Name)
			return err
		} else {
			// If the Service exists, check if it needs to be updated
			desiredPorts := getServicePorts(intf)
			if !reflect.DeepEqual(foundService.Spec.Ports, desiredPorts) {
				foundService.Spec.Ports = desiredPorts
				err = r.Update(ctx, foundService)
				if err != nil {
					// Handle update error
					r.Recorder.Eventf(f, corev1.EventTypeWarning, "UpdateFailed", "Failed to update Service %s", foundService.Name)
					log.Error(err, "Failed to update Service", "Namespace", foundService.Namespace, "Name", foundService.Name)
					return err
				}
				r.Recorder.Eventf(f, corev1.EventTypeNormal, "Updated", "Updated Service %s", foundService.Name)
			}
		}
	}

	// Delete services that are no longer specified in the Rollout spec
	serviceList := &corev1.ServiceList{}
	listOpts := []client.ListOption{client.InNamespace(f.Namespace)}
	err := r.List(ctx, serviceList, listOpts...)
	if err != nil {
		log.Error(err, "Failed to list services", "Namespace", f.Namespace)
		return err
	}

	for _, service := range serviceList.Items {
		if _, exists := expectedServices[service.Name]; !exists {
			// Service is no longer needed, delete it
			if service.Labels["managed-by"] == "framework-operator" {
				err = r.Delete(ctx, &service)
				if err != nil {
					log.Error(err, "Failed to delete service", "Namespace", service.Namespace, "Name", service.Name)
					return err
				}
				log.Info("Deleted service", "Namespace", service.Namespace, "Name", service.Name)
			}
		}
	}

	return nil
}

func (r *RolloutReconciler) serviceForRollout(f *oneclickiov1alpha1.Rollout, intf oneclickiov1alpha1.InterfaceSpec) *corev1.Service {
	labels := map[string]string{
		"rollout.one-click.dev/name": f.Name,
		"managed-by":                 "framework-operator", // Unique label to identify operator-managed services
	}
	selectorLabels := map[string]string{
		"rollout.one-click.dev/name": f.Name,
	}
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      intf.Name + "-svc", // Create a unique name for the Service
			Namespace: f.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: selectorLabels,
			Ports:    getServicePorts(intf),
			Type:     corev1.ServiceTypeClusterIP, // Default to ClusterIP, modify if needed
		},
	}

	// Set Rollout instance as the owner and controller
	ctrl.SetControllerReference(f, svc, r.Scheme)
	return svc
}

func getServicePorts(intf oneclickiov1alpha1.InterfaceSpec) []corev1.ServicePort {
	return []corev1.ServicePort{
		{
			Name:       intf.Name,
			Port:       intf.Port,
			TargetPort: intstr.FromInt(int(intf.Port)),
			Protocol:   corev1.ProtocolTCP, // Defaulting to TCP, change as needed
		},
	}
}
