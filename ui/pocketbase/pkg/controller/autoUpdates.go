package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/janlauber/one-click/pkg/image"
	"github.com/janlauber/one-click/pkg/k8s"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

// New for testing
// I want a controller which goes through all the autoUpdates and checks if there is a new tag
func AutoUpdateController(app *pocketbase.PocketBase) error {

	autoUpdates, err := app.Dao().FindRecordsByExpr("autoUpdates")
	if err != nil {
		log.Printf("Error getting autoUpdates: %v\n", err)
		return err
	}

	// Time layout constant
	const layout = "2006-01-02 15:04:05.000Z"

	var wg sync.WaitGroup
	for _, autoUpdate := range autoUpdates {
		wg.Add(1)
		go func(autoUpdate *models.Record) {
			defer wg.Done()

			creationTime, err := time.Parse(layout, autoUpdate.GetString("created"))
			if err != nil {
				log.Printf("Error parsing creation time: %v\n", err)
				return
			}

			interval, err := ParseDuration(autoUpdate.GetString("interval"))
			if err != nil {
				log.Printf("Error parsing interval: %v\n", err)
				return
			}

			timeSinceCreation := time.Since(creationTime)
			minutesSinceCreation := int(timeSinceCreation.Minutes())

			if minutesSinceCreation%int(interval.Minutes()) != 0 {
				return
			}

			err = UpdateImage(autoUpdate, app)
			if err != nil {
				log.Printf("Error updating image: %v\n", err)
				return
			}

		}(autoUpdate)
	}

	wg.Wait()
	return nil
}

func constructRegistryURL(registry string) (string, error) {
	switch {
	case strings.Contains(registry, "docker"):
		return "docker.io", nil
	case strings.Contains(registry, "ghcr"):
		return "ghcr.io", nil
	default:
		return "", fmt.Errorf("unsupported registry: %s", registry)
	}
}

func ParseDuration(duration string) (time.Duration, error) {
	unitMap := map[string]time.Duration{
		"ms": time.Millisecond,
		"s":  time.Second,
		"m":  time.Minute,
		"h":  time.Hour,
		"d":  time.Hour * 24,
		"w":  time.Hour * 24 * 7,
	}

	duration = strings.TrimSpace(duration)
	unit := duration[len(duration)-1:]
	valueStr := duration[:len(duration)-1]

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, err
	}

	unitDuration, ok := unitMap[unit]
	if !ok {
		return 0, errors.New("invalid duration unit")
	}

	return time.Duration(value) * unitDuration, nil
}

func UpdateImage(autoUpdate *models.Record, app *pocketbase.PocketBase) error {
	expr := dbx.NewExp("project = {:project}", dbx.Params{"project": autoUpdate.GetString("project")})
	rollouts, err := app.Dao().FindRecordsByExpr("rollouts", expr)
	if err != nil {
		log.Printf("Error querying rollouts: %v\n", err)
		return err
	}

	var runningRollout *models.Record
	for _, rollout := range rollouts {
		if rollout.GetString("endDate") == "" {
			runningRollout = rollout
			break
		}
	}

	if runningRollout == nil {
		return err
	}

	manifestObj := make(map[string]interface{})
	if err := json.Unmarshal([]byte(runningRollout.GetString("manifest")), &manifestObj); err != nil {
		log.Printf("Error unmarshaling manifest: %v\n", err)
		return err
	}

	specImage := manifestObj["spec"].(map[string]interface{})["image"].(map[string]interface{})
	registryURL, err := constructRegistryURL(specImage["registry"].(string))
	if err != nil {
		log.Printf("Error constructing registry URL: %v\n", err)
		return err
	}

	repository := specImage["repository"].(string)
	filterTags := image.FilterTags{
		Pattern: autoUpdate.GetString("pattern"),
		Policy:  autoUpdate.GetString("policy"),
	}

	// if repository doesn't contain a slash and the registry is docker hub, add library/ to the repository
	if !strings.Contains(repository, "/") && strings.Contains(registryURL, "docker") {
		repository = "library/" + repository
	}

	// check if username and password are set on the manifest
	username := ""
	password := ""
	if specImage["username"] != nil && specImage["password"] != nil {
		username = specImage["username"].(string)
		password = specImage["password"].(string)
	}

	tags, err := image.FilterAndSortTags(app, registryURL, repository, filterTags, username, password)
	if err != nil {
		log.Printf("Error filtering and sorting tags: %v\n", err)
		return err
	}

	if len(tags) > 0 {

		// if the latest tag is not the same as the current tag, update the manifest
		if tags[0] != specImage["tag"].(string) {
			specImage["tag"] = tags[0]
			manifest, err := json.Marshal(manifestObj)
			if err != nil {
				log.Printf("Error marshaling manifest: %v\n", err)
				return err
			}

			// create a new record for the rollout
			collection, err := app.Dao().FindCollectionByNameOrId("rollouts")
			if err != nil {
				log.Printf("Error finding collection: %v\n", err)
				return err
			}

			// Set the endDate on the running rollout
			runningRollout.Set("endDate", time.Now().UTC().Format(time.RFC3339))
			err = app.Dao().SaveRecord(runningRollout)
			if err != nil {
				log.Printf("Error saving rollout: %v\n", err)
				return err
			}

			// delete the rollout in k8s
			err = k8s.DeleteRollout(runningRollout.GetString("project"), runningRollout.GetString("id"))
			if err != nil {
				log.Printf("Error deleting rollout: %v\n", err)
				return err
			}

			// create a new record for the rollout
			record := models.NewRecord(collection)
			// set the fields
			record.Set("user", runningRollout.GetString("user"))
			record.Set("project", runningRollout.GetString("project"))
			record.Set("manifest", string(manifest))
			record.Set("startDate", time.Now().UTC().Format(time.RFC3339))

			// save the record
			err = app.Dao().SaveRecord(record)
			if err != nil {
				log.Printf("Error saving record: %v\n", err)
				return err
			}

			user, err := app.Dao().FindRecordById("users", record.GetString("user"))
			if err != nil {
				log.Printf("Error finding user: %v\n", err)
				return err
			}

			// create the rollout in k8s
			err = k8s.CreateOrUpdateRollout(record.Id, user, record.GetString("project"), record.GetString("manifest"))
			if err != nil {
				log.Printf("Error creating or updating rollout: %v\n", err)
				return err
			}
		}
	}

	return nil
}
