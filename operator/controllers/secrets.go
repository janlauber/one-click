package controllers

import (
	"strings"

	oneclickiov1 "github.com/janlauber/one-click/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *FrameworkReconciler) secretForFramework(f *oneclickiov1.Framework) *corev1.Secret {
	secretData := make(map[string]string)
	for _, secretItem := range f.Spec.Secrets {
		// Ensure the secret key is valid
		key := strings.TrimSpace(secretItem.Name)
		secretData[key] = secretItem.Value
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      f.Name + "-secrets", // Naming the secret based on the Framework name
			Namespace: f.Namespace,
		},
		StringData: secretData,
	}
}
