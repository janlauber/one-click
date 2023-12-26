package controller

import (
	"log"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/natrontech/one-click/pkg/k8s"
	"github.com/natrontech/one-click/pkg/util"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func HandleRolloutCreate(e *core.RecordCreateEvent, app *pocketbase.PocketBase) error {

	// Get project
	project, err := app.Dao().FindRecordById("projects", e.Record.GetString("project"))
	if err != nil {
		return err
	}

	// Get user
	user, err := app.Dao().FindRecordById("users", project.GetString("user"))
	if err != nil {
		return err
	}

	// Generate a rolloutId (15 characters)
	rolloutId := util.GenerateId(15)

	// Check if endDate is set, if yes, throw error
	if e.Record.GetString("endDate") != "" {
		return echo.NewHTTPError(400, "endDate is not allowed on create")
	}

	// Check if there is another rollout in the same project with no endDate
	running_rollout, err := app.Dao().FindFirstRecordByFilter("rollouts", "endDate = '' && project = {:project}",
		dbx.Params{"project": e.Record.GetString("project")},
	)
	if err != nil {
		if contains := strings.Contains(err.Error(), "no rows"); !contains {
			return err
		}
	}

	// if there is another rollout in the same project with no endDate, set endDate to now on that rollout
	if running_rollout != nil {
		running_rollout.Set("endDate", time.Now().UTC().Format(time.RFC3339))
		err = app.Dao().SaveRecord(running_rollout)
		if err != nil {
			return err
		}

		// Delete rollout in k8s
		err = k8s.DeleteRollout(project.Id, running_rollout.Id)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	// Create rollout in k8s
	err = k8s.CreateOrUpdateRollout(rolloutId, user, project.Id, e.Record.GetString("manifest"))
	if err != nil {
		log.Println(err)
		return err
	}

	// Update rollout with startDate time.Now().UTC().Format(time.RFC3339)
	e.Record.Set("startDate", time.Now().UTC().Format(time.RFC3339))

	// Update rollout with rolloutId
	e.Record.Set("id", rolloutId)

	return nil

}

func HandleRolloutUpdate(e *core.RecordUpdateEvent, app *pocketbase.PocketBase) error {

	// Get rollout
	rollout, err := app.Dao().FindRecordById("rollouts", e.Record.GetString("id"))
	if err != nil {
		return err
	}

	// Get project
	project, err := app.Dao().FindRecordById("projects", rollout.GetString("project"))
	if err != nil {
		return err
	}

	// Get user
	user, err := app.Dao().FindRecordById("users", project.GetString("user"))
	if err != nil {
		return err
	}

	// Check if endDate is set, if yes, delete rollout
	if e.Record.GetString("endDate") != "" {
		err = k8s.DeleteRollout(project.Id, rollout.Id)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	} else if rollout.GetString("endDate") != "" {

		// Check if there is another rollout in the same project with no endDate
		running_rollout, err := app.Dao().FindFirstRecordByFilter("rollouts", "endDate = '' && project = {:project}",
			dbx.Params{"project": rollout.GetString("project")},
		)
		if err != nil {
			// only throw error string if it doesn't contain "no rows"
			if contains := strings.Contains(err.Error(), "no rows"); !contains {
				return err
			}
		}

		// if there is another rollout in the same project with no endDate, set endDate to now on that rollout
		if running_rollout != nil && running_rollout.Id != rollout.Id {
			running_rollout.Set("endDate", time.Now().UTC().Format(time.RFC3339))
			err = app.Dao().SaveRecord(running_rollout)
			if err != nil {
				return err
			}

			// Delete rollout in k8s
			err = k8s.DeleteRollout(project.Id, running_rollout.Id)
			if err != nil {
				log.Println(err)
				return err
			}
		}

		e.Record.Set("startDate", time.Now().UTC().Format(time.RFC3339))

		// If endDate was set before, but is not set anymore, create rollout again
		err = k8s.CreateOrUpdateRollout(rollout.Id, user, project.Id, e.Record.GetString("manifest"))
		if err != nil {
			log.Println(err)
			return err
		}

		return nil

	} else {
		e.Record.Set("startDate", time.Now().UTC().Format(time.RFC3339))
		// Update rollout in k8s
		err = k8s.CreateOrUpdateRollout(rollout.Id, user, project.Id, e.Record.GetString("manifest"))
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func HandleRolloutDelete(e *core.RecordDeleteEvent, app *pocketbase.PocketBase) error {

	// Get rollout
	rollout, err := app.Dao().FindRecordById("rollouts", e.Record.GetString("id"))
	if err != nil {
		return err
	}

	// Get project
	project, err := app.Dao().FindRecordById("projects", rollout.GetString("project"))
	if err != nil {
		return err
	}

	// Check if endDate is set, if no, delete rollout
	if rollout.GetString("endDate") == "" {
		err = k8s.DeleteRollout(project.Id, rollout.Id)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func HandleRolloutStatus(c echo.Context, app *pocketbase.PocketBase, projectId string, rollout string) error {

	// Get rollout status
	status, err := k8s.GetRolloutStatus(projectId, rollout)
	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(200, status)
}

func HandleRolloutMetrics(c echo.Context, app *pocketbase.PocketBase, projectId string, rollout string) error {

	// Get rollout metrics
	metrics, err := k8s.GetRolloutMetrics(projectId, rollout)
	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(200, metrics)
}

func HandleRolloutEvents(c echo.Context, app *pocketbase.PocketBase, projectId string, rollout string) error {

	// Get rollout events
	events, err := k8s.GetRolloutEvents(projectId, rollout)
	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(200, events)
}
