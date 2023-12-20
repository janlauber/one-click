package controller

import (
	"io"
	"log"
	"strings"

	"encoding/json"

	"github.com/labstack/echo/v5"
	"github.com/natrontech/one-click/pkg/image"
	"github.com/pocketbase/pocketbase"
)

type UpdateRequest struct {
	Pattern    string `json:"pattern"`
	Policy     string `json:"policy"`
	Registry   string `json:"registry"`
	Repository string `json:"repository"`
}

func HandleAutoUpdate(c echo.Context, app *pocketbase.PocketBase) error {
	// call image.AutoUpdate

	// read the request body
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Println(err)
		return err
	}

	// parse the request body
	var updateReq UpdateRequest
	err = json.Unmarshal(body, &updateReq)
	if err != nil {
		log.Println(err)
		return err
	}

	// get the filter tags (pattern, extract) from the request body
	pattern := updateReq.Pattern // required
	if pattern == "" {
		pattern = ".*"
	}

	// check if the pattern contains spaces and remove them
	pattern = strings.ReplaceAll(pattern, " ", "")

	// get the policy (semver (range), alphabetical (asc, desc), latest) from the request body
	policy := updateReq.Policy // required
	if policy == "" {
		policy = "semver"
	}

	// get the registry from the request body
	registry := updateReq.Registry // required
	if registry == "" {
		registry = "https://registry.hub.docker.com/v2/repositories"
	}

	// get the repository from the request body
	repository := updateReq.Repository // required
	if repository == "" {
		repository = "janlauber/test-image"
	}

	filterTags := image.FilterTags{
		Pattern: pattern,
		Policy:  policy,
	}

	tags, err := image.FilterAndSortTags(app, registry, repository, filterTags)
	if err != nil {
		log.Println(err)
		return err
	}

	// if no tags are found, return an empty response
	var latestTag string
	if len(tags) > 0 {
		latestTag = tags[0]
	} else {
		tags = []string{}
	}

	return c.JSON(200, map[string]interface{}{
		"tags":      tags,
		"latestTag": latestTag,
	})
}
