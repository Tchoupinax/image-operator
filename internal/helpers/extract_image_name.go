package helpers

import (
	"fmt"
	"regexp"
)

type ImageDetails struct {
	Registry string
	Image    string
	Version  string
}

func ExtractImageName(input string) (*ImageDetails, error) {
	pattern := `Failed to pull image "(?P<registry>[a-zA-Z0-9._-]+(?:\.[a-zA-Z0-9._-]+)+(?:/[a-zA-Z0-9._-]+)*)/(?P<image>[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+):(?P<version>[a-zA-Z0-9_.-]+)"`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(input)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no matches found in the input string")
	}

	result := &ImageDetails{
		Registry: matches[1],
		Image:    matches[2],
		Version:  matches[3],
	}

	return result, nil
}
