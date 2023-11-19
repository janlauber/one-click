package controller

import (
	"log"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v5"
	"github.com/natrontech/one-click/pkg/env"
	"github.com/natrontech/one-click/pkg/models"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"gopkg.in/yaml.v2"
)

func HandleRolloutCreate(e *core.RecordCreateEvent, app *pocketbase.PocketBase) error {

	log.Println("HandleRolloutCreate")

	return nil
}

func HandleRolloutUpdate(e *core.RecordUpdateEvent, app *pocketbase.PocketBase) error {

	log.Println("HandleRolloutUpdate")

	return nil
}

func HandleRolloutDelete(e *core.RecordDeleteEvent, app *pocketbase.PocketBase) error {

	log.Println("HandleRolloutDelete")

	return nil
}

func HandleRolloutGet(c echo.Context, app *pocketbase.PocketBase, project string, revision string) error {

	filePath := filepath.Join(env.Config.DefaultRolloutDir, project, revision+".yaml")
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var frameworkRevision models.Framework
	err = yaml.Unmarshal(yamlFile, &frameworkRevision)
	if err != nil {
		return err
	}

	return c.JSON(200, frameworkRevision)
}

func HandleRolloutGetAll(c echo.Context, app *pocketbase.PocketBase, project string) error {

	// get all revisions of project
	filePath := filepath.Join(env.Config.DefaultRolloutDir, project)
	rolloutDir, err := os.ReadDir(filePath)
	if err != nil {
		return err
	}

	// get all revisions of project
	var revisions []string
	for _, rollout := range rolloutDir {
		// remove .yaml extension
		revisions = append(revisions, rollout.Name()[:len(rollout.Name())-5])
	}

	return c.JSON(200, revisions)
}

func HandleRolloutPost(c echo.Context, app *pocketbase.PocketBase, project string, revision string) error {
	// Create project dir if not exists
	filePath := filepath.Join(env.Config.DefaultRolloutDir, project)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.MkdirAll(filePath, 0755)
		if err != nil {
			return err
		}
	}

	// Create rollout file
	filePath = filepath.Join(filePath, revision+".yaml")
	rolloutFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer rolloutFile.Close()

	// Get body and parse to yaml
	framework := new(models.Framework)
	if err := c.Bind(framework); err != nil {
		return err
	}

	// Marshal struct to YAML
	yamlData, err := yaml.Marshal(framework)
	if err != nil {
		return err
	}

	// Write YAML to file
	if _, err := rolloutFile.Write(yamlData); err != nil {
		return err
	}

	return nil
}
