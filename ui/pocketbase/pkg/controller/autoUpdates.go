package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/natrontech/one-click/pkg/image"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

func HandleAutoUpdateCreate(e *core.RecordCreateEvent, app *pocketbase.PocketBase) error {

	// get current running rollout -> no endDate set

	expr := dbx.NewExp("project = {:project}", dbx.Params{"project": e.Record.GetString("project")})
	rollouts, err := app.Dao().FindRecordsByExpr("rollouts", expr)
	if err != nil {
		return err
	}

	// search rollouts for the one with no endDate
	var runningRollout *models.Record
	for _, rollout := range rollouts {
		if rollout.GetString("endDate") == "" {
			runningRollout = rollout
			break
		}
	}

	// get the registry, repository, tag from the rollout manifest under spec.image

	manifest := runningRollout.GetString("manifest")
	// parse the manifest to a traversable object
	var manifestObj map[string]interface{}
	err = json.Unmarshal([]byte(manifest), &manifestObj)
	if err != nil {
		return err
	}

	// get the spec.image.registry, spec.image.repository, spec.image.tag
	registry := manifestObj["spec"].(map[string]interface{})["image"].(map[string]interface{})["registry"].(string)
	repository := manifestObj["spec"].(map[string]interface{})["image"].(map[string]interface{})["repository"].(string)
	tag := manifestObj["spec"].(map[string]interface{})["image"].(map[string]interface{})["tag"].(string)

	fmt.Println("registry", registry)
	fmt.Println("repository", repository)
	fmt.Println("tag", tag)

	// get the interval, pattern and policy from the autoUpdate
	interval := e.Record.GetString("interval")
	pattern := e.Record.GetString("pattern")
	policy := e.Record.GetString("policy")

	fmt.Println("interval", interval)
	fmt.Println("pattern", pattern)
	fmt.Println("policy", policy)

	// start a go routine with the interval to check for new tags
	filterTags := image.FilterTags{
		Pattern: pattern,
		Policy:  policy,
	}

	if strings.Contains(registry, "docker") {
		registry = "https://registry.hub.docker.com/v2/repositories"
		// if there is no / in the repository, add library/
		if !strings.Contains(repository, "/") {
			repository = "library/" + repository
		}
	}

	// ghcr
	if strings.Contains(registry, "ghcr") {
		registry = "https://ghcr.io/v2"
	}

	// parse interval to time.Duration
	// get if its 1m, 1h, 1d, 1w, or 1m30s, 1h30m, 1d12h, 1w3d, 1m30s30ms
	intervalDuration, err := ParseDuration(interval)
	if err != nil {
		return err
	}

	go func() {
		for {
			// check for new tags
			tags, err := image.FilterAndSortTags(app, registry, repository, filterTags)
			if err != nil {
				log.Println(err)
				return
			}

			// if no tags are found, return an empty response
			var latestTag string
			if len(tags) > 0 {
				latestTag = tags[0]
			} else {
				tags = []string{}
			}

			fmt.Println("latestTag", latestTag)

			// sleep for the interval
			time.Sleep(intervalDuration)
		}
	}()

	// if new tag found, create a new rollout with the new tag and the same manifest

	return nil
}

func HandleAutoUpdateUpdate(e *core.RecordUpdateEvent, app *pocketbase.PocketBase) error {
	return nil
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
