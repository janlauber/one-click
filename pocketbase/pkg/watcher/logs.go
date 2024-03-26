package watcher

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/janlauber/one-click/pkg/k8s"
	"github.com/labstack/echo/v5"
)

type RolloutLogsRequest struct {
	RolloutId string `json:"rolloutId"`
	PodName   string `json:"podName"`
}

func WsK8sRolloutLogsHandler(c echo.Context) error {
	ws, err := Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return err
	}

	defer ws.Close()

	// Wait for a single message that contains the rolloutId and podName
	_, msg, err := ws.ReadMessage()
	if err != nil {
		log.Println("WebSocket read error:", err)
		return err
	}

	// Unmarshal JSON to a struct and handle the message
	var request RolloutLogsRequest
	err = json.Unmarshal(msg, &request)
	if err != nil {
		log.Println("WebSocket unmarshal error:", err)
		ws.WriteMessage(websocket.TextMessage, []byte("{\"error\":\"Invalid request format\"}"))
		return err
	}

	go k8s.WatchK8sLogsAndSendUpdates(ws, request.RolloutId, request.PodName)

	// Keep the WebSocket connection open
	for {
		if _, _, err := ws.NextReader(); err != nil {
			ws.Close()
			break
		}
	}

	return nil
}
