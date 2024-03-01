package controller

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
)

func HandleBlueprint(c echo.Context, app *pocketbase.PocketBase, blueprintId string) error {

	blueprint, err := app.Dao().FindRecordById("blueprints", blueprintId)
	if err != nil {
		return err
	}

	userRecord, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
	if userRecord == nil {
		return c.JSON(401, "Unauthorized")
	}

	return c.JSON(200, blueprint)
}

func HandleBlueprintAdd(c echo.Context, app *pocketbase.PocketBase, blueprintId string) error {
	blueprint, err := app.Dao().FindRecordById("blueprints", blueprintId)
	if err != nil {
		return err
	}

	userRecord, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
	if userRecord == nil {
		return c.JSON(401, "Unauthorized")
	}

	// Add the user to the blueprint under the "users" field
	currentUsers := blueprint.GetStringSlice("users")

	// Check if the user is already "owner" of the blueprint
	if blueprint.GetString("owner") == userRecord.Id {
		return c.JSON(200, blueprint)
	}

	// Check if the user is already in the list
	for _, u := range currentUsers {
		if u == userRecord.Id {
			return c.JSON(200, blueprint)
		}
	}

	currentUsers = append(currentUsers, userRecord.Id)
	blueprint.Set("users", currentUsers)

	if err := app.Dao().SaveRecord(blueprint); err != nil {
		return err
	}

	return c.JSON(200, blueprint)
}
