package k8s

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"

	yaml2 "github.com/ghodss/yaml"
	"github.com/janlauber/one-click/pkg/models"
	pb_models "github.com/pocketbase/pocketbase/models"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

func CreateOrUpdateRollout(rolloutId string, user *pb_models.Record, projectId string, deploymentId string, manifest string) error {
	if rolloutId == "" {
		return fmt.Errorf("rolloutId is required")
	}

	if user == nil {
		return fmt.Errorf("user is required")
	}

	if projectId == "" {
		return fmt.Errorf("projectId is required")
	}

	if deploymentId == "" {
		return fmt.Errorf("deploymentId is required")
	}

	if manifest == "" {
		return fmt.Errorf("manifest is required")
	}

	// manifest is a json string, convert to yaml
	yamlBytes, err := yaml2.JSONToYAML([]byte(manifest))
	if err != nil {
		return fmt.Errorf("error converting JSON to YAML: %w", err)
	}

	// Decode the YAML manifest into an unstructured object
	decUnstructured := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	_, _, err = decUnstructured.Decode(yamlBytes, nil, obj)
	if err != nil {
		return fmt.Errorf("error decoding YAML: %w", err)
	}

	// Set the name and namespace to rolloutId and deploymentId
	obj.SetName(deploymentId)
	obj.SetNamespace(projectId)
	// Define the labels for the Rollout object
	labels := map[string]string{
		"one-click.dev/projectId":    projectId,
		"one-click.dev/deploymentId": deploymentId,
		"one-click.dev/rolloutId":    rolloutId,
	}
	obj.SetLabels(labels)

	// Define the GroupVersionResource for the Rollout object
	rolloutGVR := schema.GroupVersionResource{
		Group:    "one-click.dev",
		Version:  "v1alpha1",
		Resource: "rollouts",
	}

	// Create Namespace
	nparams := NamespaceParams{
		Name:       projectId,
		UserRecord: user,
	}
	err = CreateNamespace(nparams)
	if err != nil {
		log.Println(err)
	}

	// Try to get the existing Rollout
	existingRollout, err := DynamicClient.Resource(rolloutGVR).Namespace(projectId).Get(Ctx, deploymentId, metav1.GetOptions{})
	if err != nil {
		// If not found, create it
		_, err = DynamicClient.Resource(rolloutGVR).Namespace(projectId).Create(Ctx, obj, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("error creating Rollout: %w", err)
		}
	} else {
		// Set the resourceVersion in the object to be updated
		obj.SetResourceVersion(existingRollout.GetResourceVersion())

		// Update the Rollout
		_, err = DynamicClient.Resource(rolloutGVR).Namespace(projectId).Update(Ctx, obj, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("error updating Rollout: %w", err)
		}
	}

	return nil
}

func DeleteRollout(projectId string, deploymentId string) error {
	// Define the GroupVersionResource for the Rollout object
	rolloutGVR := schema.GroupVersionResource{
		Group:    "one-click.dev",
		Version:  "v1alpha1",
		Resource: "rollouts",
	}

	// Delete the Rollout
	err := DynamicClient.Resource(rolloutGVR).Namespace(projectId).Delete(Ctx, deploymentId, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("error deleting Rollout: %w", err)
	}

	return nil
}

func GetRolloutStatus(projectId string, deploymentId string) (*models.RolloutStatus, error) {
	// Define the GroupVersionResource for the Rollout object
	rolloutGVR := schema.GroupVersionResource{
		Group:    "one-click.dev",
		Version:  "v1alpha1",
		Resource: "rollouts",
	}

	// Get the Rollout
	rollout, err := DynamicClient.Resource(rolloutGVR).Namespace(projectId).Get(Ctx, deploymentId, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting Rollout: %w", err)
	}

	// Extract the status field
	statusObj, found, err := unstructured.NestedMap(rollout.Object, "status")
	if err != nil || !found {
		return nil, fmt.Errorf("error getting Rollout status: %w", err)
	}

	// Convert the status to the RolloutStatus struct
	var rolloutStatus models.RolloutStatus
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(statusObj, &rolloutStatus)
	if err != nil {
		return nil, fmt.Errorf("error converting status to RolloutStatus: %w", err)
	}

	return &rolloutStatus, nil
}

func GetRolloutMetrics(projectId string, deploymentId string) (*models.PodMetricsResponse, error) {

	// List all pods in the projectId namespaced controlled by the deploymentId
	pods, err := Clientset.CoreV1().Pods(projectId).List(Ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("one-click.dev/deploymentId=%s", projectId),
	})

	if err != nil {
		return nil, fmt.Errorf("error getting pods: %w", err)
	}

	var podMetricsResponse models.PodMetricsResponse
	var podMetrics []models.PodMetrics

	for _, pod := range pods.Items {
		metrics, err := MetricsClient.MetricsV1beta1().PodMetricses(projectId).Get(Ctx, pod.Name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("error getting pod metrics: %w", err)
		}

		podMetrics = append(podMetrics, models.PodMetrics{
			Name:   pod.Name,
			CPU:    metrics.Containers[0].Usage.Cpu().AsDec().String(),
			Memory: metrics.Containers[0].Usage.Memory().AsDec().String(),
		})
	}

	podMetricsResponse.Metrics = podMetrics

	return &podMetricsResponse, nil
}

func GetRolloutEvents(projectId string, deploymentId string) (*models.EventResponse, error) {
	// List all events in the projectId namespaced controlled by the deploymentId deployment
	events, err := Clientset.CoreV1().Events(projectId).List(Ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s", deploymentId),
	})

	if err != nil {
		return nil, fmt.Errorf("error getting events: %w", err)
	}

	var eventResponse models.EventResponse
	var eventList []models.Event

	for _, event := range events.Items {
		eventList = append(eventList, models.Event{
			Reason:  event.Reason,
			Message: event.Message,
			Typus:   event.Type,
		})
	}

	eventResponse.Events = eventList

	// check if there is null
	if eventResponse.Events == nil {
		eventResponse.Events = []models.Event{}
	}

	return &eventResponse, nil
}

func GetRolloutLogs(w http.ResponseWriter, projectId string, podName string) error {

	// Get live logs
	liveLogOptions := &corev1.PodLogOptions{
		Follow: true,
	}
	liveReq := Clientset.CoreV1().Pods(projectId).GetLogs(podName, liveLogOptions)
	liveLogs, err := liveReq.Stream(context.Background())
	if err != nil {
		return err
	}

	// Write live logs line by line
	scanner := bufio.NewScanner(liveLogs)
	for scanner.Scan() {
		line := scanner.Text()
		_, err := fmt.Fprintf(w, "%s\n\n", line)
		if err != nil {
			return err
		}
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}

	return nil
}
