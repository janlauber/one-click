package k8s

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func WatchPodsAndSendUpdates(ws *websocket.Conn, projectId string, rolloutId string) {
	// Prepare the label selector string
	labelSelector := fmt.Sprintf("rollout.one-click.dev/name=%s", rolloutId)

	// Start watching Pods with the specified label selector
	watchInterface, err := Clientset.CoreV1().Pods(projectId).Watch(Ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Fatalf("Error setting up watch: %v", err)
	}

	defer watchInterface.Stop()

	for event := range watchInterface.ResultChan() {
		switch event.Type {
		case watch.Added, watch.Modified, watch.Deleted, watch.Error:
			pod, ok := event.Object.(*corev1.Pod)
			if !ok {
				log.Printf("Unexpected type")
				continue
			}

			// Prepare and marshal the pod information
			podInfo := map[string]interface{}{
				"name":      pod.Name,
				"namespace": pod.Namespace,
				"status":    pod.Status.Phase,
				"labels":    pod.Labels,
			}
			msg, err := json.Marshal(podInfo)
			if err != nil {
				log.Printf("Error marshaling pod info: %v", err)
				continue
			}

			// Send the message over WebSocket
			if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("WebSocket write error: %v", err)
				// Handle WebSocket disconnection or errors as needed
				return
			}
		}
	}
}
