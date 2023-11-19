package controller

import (
	"log"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v5"
	"github.com/natrontech/one-click/pkg/env"
	"github.com/natrontech/one-click/pkg/k8s"
	"github.com/natrontech/one-click/pkg/models"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"gopkg.in/yaml.v2"
)

func HandleRolloutCreate(e *core.RecordCreateEvent, app *pocketbase.PocketBase) error {

	return nil
}

func HandleRolloutUpdate(e *core.RecordUpdateEvent, app *pocketbase.PocketBase) error {

	return nil
}

func HandleRolloutDelete(e *core.RecordDeleteEvent, app *pocketbase.PocketBase) error {

	// Get rollout
	rollout, err := app.Dao().FindRecordById("rollouts", e.Record.GetString("id"))
	if err != nil {
		return err
	}

	log.Println("Deleting rollout: " + rollout.GetString("id"))

	// Get project
	project, err := app.Dao().FindRecordById("projects", e.Record.GetString("projectId"))
	if err != nil {
		return err
	}

	log.Println("Deleting rollout: " + rollout.GetString("rolloutId"))
	log.Println("Deleting project: " + project.GetString("projectId"))

	// err = k8s.DeleteRollout(project.Id, rollout.Id)
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }

	return nil
}

func HandleRolloutGet(c echo.Context, app *pocketbase.PocketBase, projectId string, rolloutId string) error {

	filePath := filepath.Join(env.Config.DefaultRolloutDir, projectId, rolloutId+".yaml")
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var rolloutRevision models.Rollout
	err = yaml.Unmarshal(yamlFile, &rolloutRevision)
	if err != nil {
		return err
	}

	return c.JSON(200, rolloutRevision)
}

func HandleRolloutGetAll(c echo.Context, app *pocketbase.PocketBase, projectId string) error {

	// get all revisions of projectId
	filePath := filepath.Join(env.Config.DefaultRolloutDir, projectId)
	rolloutDir, err := os.ReadDir(filePath)
	if err != nil {
		return err
	}

	// get all rollouts of projectId
	var rollouts []string
	for _, rollout := range rolloutDir {
		// remove .yaml extension
		rollouts = append(rollouts, rollout.Name()[:len(rollout.Name())-5])
	}

	return c.JSON(200, rollouts)
}

func HandleRolloutPost(c echo.Context, app *pocketbase.PocketBase, projectId string, rollout string) error {
	// Create project dir if not exists
	filePath := filepath.Join(env.Config.DefaultRolloutDir, projectId)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.MkdirAll(filePath, 0755)
		if err != nil {
			return err
		}
	}

	// Create rollout file
	filePath = filepath.Join(filePath, rollout+".yaml")
	rolloutFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer rolloutFile.Close()

	// Get body and parse to yaml
	rolloutBody := new(models.Rollout)
	if err := c.Bind(rolloutBody); err != nil {
		return err
	}

	// Marshal struct to YAML
	yamlData, err := yaml.Marshal(rolloutBody)
	if err != nil {
		return err
	}

	// Write YAML to file
	if _, err := rolloutFile.Write(yamlData); err != nil {
		return err
	}

	// Apply the rollout
	err = k8s.CreateOrUpdateRollout(projectId, rollout)
	if err != nil {
		log.Println(err)
		return err
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
