package controller

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
)

func HandleBlueprint(c echo.Context, app *pocketbase.PocketBase, blueprintId string) error {

	blueprint, err := app.Dao().FindRecordById("blueprints", blueprintId)
	if err != nil {
		return err
	}

	return c.JSON(200, blueprint)
}
