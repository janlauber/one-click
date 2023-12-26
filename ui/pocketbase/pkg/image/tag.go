package image

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/pocketbase/pocketbase"
)

type FilterTags struct {
	Pattern string `json:"pattern"`
	Policy  string `json:"policy"` // asc, desc, latest, semver, timestamp
}

func FilterAndSortTags(app *pocketbase.PocketBase, registry string, repository string, filterTags FilterTags) ([]string, error) {

	tags, err := fetchTags(registry, repository)
	if err != nil {
		log.Printf("Failed to fetch tags: %v", err)
		return nil, fmt.Errorf("failed to fetch tags: %v", err)
	}

	filteredTags := filterTagsByPattern(tags, filterTags.Pattern)
	// Parse the policy for sorting method and optional semver constraint
	sortMethod, semverConstraint := parsePolicy(filterTags.Policy)

	// Sort the tags based on the sorting method
	sortedTags := sortTags(filteredTags, sortMethod)

	if semverConstraint != "" {
		constraint, err := semver.NewConstraint(semverConstraint)
		if err != nil {
			log.Printf("Invalid semver constraint: %v", err)
			return nil, err
		}

		var constrainedTags []string
		for _, tag := range sortedTags {
			baseVersion, err := extractBaseVersion(tag)
			if err != nil {
				continue
			}
			version, err := semver.NewVersion(baseVersion)
			if err == nil && constraint.Check(version) {
				constrainedTags = append(constrainedTags, tag)
			}
		}
		sortedTags = constrainedTags
	}

	return sortedTags, nil
}

func extractBaseVersion(tag string) (string, error) {
	// Extracts the base version from a tag
	re := regexp.MustCompile(`^\d+\.\d+\.\d+`)
	baseVersion := re.FindString(tag)
	if baseVersion == "" {
		return "", fmt.Errorf("no base version found")
	}
	return baseVersion, nil
}

func parsePolicy(policy string) (string, string) {
	parts := strings.Split(policy, ":")
	if len(parts) == 2 && parts[0] == "semver" {
		return parts[0], parts[1] // e.g., "semver", ">1.0.0"
	}
	return policy, "" // no semver constraint
}

func fetchTags(registry string, repository string) ([]string, error) {
	url := fmt.Sprintf("%s/%s/tags?page_size=1000", registry, repository)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags: %v", err)
	}
	defer resp.Body.Close()

	var data struct {
		Results []struct {
			Name string `json:"name"`
		} `json:"results"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	tags := make([]string, len(data.Results))
	for i, result := range data.Results {
		tags[i] = result.Name
	}

	return tags, nil
}

func filterTagsByPattern(tags []string, pattern string) []string {

	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Printf("Failed to compile pattern: %v", err)
		return nil
	}

	filteredTags := make([]string, 0)
	for _, tag := range tags {
		if re.MatchString(tag) {
			filteredTags = append(filteredTags, tag)
		}
	}

	return filteredTags
}

func sortTags(tags []string, policy string) []string {
	switch policy {
	case "asc":
		sort.Strings(tags)
	case "desc":
		sort.Sort(sort.Reverse(sort.StringSlice(tags)))
	case "latest":
		sort.Slice(tags, func(i, j int) bool {
			return tags[i] == "latest"
		})
	case "semver":
		return sortSemverTags(tags)
	case "timestamp":
		return sortTimestampTags(tags)
	}
	return tags
}

func extractSemver(tag string) (*semver.Version, error) {
	// Assuming the tag is a straightforward semver string like "1.2.3"
	version, err := semver.NewVersion(tag)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func sortSemverTags(tags []string) []string {
	semverTags := make([]*semver.Version, 0, len(tags))
	tagMap := make(map[string]string)

	for _, tag := range tags {
		version, err := extractSemver(tag)
		if err == nil {
			semverTags = append(semverTags, version)
			tagMap[version.String()] = tag
		} else {
			log.Printf("Failed to parse semver for tag '%s': %v", tag, err)
		}
	}

	sort.Slice(semverTags, func(i, j int) bool {
		return semverTags[i].GreaterThan(semverTags[j])
	})

	sortedTags := make([]string, len(semverTags))
	for i, version := range semverTags {
		sortedTags[i] = tagMap[version.String()]
	}

	return sortedTags
}

func sortTimestampTags(tags []string) []string {
	timestampTags := make([]int64, 0, len(tags))
	tagMap := make(map[int64]string)

	for _, tag := range tags {
		timestamp, err := extractTimestamp(tag)
		if err == nil {
			timestampTags = append(timestampTags, timestamp)
			tagMap[timestamp] = tag
		}
	}

	sort.Slice(timestampTags, func(i, j int) bool {
		return timestampTags[i] < timestampTags[j]
	})

	sortedTags := make([]string, len(timestampTags))
	for i, timestamp := range timestampTags {
		sortedTags[i] = tagMap[timestamp]
	}

	return sortedTags
}

func extractTimestamp(tag string) (int64, error) {
	// Extract timestamp part from the tag
	// This regex matches typical timestamp patterns in the tag
	re := regexp.MustCompile(`\b\d{10}\b`)
	timestampStr := re.FindString(tag)
	if timestampStr == "" {
		return 0, fmt.Errorf("no timestamp found")
	}

	return strconv.ParseInt(timestampStr, 10, 64)
}
