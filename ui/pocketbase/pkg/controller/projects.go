package controller

import (
	"log"

	"github.com/natrontech/one-click/pkg/k8s"
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

		// Delete rollout in k8s
		err = k8s.DeleteRollout(e.Record.GetString("id"), rollout.Id)
		if err != nil {
			log.Println(err)
		}
	}

	// delete namespace in k8s
	err = k8s.DeleteNamespace(e.Record.GetString("id"))
	if err != nil {
		log.Println(err)
	}

	return nil
}
