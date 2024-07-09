package watcher

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

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
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not open WebSocket connection")
	}
	defer ws.Close()

	_, msg, err := ws.ReadMessage()
	if err != nil {
		log.Println("WebSocket read error:", err)
		ws.WriteMessage(websocket.TextMessage, []byte("{\"error\":\"Read error\"}")) // Better error transmission
		return err
	}

	var request RolloutLogsRequest
	if err = json.Unmarshal(msg, &request); err != nil {
		log.Println("WebSocket unmarshal error:", err)
		ws.WriteMessage(websocket.TextMessage, []byte("{\"error\":\"Unmarshal error\"}")) // Better error transmission
		return err
	}

	ctx, cancel := context.WithCancel(context.Background()) // Create a new context that we can cancel
	defer cancel()                                          // Ensure resources are freed on handler exit

	// Offload log watching to a goroutine, passing context to manage its lifecycle
	go func() {
		defer cancel() // Ensure calling cancel when the goroutine finishes
		k8s.WatchK8sLogsAndSendUpdates(ws, request.RolloutId, request.PodName, ctx)
	}()

	// Await for close event in the main function to properly manage connection lifecycle
	_, _, err = ws.ReadMessage() // This blocks until client disconnects or errors out
	if err != nil {
	}

	return nil
}
