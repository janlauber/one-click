package k8s

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// Watch all k8s resources with the specified label selector and send updates over WebSocket
func WatchK8sResourcesAndSendUpdates(ws *websocket.Conn, projectId string, rolloutId string) {

	// resources to watch "deployments", "replicasets", "pods", "services", "ingresses", "configmaps", "secrets", "persistentvolumeclaims", "persistentvolumes", "serviceaccounts"

	// Start watching resources with the specified label selector
	labelSelector := fmt.Sprintf("rollout.one-click.dev/name=%s", rolloutId)

	// deploymentWatchInterface, err := Clientset.AppsV1().Deployments(projectId).Watch(Ctx, metav1.ListOptions{
	// 	LabelSelector: labelSelector,
	// })
	// if err != nil {
	// 	log.Fatalf("Error setting up watch: %v", err)
	// }
	// defer deploymentWatchInterface.Stop()

	// replicaSetWatchInterface, err := Clientset.AppsV1().ReplicaSets(projectId).Watch(Ctx, metav1.ListOptions{
	// 	LabelSelector: labelSelector,
	// })
	// if err != nil {
	// 	log.Fatalf("Error setting up watch: %v", err)
	// }
	// defer replicaSetWatchInterface.Stop()

	podWatchInterface, err := Clientset.CoreV1().Pods(projectId).Watch(Ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Fatalf("Error setting up watch: %v", err)
	}
	defer podWatchInterface.Stop()

	serviceWatchInterface, err := Clientset.CoreV1().Services(projectId).Watch(Ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Fatalf("Error setting up watch: %v", err)
	}
	defer serviceWatchInterface.Stop()

	ingressWatchInterface, err := Clientset.NetworkingV1().Ingresses(projectId).Watch(Ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Fatalf("Error setting up watch: %v", err)
	}
	defer ingressWatchInterface.Stop()

	configMapWatchInterface, err := Clientset.CoreV1().ConfigMaps(projectId).Watch(Ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Fatalf("Error setting up watch: %v", err)
	}
	defer configMapWatchInterface.Stop()

	secretWatchInterface, err := Clientset.CoreV1().Secrets(projectId).Watch(Ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Fatalf("Error setting up watch: %v", err)
	}
	defer secretWatchInterface.Stop()

	pvcWatchInterface, err := Clientset.CoreV1().PersistentVolumeClaims(projectId).Watch(Ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Fatalf("Error setting up watch: %v", err)
	}
	defer pvcWatchInterface.Stop()

	pvWatchInterface, err := Clientset.CoreV1().PersistentVolumes().Watch(Ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Fatalf("Error setting up watch: %v", err)
	}
	defer pvWatchInterface.Stop()

	saWatchInterface, err := Clientset.CoreV1().ServiceAccounts(projectId).Watch(Ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Fatalf("Error setting up watch: %v", err)
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
			log.Printf("WebSocket write error: %v", err)
			return
		}
	}
}
