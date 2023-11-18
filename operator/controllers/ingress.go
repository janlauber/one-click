package controllers

import (
	"context"
	"reflect"

	oneclickiov1 "github.com/janlauber/one-click/api/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *FrameworkReconciler) reconcileIngress(ctx context.Context, f *oneclickiov1.Framework) error {
	log := log.FromContext(ctx)

	for _, intf := range f.Spec.Interfaces {
		// Process each interface
		ingress := r.ingressForFramework(f, intf)

		foundIngress := &networkingv1.Ingress{}
		err := r.Get(ctx, types.NamespacedName{Name: ingress.Name, Namespace: f.Namespace}, foundIngress)
		if err != nil && errors.IsNotFound(err) {
			// If the Ingress is not found, create a new one
			log.Info("Creating a new Ingress", "Namespace", ingress.Namespace, "Name", ingress.Name)
			err = r.Create(ctx, ingress)
			if err != nil {
				// Handle creation error
				log.Error(err, "Failed to create Ingress", "Namespace", ingress.Namespace, "Name", ingress.Name)
				return err
			}
		} else if err != nil {
			// Handle other errors
			log.Error(err, "Failed to get Ingress", "Namespace", ingress.Namespace, "Name", ingress.Name)
			return err
		} else {
			// If the Ingress exists, check if it needs to be updated
			desiredRules := getIngressRules(intf)
			desiredTLS := getIngressTLS(intf)
			if !reflect.DeepEqual(foundIngress.Spec.Rules, desiredRules) || !reflect.DeepEqual(foundIngress.Spec.TLS, desiredTLS) {
				foundIngress.Spec.Rules = desiredRules
				foundIngress.Spec.TLS = desiredTLS
				err = r.Update(ctx, foundIngress)
				if err != nil {
					// Handle update error
					log.Error(err, "Failed to update Ingress", "Namespace", foundIngress.Namespace, "Name", foundIngress.Name)
					return err
				}
			}
		}
	}

	return nil
}

func (r *FrameworkReconciler) ingressForFramework(f *oneclickiov1.Framework, intf oneclickiov1.InterfaceSpec) *networkingv1.Ingress {
	labels := map[string]string{"app": f.Name}
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        intf.Name + "-ingress", // Create a unique name for the Ingress
			Namespace:   f.Namespace,
			Labels:      labels,
			Annotations: make(map[string]string),
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{},
			TLS:   []networkingv1.IngressTLS{},
		},
	}

	// Add ingress class if defined
	if intf.Ingress.IngressClass != "" {
		ingress.Spec.IngressClassName = &intf.Ingress.IngressClass
	}

	// Add annotations if defined
	if len(intf.Ingress.Annotations) > 0 {
		for k, v := range intf.Ingress.Annotations {
			ingress.Annotations[k] = v
		}
	}

	// Define the ingress rules
	for _, rule := range intf.Ingress.Rules {
		ingressRule := networkingv1.IngressRule{
			Host: rule.Host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{
						{
							Path: rule.Path,
							PathType: func() *networkingv1.PathType {
								pt := networkingv1.PathTypeImplementationSpecific // or PathTypeExact or PathTypePrefix
								return &pt
							}(),
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: intf.Name + "-svc",
									Port: networkingv1.ServiceBackendPort{
										Number: intf.Port,
									},
								},
							},
						},
					},
				},
			},
		}
		ingress.Spec.Rules = append(ingress.Spec.Rules, ingressRule)

		// Add TLS configuration if TLS is enabled for this ingress path
		if rule.TLS {
			tls := networkingv1.IngressTLS{
				Hosts:      []string{rule.Host},
				SecretName: intf.Name + "-tls-secret", // Name of the TLS secret
			}
			ingress.Spec.TLS = append(ingress.Spec.TLS, tls)
		}
	}

	// Set Framework instance as the owner and controller
	ctrl.SetControllerReference(f, ingress, r.Scheme)
	return ingress
}

func getIngressRules(intf oneclickiov1.InterfaceSpec) []networkingv1.IngressRule {
	var rules []networkingv1.IngressRule

	for _, rule := range intf.Ingress.Rules {
		ingressRule := networkingv1.IngressRule{
			Host: rule.Host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{
						{
							Path: rule.Path,
							PathType: func() *networkingv1.PathType {
								pt := networkingv1.PathTypeImplementationSpecific // or PathTypeExact or PathTypePrefix
								return &pt
							}(),
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: intf.Name + "-svc",
									Port: networkingv1.ServiceBackendPort{
										Number: intf.Port,
									},
								},
							},
						},
					},
				},
			},
		}
		rules = append(rules, ingressRule)
	}

	return rules
}

func getIngressTLS(intf oneclickiov1.InterfaceSpec) []networkingv1.IngressTLS {
	var tlsConfigs []networkingv1.IngressTLS

	// Loop over each rule defined in the ingress path
	for _, rule := range intf.Ingress.Rules {
		// Add TLS configuration if TLS is enabled for this ingress path
		if rule.TLS {
			tls := networkingv1.IngressTLS{
				Hosts:      []string{rule.Host},
				SecretName: intf.Name + "-tls-secret", // Name of the TLS secret
			}
			tlsConfigs = append(tlsConfigs, tls)
		}
	}

	return tlsConfigs
}
