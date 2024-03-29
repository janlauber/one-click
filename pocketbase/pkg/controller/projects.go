package controller

import (
	"log"

	"github.com/janlauber/one-click/pkg/k8s"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func HandleProjectDelete(e *core.RecordDeleteEvent, app *pocketbase.PocketBase) error {
	// Get project
	expr := dbx.NewExp("project = {:project}", dbx.Params{"project": e.Record.GetString("id")})
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
			err = k8s.DeleteRollout(e.Record.GetString("id"), rollout.Id)
			if err != nil {
				log.Println(err)
			}
		}
	}

	// Get autoUpdate from autoUpdates collection
	expr = dbx.NewExp("project = {:project}", dbx.Params{"project": e.Record.GetString("id")})
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

	// delete namespace in k8s
	err = k8s.DeleteNamespace(e.Record.GetString("id"))
	if err != nil {
		log.Println(err)
	}

	return nil
}
