package watcher

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/janlauber/one-click/pkg/k8s"
	"github.com/labstack/echo/v5"
)

type RolloutEventsRequest struct {
	RolloutId string `json:"rolloutId"`
	Kind      string `json:"kind"`
	Name      string `json:"name"`
}

func WsK8sRolloutEventsHandler(c echo.Context) error {
	conn, err := Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return err
	}
	defer conn.Close()

	// Wait for a single message that contains the rolloutId
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("WebSocket read error:", err)
		return err
	}

	// Unmarshal JSON to a struct and handle the message
	var request RolloutEventsRequest
	err = json.Unmarshal(msg, &request)
	if err != nil {
		log.Println("WebSocket unmarshal error:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("{\"error\":\"Invalid request format\"}"))
		return err
	}

	go k8s.WatchK8sEventsAndSendUpdates(conn, request.RolloutId, request.Kind, request.Name)

	// Keep the WebSocket connection open
	for {
		if _, _, err := conn.NextReader(); err != nil {
			conn.Close()
			break
		}
	}

	return nil
}
