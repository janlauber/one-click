package controller

import (
	"log"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
)

type UpdateRequest struct {
	Pattern    string `json:"pattern"`
	Policy     string `json:"policy"`
	Registry   string `json:"registry"`
	Repository string `json:"repository"`
}

func HandleAutoUpdate(c echo.Context, app *pocketbase.PocketBase, autoUpdateId string) error {
	autoUpdate, err := app.Dao().FindRecordById("autoUpdates", autoUpdateId)
	if err != nil {
		log.Printf("Error getting autoUpdate: %v\n", err)
		return err
	}

	err = UpdateImage(autoUpdate, app)
	if err != nil {
		log.Printf("Error updating image: %v\n", err)
		return err
	}

	return c.JSON(200, autoUpdate)
}
