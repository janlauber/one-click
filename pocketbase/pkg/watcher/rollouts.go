package watcher

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/janlauber/one-click/pkg/k8s"
	"github.com/labstack/echo/v5"
)

type RolloutStatusRequest struct {
	RolloutId string `json:"rolloutId"`
}

func WsK8sRolloutsHandler(c echo.Context) error {
	ws, err := Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return err
	}
	defer ws.Close()

	// Wait for a single message that contains the rolloutId
	_, msg, err := ws.ReadMessage()
	if err != nil {
		log.Println("WebSocket read error:", err)
		return err // or handle the error as appropriate
	}

	// Unmarshal JSON to a struct and handle the message
	var request RolloutStatusRequest
	err = json.Unmarshal(msg, &request)
	if err != nil {
		log.Println("WebSocket unmarshal error:", err)
		ws.WriteMessage(websocket.TextMessage, []byte("{\"error\":\"Invalid request format\"}"))
		return err // or handle the error as appropriate
	}

	go k8s.WatchK8sResourcesAndSendUpdates(ws, request.RolloutId, request.RolloutId) // Assuming "default" namespace; adjust as needed

	// Keep the WebSocket connection open
	for {
		if _, _, err := ws.NextReader(); err != nil {
			ws.Close()
			break
		}
		// Optionally, you could handle additional messages from the client here
	}

	return nil
}
