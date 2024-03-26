package controller

import (
	"net/http"

	"github.com/janlauber/one-click/pkg/k8s"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
)

func HandleClusterInfo(c echo.Context, app *pocketbase.PocketBase) error {
	clusterInfo, err := k8s.GetClusterInfo()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, clusterInfo)
}
