package watcher

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/janlauber/one-click/pkg/k8s"
	"github.com/labstack/echo/v5"
)

type RolloutTerminalRequest struct {
	ProjectId string `json:"projectId"`
	PodName   string `json:"podName"`
}

func WsK8sTerminalHandler(c echo.Context) error {
	conn, err := Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return err
	}
	defer func() {
		conn.Close()
		log.Println("WebSocket connection closed")
	}()

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("WebSocket read error:", err)
		return err
	}

	log.Printf("Received initial message: %s", string(msg))

	var request RolloutTerminalRequest
	err = json.Unmarshal(msg, &request)
	if err != nil {
		log.Println("WebSocket unmarshal error:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("{\"error\":\"Invalid request format\"}"))
		return err
	}

	log.Printf("Parsed request: projectId=%s, podName=%s", request.ProjectId, request.PodName)

	go k8s.WatchK8sTerminalAndSendUpdates(conn, request.ProjectId, request.PodName)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket connection closed:", err)
			break
		}
		log.Printf("WebSocket received message: %s", string(msg))
	}

	return nil
}
