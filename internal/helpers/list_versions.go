package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func ListVersion(sourceName string, matchingString string, allowCandidateRelease bool) []string {
	repository := strings.Join(strings.SplitN(sourceName, "/", 3)[1:], "/")
	url := fmt.Sprintf("https://quay.io/api/v1/repository/%s/tag?limit=100&page=2", repository)
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

	regex := GenerateRegex(matchingString, allowCandidateRelease)
	re := regexp.MustCompile(regex)

	fmt.Println(fmt.Printf("Available versions starting from %s (%s)", matchingString, regex))

	matchedTags := []string{}
	for _, tag := range result.Tags {
		if re.MatchString(tag.Name) {
			if !contains(matchedTags, tag.Name) {
				matchedTags = append(matchedTags, tag.Name)
			}
		}
	}

	return matchedTags
}

func GenerateRegex(input string, allowReleaseCandidate bool) string {
	if strings.Contains(input, ".x") {
		regex := strings.Replace(input, ".x", `.\d+`, 1)

		if allowReleaseCandidate {
			regex = fmt.Sprintf(`%s(-rc\d)?`, regex)
		}

		return fmt.Sprintf(`%s$`, regex)
	}
	return fmt.Sprintf(`%s$`, input)
}

func contains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
