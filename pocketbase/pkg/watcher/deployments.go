package watcher

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/janlauber/one-click/pkg/k8s"
	"github.com/labstack/echo/v5"
)

type DeploymentStatusRequest struct {
	DeploymentId string `json:"deploymentId"`
	ProjectId    string `json:"projectId"`
}

func WsK8sDeploymentsHandler(c echo.Context) error {
	ws, err := Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return err
	}
	defer ws.Close()

	// Wait for a single message that contains the deploymentId
	_, msg, err := ws.ReadMessage()
	if err != nil {
		log.Println("WebSocket read error:", err)
		return err // or handle the error as appropriate
	}

	// Unmarshal JSON to a struct and handle the message
	var request DeploymentStatusRequest
	err = json.Unmarshal(msg, &request)
	if err != nil {
		log.Println("WebSocket unmarshal error:", err)
		ws.WriteMessage(websocket.TextMessage, []byte("{\"error\":\"Invalid request format\"}"))
		return err // or handle the error as appropriate
	}

	var labelSelector string

	if request.DeploymentId == "" && request.ProjectId == "" {
		log.Println("Invalid request:", request)
		ws.WriteMessage(websocket.TextMessage, []byte("{\"error\":\"Invalid request\"}"))
		return nil
	}

	if request.ProjectId != "" && request.DeploymentId == "" {
		labelSelector = fmt.Sprintf("one-click.dev/projectId=%s", request.ProjectId)
	} else if request.ProjectId != "" && request.DeploymentId != "" {
		labelSelector = fmt.Sprintf("one-click.dev/deploymentId=%s", request.DeploymentId)
	}

	go k8s.WatchK8sResourcesAndSendUpdates(ws, request.ProjectId, labelSelector) // Assuming "default" namespace; adjust as needed

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
