package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
)

type DockerHubAuth struct {
	username string
	password string
}

type AWSPublicECR struct {
	token string
}

func ListVersions(
	logger logr.Logger,
	sourceName string,
	matchingString string,
	allowCandidateRelease bool,
	dockerhubAuth DockerHubAuth,
	awsPublicECR AWSPublicECR,
) []string {
	logger.Info(fmt.Sprintf("Looking for versions for %s:%s", sourceName, matchingString))

	isQuay := strings.HasPrefix(sourceName, "quay.io/")
	isAWSPublicECR := strings.HasPrefix(sourceName, "public.ecr.aws/")

	repoParts := strings.SplitN(sourceName, "/", 2)
	if len(repoParts) != 2 {
		repoParts = strings.SplitN("library/"+sourceName, "/", 2)
	}
	repository := repoParts[1]

	matchedTags := []string{}
	page := 1
	var nextToken string

	for {
		var url string
		var req *http.Request

		if isQuay {
			url = fmt.Sprintf("https://quay.io/api/v1/repository/%s/tag/?limit=100&page=%d", repository, page)
			req, _ = http.NewRequest("GET", url, nil)
		} else if isAWSPublicECR {
			url = "https://api.us-east-1.gallery.ecr.aws/describeImageTags"

			var parts = strings.Split(repository, "/")

			awsPageCount := 1000
			if os.Getenv("API_AWS_PAGE_ITEMS_COUNT") != "" {
				number, err := strconv.Atoi(os.Getenv("API_AWS_PAGE_ITEMS_COUNT"))
				if err == nil {
					awsPageCount = number
				}
			}

			var jsonData string
			if nextToken != "" {
				jsonData = fmt.Sprintf(`{
					"registryAliasName":"%s",
					"repositoryName":"%s",
					"nextToken": "%s",
					"maxResults": %d
				}`, parts[0], parts[1]+"/"+parts[2], nextToken, awsPageCount)
			} else {
				jsonData = fmt.Sprintf(`{
					"registryAliasName":"%s",
					"repositoryName":"%s",
					"maxResults": %d
				}`, parts[0], parts[1]+"/"+parts[2], awsPageCount)
			}

			req, _ = http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
			req.Header.Set("Cache-Control", "no-cache")
			req.Header.Set("TE", "trailers")
			req.Header.Set("Content-Type", "application/json")
		} else {
			if !strings.Contains(sourceName, "/") {
				sourceName = "library/" + sourceName
			}

			url = fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags/?page_size=100&page=%d", sourceName, page)
			req, _ = http.NewRequest("GET", url, nil)
		}

		if !isQuay && !isAWSPublicECR && dockerhubAuth.username != "" && dockerhubAuth.password != "" {
			req.SetBasicAuth(dockerhubAuth.username, dockerhubAuth.password)
		}

		if isAWSPublicECR && awsPublicECR.token != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Basic %s", awsPublicECR.token))
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("Error fetching data: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error: received status code %d\n", resp.StatusCode)
			fmt.Println(resp.Body)
			break
		}

		var result struct {
			Tags []struct {
				Name string `json:"name"`
			} `json:"tags"`
			Results []struct {
				Name string `json:"name"`
			} `json:"results"`
			ImageTagDetails []struct {
				ImageTag string `json:"imageTag"`
			} `json:"imageTagDetails"`
			NextToken string `json:"nextToken"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Printf("Error decoding JSON: %v\n", err)
			os.Exit(1)
		}

		if isAWSPublicECR {
			nextToken = result.NextToken
		}

		var tags []string
		if isQuay {
			if page > 10 || len(result.Tags) == 0 {
				break
			}
			for _, tag := range result.Tags {
				tags = append(tags, tag.Name)
			}
		} else if isAWSPublicECR {
			pageMax := 4
			if os.Getenv("API_AWS_PAGE_MAX") != "" {
				number, err := strconv.Atoi(os.Getenv("API_AWS_PAGE_MAX"))
				if err == nil {
					pageMax = number
				}
			}

			// Are made 1000 by 1000
			if page > pageMax || len(result.ImageTagDetails) == 0 {
				break
			}

			for _, image := range result.ImageTagDetails {
				tags = append(tags, image.ImageTag)
			}
		} else {
			if page > 10 || len(result.Results) == 0 {
				break
			}
			for _, result := range result.Results {
				tags = append(tags, result.Name)
			}
		}

		regex := GenerateRegex(matchingString, allowCandidateRelease)
		re := regexp.MustCompile(regex)

		for _, tag := range tags {
			if re.MatchString(tag) && !contains(matchedTags, tag) {
				matchedTags = append(matchedTags, tag)
			}
		}

		page++
	}

	logger.Info(fmt.Sprintf("%d images detected", len(matchedTags)))
	sort.Strings(matchedTags)

	return matchedTags
}

func GenerateRegex(input string, allowReleaseCandidate bool) string {
	if strings.Contains(input, ".x") {
		regex := strings.ReplaceAll(input, ".x", `.\d+`)

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
