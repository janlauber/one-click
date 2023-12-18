package controller

import (
	"log"

	"github.com/labstack/echo/v5"
	"github.com/natrontech/one-click/pkg/image"
	"github.com/pocketbase/pocketbase"
)

func HandleAutoUpdate(c echo.Context, app *pocketbase.PocketBase) error {
	// call image.AutoUpdate
	err := image.AutoUpdate(app)
	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(200, "ok")
}
