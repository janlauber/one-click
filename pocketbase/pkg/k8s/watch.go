package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// Watch all k8s resources with the specified label selector and send updates over WebSocket
func WatchK8sResourcesAndSendUpdates(ws *websocket.Conn, projectId string, labelSelector string) {

	if projectId == "" {
		log.Printf("Error: projectId must not be empty")
		return
	}

	// deploymentWatchInterface, err := Clientset.AppsV1().Deployments(projectId).Watch(Ctx, metav1.ListOptions{
	// 	LabelSelector: labelSelector,
	// })
	// if err != nil {
	// 	log.Printf("Error setting up watch: %v", err)
	// }
	// defer deploymentWatchInterface.Stop()

	// replicaSetWatchInterface, err := Clientset.AppsV1().ReplicaSets(projectId).Watch(Ctx, metav1.ListOptions{
	// 	LabelSelector: labelSelector,
	// })
	// if err != nil {
	// 	log.Printf("Error setting up watch: %v", err)
	// }
	// defer replicaSetWatchInterface.Stop()

	podWatchInterface, err := Clientset.CoreV1().Pods(projectId).Watch(Ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Printf("Error setting up watch: %v", err)
		return
	}
	defer podWatchInterface.Stop()

	serviceWatchInterface, err := Clientset.CoreV1().Services(projectId).Watch(Ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Printf("Error setting up watch: %v", err)
		return
	}
	defer serviceWatchInterface.Stop()

	ingressWatchInterface, err := Clientset.NetworkingV1().Ingresses(projectId).Watch(Ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Printf("Error setting up watch: %v", err)
		return
	}
	defer ingressWatchInterface.Stop()

	configMapWatchInterface, err := Clientset.CoreV1().ConfigMaps(projectId).Watch(Ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Printf("Error setting up watch: %v", err)
		return
	}
	defer configMapWatchInterface.Stop()

	secretWatchInterface, err := Clientset.CoreV1().Secrets(projectId).Watch(Ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Printf("Error setting up watch: %v", err)
		return
	}
	defer secretWatchInterface.Stop()

	pvcWatchInterface, err := Clientset.CoreV1().PersistentVolumeClaims(projectId).Watch(Ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Printf("Error setting up watch: %v", err)
		return
	}
	defer pvcWatchInterface.Stop()

	pvWatchInterface, err := Clientset.CoreV1().PersistentVolumes().Watch(Ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Printf("Error setting up watch: %v", err)
		return
	}
	defer pvWatchInterface.Stop()

	saWatchInterface, err := Clientset.CoreV1().ServiceAccounts(projectId).Watch(Ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Printf("Error setting up watch: %v", err)
		return
	}
	defer saWatchInterface.Stop()

	for {
		select {
		// case event, ok := <-deploymentWatchInterface.ResultChan():
		// 	if !ok {
		// 		return
		// 	}
		// 	sendResourceUpdate(ws, event, "deployment")

		// case event, ok := <-replicaSetWatchInterface.ResultChan():
		// 	if !ok {
		// 		return
		// 	}
		// 	sendResourceUpdate(ws, event, "replicaset")

		case event, ok := <-podWatchInterface.ResultChan():
			if !ok {
				return
			}
			sendResourceUpdate(ws, event, "pod")

		case event, ok := <-serviceWatchInterface.ResultChan():
			if !ok {
				return
			}
			sendResourceUpdate(ws, event, "service")

		case event, ok := <-ingressWatchInterface.ResultChan():
			if !ok {
				return
			}
			sendResourceUpdate(ws, event, "ingress")

		case event, ok := <-configMapWatchInterface.ResultChan():
			if !ok {
				return
			}
			sendResourceUpdate(ws, event, "configmap")

		case event, ok := <-secretWatchInterface.ResultChan():
			if !ok {
				return
			}
			sendResourceUpdate(ws, event, "secret")

		case event, ok := <-pvcWatchInterface.ResultChan():
			if !ok {
				return
			}
			sendResourceUpdate(ws, event, "pvc")

		case event, ok := <-pvWatchInterface.ResultChan():
			if !ok {
				return
			}
			sendResourceUpdate(ws, event, "pv")

		case event, ok := <-saWatchInterface.ResultChan():
			if !ok {
				return
			}
			sendResourceUpdate(ws, event, "sa")
		}
	}
}

func sendResourceUpdate(ws *websocket.Conn, event watch.Event, resourceType string) {
	switch event.Type {
	case watch.Added, watch.Modified, watch.Deleted, watch.Error:
		resource, ok := event.Object.(metav1.Object)
		if !ok {
			log.Printf("Unexpected type")
			return
		}

		object, ok := event.Object.(metav1.Object)
		if !ok {
			log.Printf("Unexpected type")
			return
		}

		// Prepare and marshal the resource information
		resourceInfo := map[string]interface{}{
			"kind":      resourceType,
			"name":      resource.GetName(),
			"namespace": resource.GetNamespace(),
			"labels":    resource.GetLabels(),
			"status":    event.Type,
			"object":    object,
		}
		msg, err := json.Marshal(resourceInfo)
		if err != nil {
			log.Printf("Error marshaling resource info: %v", err)
			return
		}

		// Send the message over WebSocket
		if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

// Watch the logs of a pod and send updates over WebSocket
func WatchK8sLogsAndSendUpdates(ws *websocket.Conn, projectId string, podName string, ctx context.Context) {
	if projectId == "" || podName == "" {
		ws.WriteMessage(websocket.TextMessage, []byte("Error: projectId and podName must not be empty"))
		log.Printf("Error: projectId and podName must not be empty")
		return
	}

	// Start watching logs of the specified pod with the provided context
	req := Clientset.CoreV1().Pods(projectId).GetLogs(podName, &v1.PodLogOptions{
		Follow: true,
	})
	readCloser, err := req.Stream(ctx)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error getting logs: "+err.Error()))
		log.Printf("Error getting logs: %v", err)
		return // Return here to prevent further execution when there's an error
	}
	if readCloser == nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error: readCloser is nil"))
		log.Printf("Error: readCloser is nil")
		return // Prevent nil pointer dereference
	}
	defer readCloser.Close()

	buf := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			ws.WriteMessage(websocket.TextMessage, []byte("Stream closed by context cancellation"))
			return // Exit if context is cancelled
		default:
			n, err := readCloser.Read(buf)
			if err != nil {
				if err == io.EOF {
					// End of file reached, stop reading
					return
				}
				ws.WriteMessage(websocket.TextMessage, []byte("Error reading from readCloser: "+err.Error()))
				log.Printf("Error reading from readCloser: %v", err)
				return
			}

			// Send the message over WebSocket
			if err := ws.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
				log.Printf("Error sending logs over WebSocket: %v", err)
				return
			}
		}
	}
}

// Watch the events of a rollout and send updates over WebSocket
func WatchK8sEventsAndSendUpdates(ws *websocket.Conn, projectId string, kind string, name string) {

	if projectId == "" || kind == "" || name == "" {
		log.Printf("Error: projectId, kind and name must not be empty")
		return
	}

	// Start watching events of the specified rollout object with the specified kind and name
	req, err := Clientset.CoreV1().Events(projectId).Watch(Ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.kind=%s,involvedObject.name=%s", kind, name),
	})
	if err != nil {
		log.Printf("Error setting up watch: %v", err)
		return
	}
	defer req.Stop()

	for event := range req.ResultChan() {
		sendEventUpdate(ws, event)
	}
}

func sendEventUpdate(ws *websocket.Conn, event watch.Event) {
	switch event.Type {
	case watch.Added, watch.Modified, watch.Deleted, watch.Error:
		eventObj, ok := event.Object.(*v1.Event)
		if !ok {
			log.Printf("Unexpected type")
			return
		}

		// Prepare and marshal the event information
		eventInfo := map[string]interface{}{
			"reason":  eventObj.Reason,
			"message": eventObj.Message,
			"typus":   eventObj.Type,
		}
		msg, err := json.Marshal(eventInfo)
		if err != nil {
			log.Printf("Error marshaling event info: %v", err)
			return
		}

		// Send the message over WebSocket
		if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}
