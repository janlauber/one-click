package controllers

import (
	"context"
	"reflect"

	oneclickiov1 "github.com/janlauber/one-click/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *FrameworkReconciler) reconcileService(ctx context.Context, f *oneclickiov1.Framework) error {
	log := log.FromContext(ctx)

	for _, intf := range f.Spec.Interfaces {
		// Process each interface
		service := r.serviceForFramework(f, intf)

		foundService := &corev1.Service{}
		err := r.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: f.Namespace}, foundService)
		if err != nil && errors.IsNotFound(err) {
			// If the Service is not found, create a new one
			log.Info("Creating a new Service", "Namespace", service.Namespace, "Name", service.Name)
			err = r.Create(ctx, service)
			if err != nil {
				// Handle creation error
				log.Error(err, "Failed to create Service", "Namespace", service.Namespace, "Name", service.Name)
				return err
			}
		} else if err != nil {
			// Handle other errors
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
					log.Error(err, "Failed to update Service", "Namespace", foundService.Namespace, "Name", foundService.Name)
					return err
				}
			}
		}
	}

	return nil
}

func (r *FrameworkReconciler) serviceForFramework(f *oneclickiov1.Framework, intf oneclickiov1.InterfaceSpec) *corev1.Service {
	labels := map[string]string{"app": f.Name}
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      intf.Name + "-svc", // Create a unique name for the Service
			Namespace: f.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports:    getServicePorts(intf),
			Type:     corev1.ServiceTypeClusterIP, // Default to ClusterIP, modify if needed
		},
	}

	// Set Framework instance as the owner and controller
	ctrl.SetControllerReference(f, svc, r.Scheme)
	return svc
}

func getServicePorts(intf oneclickiov1.InterfaceSpec) []corev1.ServicePort {
	return []corev1.ServicePort{
		{
			Name:       intf.Name,
			Port:       intf.Port,
			TargetPort: intstr.FromInt(int(intf.Port)),
			Protocol:   corev1.ProtocolTCP, // Defaulting to TCP, change as needed
		},
	}
}
