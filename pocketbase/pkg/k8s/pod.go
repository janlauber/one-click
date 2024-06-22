package k8s

import (
	"log"

	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

// DeletePod is a helper function to delete a pod
func DeletePod(namespace string, podName string) error {
	return Clientset.CoreV1().Pods(namespace).Delete(Ctx, podName, metav1.DeleteOptions{})
}

// WatchK8sTerminalAndSendUpdates watches the Kubernetes terminal and sends updates over WebSocket
func WatchK8sTerminalAndSendUpdates(conn *websocket.Conn, projectId string, podName string) {
	log.Printf("WatchK8sTerminalAndSendUpdates called with projectId=%s, podName=%s", projectId, podName)

	req := Clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(projectId).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Command: []string{"/bin/sh"},
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(Kubeconfig, "POST", req.URL())
	if err != nil {
		log.Println("Error creating executor:", err)
		return
	}

	streamOptions := remotecommand.StreamOptions{
		Stdin:  &wsStream{Conn: conn},
		Stdout: &wsStream{Conn: conn},
		Stderr: &wsStream{Conn: conn},
		Tty:    true,
	}

	err = exec.StreamWithContext(Ctx, streamOptions)
	if err != nil {
		log.Println("Error starting executor:", err)
		return
	}
}

type wsStream struct {
	Conn    *websocket.Conn
	readBuf []byte
}

func (w *wsStream) Read(p []byte) (n int, err error) {
	if len(w.readBuf) == 0 {
		_, message, err := w.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from WebSocket:", err)
			return 0, err
		}
		w.readBuf = append(w.readBuf, message...)
		log.Printf("Read %d bytes from WebSocket: %s", len(message), message)
	}

	n = copy(p, w.readBuf)
	w.readBuf = w.readBuf[n:]
	return n, nil
}

func (w *wsStream) Write(p []byte) (n int, err error) {
	err = w.Conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		log.Println("Error writing message to WebSocket:", err)
		return 0, err
	}
	log.Printf("Sent output: %s", string(p))
	return len(p), nil
}
