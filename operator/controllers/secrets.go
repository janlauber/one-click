package controllers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	oneclickiov1 "github.com/janlauber/one-click/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *FrameworkReconciler) reconcileSecret(ctx context.Context, f *oneclickiov1.Framework) error {
	// Get the desired state of the Secret from the helper function
	desiredSecret := r.secretForFramework(f)

	// Check if the Secret already exists
	foundSecret := &corev1.Secret{}
	err := r.Get(ctx, types.NamespacedName{Name: desiredSecret.Name, Namespace: f.Namespace}, foundSecret)
	if err != nil && errors.IsNotFound(err) {
		// Create the Secret
		err = r.Create(ctx, desiredSecret)
		if err != nil {
			r.Recorder.Eventf(f, corev1.EventTypeWarning, "CreationFailed", "Failed to create Secret %s", f.Name)
			log.Log.Error(err, "Failed to create Secret", "Namespace", desiredSecret.Namespace, "Name", desiredSecret.Name)
			return err
		}
		r.Recorder.Eventf(f, corev1.EventTypeNormal, "Created", "Created Secret %s", f.Name)
	} else if err != nil {
		r.Recorder.Eventf(f, corev1.EventTypeWarning, "GetFailed", "Failed to get Secret %s", f.Name)
		log.Log.Error(err, "Failed to get Secret", "Namespace", desiredSecret.Namespace, "Name", desiredSecret.Name)
		return err
	} else {
		// Secret exists - check if it needs an update
		if needsSecretUpdate(foundSecret, f) {
			updateSecret(foundSecret, f)
			err = r.Update(ctx, foundSecret)
			if err != nil {
				r.Recorder.Eventf(f, corev1.EventTypeWarning, "UpdateFailed", "Failed to update Secret %s", foundSecret.Name)
				log.Log.Error(err, "Failed to update Secret", "Namespace", foundSecret.Namespace, "Name", foundSecret.Name)
				return err
			}
			r.Recorder.Eventf(f, corev1.EventTypeNormal, "Updated", "Updated Secret %s", foundSecret.Name)
		}
	}

	return nil
}

func needsSecretUpdate(foundSecret *corev1.Secret, f *oneclickiov1.Framework) bool {
	// If no secrets are specified in the Framework spec and the found secret is not empty, it needs to be deleted
	if len(f.Spec.Secrets) == 0 && len(foundSecret.Data) > 0 {
		return true
	}

	// If the number of keys in the secret doesn't match the number of specified secrets, update is needed
	if len(foundSecret.Data) != len(f.Spec.Secrets) {
		return true
	}

	// Check if the content of the secret matches the Framework spec
	for _, secretItem := range f.Spec.Secrets {
		key := strings.TrimSpace(secretItem.Name)
		if string(foundSecret.Data[key]) != secretItem.Value {
			return true
		}
	}

	// Check if the secrets checksum annotation needs an update
	secretsChecksum := calculateSecretsChecksum(f.Spec.Secrets)
	currentChecksum, exists := foundSecret.Annotations["oneclick.io/secrets-checksum"]
	if !exists || currentChecksum != secretsChecksum {
		return true
	}

	// No update needed
	return false
}

func calculateSecretsChecksum(secrets []oneclickiov1.SecretItem) string {
	hash := sha256.New()
	for _, secret := range secrets {
		data := secret.Name + "=" + secret.Value
		hash.Write([]byte(data))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func updateSecret(foundSecret *corev1.Secret, f *oneclickiov1.Framework) {
	if len(f.Spec.Secrets) == 0 {
		// Clear the secret data if no secrets are specified
		foundSecret.Data = make(map[string][]byte)
		return
	}

	// Update or add new data to the secret
	foundSecret.Data = make(map[string][]byte)
	for _, secretItem := range f.Spec.Secrets {
		key := strings.TrimSpace(secretItem.Name)
		foundSecret.Data[key] = []byte(secretItem.Value)
	}
}

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
