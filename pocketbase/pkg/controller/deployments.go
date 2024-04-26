package controller

import (
	"log"

	"github.com/janlauber/one-click/pkg/k8s"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func HandleDeploymentDelete(e *core.RecordDeleteEvent, app *pocketbase.PocketBase) error {
	// Get deployment
	expr := dbx.NewExp("deployment = {:deployment}", dbx.Params{"deployment": e.Record.GetString("id")})
	rollouts, err := app.Dao().FindRecordsByExpr("rollouts", expr)
	if err != nil {
		return err
	}

	// Delete all rollouts in database
	for _, rollout := range rollouts {
		err = app.Dao().DeleteRecord(rollout)
		if err != nil {
			return err
		}

		// Delete rollout in k8s with no endDate
		if rollout.GetString("endDate") == "" {
			err = k8s.DeleteRollout(e.Record.GetString("project"), e.Record.GetString("id"))
			if err != nil {
				log.Println(err)
			}
		}
	}

	// Get autoUpdate from autoUpdates collection
	expr = dbx.NewExp("deployment = {:deployment}", dbx.Params{"deployment": e.Record.GetString("id")})
	autoUpdates, err := app.Dao().FindRecordsByExpr("autoUpdates", expr)
	if err != nil {
		return err
	}

	// Delete all autoUpdates in database
	for _, autoUpdate := range autoUpdates {
		err = app.Dao().DeleteRecord(autoUpdate)
		if err != nil {
			return err
		}
	}

	return nil
}

func HandleDeploymentStatus(c echo.Context, app *pocketbase.PocketBase, projectId string, deploymentId string) error {
	// We get the status of the kind: Rollout with the name: deploymentId (because we don't want to recreate the object on each new rollout)
	status, err := k8s.GetRolloutStatus(projectId, deploymentId)
	if err != nil {
		return err
	}

	return c.JSON(200, status)
}

func HandleDeploymentMetrics(c echo.Context, app *pocketbase.PocketBase, projectId string, deploymentId string) error {
	// We get the metrics of the kind: Rollout with the name: deploymentId (because we don't want to recreate the object on each new rollout)
	metrics, err := k8s.GetRolloutMetrics(projectId, deploymentId)
	if err != nil {
		return err
	}

	return c.JSON(200, metrics)
}

func HandleDeploymentEvents(c echo.Context, app *pocketbase.PocketBase, projectId string, deploymentId string) error {
	// We get the events of the kind: Rollout with the name: deploymentId (because we don't want to recreate the object on each new rollout)
	events, err := k8s.GetRolloutEvents(projectId, deploymentId)
	if err != nil {
		return err
	}

	return c.JSON(200, events)
}
