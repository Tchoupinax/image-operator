package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/go-logr/logr"
)

func ListVersion(
	logger logr.Logger,
	sourceName string,
	matchingString string,
	allowCandidateRelease bool,
) []string {
	logger.Info(fmt.Sprintf("Looking for version for %s:%s", sourceName, matchingString))

	repository := strings.Join(strings.SplitN(sourceName, "/", 3)[1:], "/")
	matchedTags := []string{}
	page := 1

	for {
		// Update the URL to fetch the current page
		url := fmt.Sprintf("https://quay.io/api/v1/repository/%s/tag/?limit=100&page=%d", repository, page)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error fetching data: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error: received status code %d\n", resp.StatusCode)
			os.Exit(1)
		}

		// Decode the JSON response
		var result struct {
			Tags []struct {
				Name string `json:"name"`
			} `json:"tags"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Printf("Error decoding JSON: %v\n", err)
			os.Exit(1)
		}

		if len(result.Tags) == 0 {
			break
		}

		// Generate the regex
		regex := GenerateRegex(matchingString, allowCandidateRelease)
		re := regexp.MustCompile(regex)

		for _, tag := range result.Tags {
			if re.MatchString(tag.Name) && !contains(matchedTags, tag.Name) {
				matchedTags = append(matchedTags, tag.Name)
			}
		}

		page++
	}

	return matchedTags
}

func GenerateRegex(input string, allowReleaseCandidate bool) string {
	if strings.Contains(input, ".x") {
		regex := strings.Replace(input, ".x", `.\d+`, 1)

		if allowReleaseCandidate {
			regex = fmt.Sprintf(`%s(-rc\d)?`, regex)
		}

		return fmt.Sprintf(`^v?%s$`, regex)
	}

	return fmt.Sprintf(`^v?%s$`, input)
}

func contains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
