package k8s

import (
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

func CreateOrUpdateRollout(projectId string, rollout string) error {
	// Define the path to the YAML file
	yamlFilePath := filepath.Join(".rollouts", projectId, rollout+".yaml")

	// Read the YAML manifest
	yamlFile, err := os.ReadFile(yamlFilePath)
	if err != nil {
		return fmt.Errorf("error reading YAML file: %w", err)
	}

	// Decode the YAML manifest into an unstructured object
	decUnstructured := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	_, _, err = decUnstructured.Decode(yamlFile, nil, obj)
	if err != nil {
		return fmt.Errorf("error decoding YAML: %w", err)
	}

	// Define the GroupVersionResource for the Rollout object
	rolloutGVR := schema.GroupVersionResource{
		Group:    "one-click.io",
		Version:  "v1alpha1",
		Resource: "rollouts",
	}

	namespace := obj.GetNamespace()
	rolloutName := obj.GetName()

	// Try to get the existing Rollout
	existingRollout, err := DynamicClient.Resource(rolloutGVR).Namespace(namespace).Get(Ctx, rolloutName, metav1.GetOptions{})
	if err != nil {
		// If not found, create it
		_, err = DynamicClient.Resource(rolloutGVR).Namespace(namespace).Create(Ctx, obj, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("error creating Rollout: %w", err)
		}
	} else {
		// Set the resourceVersion in the object to be updated
		obj.SetResourceVersion(existingRollout.GetResourceVersion())

		// Update the Rollout
		_, err = DynamicClient.Resource(rolloutGVR).Namespace(namespace).Update(Ctx, obj, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("error updating Rollout: %w", err)
		}
	}

	return nil
}
