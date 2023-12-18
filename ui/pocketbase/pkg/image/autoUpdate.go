package image

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"

	"github.com/Masterminds/semver"
	"github.com/pocketbase/pocketbase"
)

type registryResponse struct {
	Tags []string `json:"tags"`
}

func AutoUpdate(app *pocketbase.PocketBase, regex string) error {
	registry := "https://registry.hub.docker.com/v2/repositories"
	repository := "library/nginx" // Example: for official nginx image

	// Fetch tags from the registry
	tags, err := fetchTags(registry, repository)
	if err != nil {
		return err
	}

	// Filter and sort tags
	validTags := make([]string, 0)
	for _, tag := range tags {
		if validateTagWithRegex(tag, regex) && validateTagWithSemver(tag) {
			validTags = append(validTags, tag)
		}
	}

	// Sort by SemVer
	sortedTags := sortTagsBySemver(validTags)

	// Implement the logic to handle the sorted tags

	fmt.Println(sortedTags)

	return nil
}

func fetchTags(registry string, repository string) ([]string, error) {
	resp, err := http.Get(registry + "/" + repository + "/tags/list")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch tags")
	}

	var r registryResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	return r.Tags, nil
}

func validateTagWithRegex(tag string, regex string) bool {
	matched, _ := regexp.MatchString(regex, tag)
	return matched
}

func validateTagWithSemver(tag string) bool {
	_, err := semver.NewVersion(tag)
	return err == nil
}

func sortTagsBySemver(tags []string) []string {
	vs := make([]*semver.Version, len(tags))
	for i, tag := range tags {
		v, _ := semver.NewVersion(tag)
		vs[i] = v
	}

	sort.Sort(semver.Collection(vs))

	sortedTags := make([]string, len(tags))
	for i, v := range vs {
		sortedTags[i] = v.Original()
	}
	return sortedTags
}
