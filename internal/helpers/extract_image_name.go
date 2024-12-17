package helpers

import (
	"fmt"
	"regexp"
	"strings"
)

var dockerImageWithoutPrefix = []string{
	"alpine",
	"busybox",
	"debian",
	"ubuntu",
}

type ImageDetails struct {
	Registry string
	Image    string
	Version  string
}

func ExtractImageName(input string) (*ImageDetails, error) {
	pattern := `([a-z0-9-._]*/[a-z0-9-]*/[a-z0-9-]*(/[a-z0-9-]*)?(:[a-z0-9]*)?)`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(input)
	fmt.Println(matches)
	if len(matches) < 1 {
		return nil, fmt.Errorf("no matches found in the input string")
	}

	var imageName = matches[0]
	// Handle the case where there is no tag
	// We add latest explicitly
	if !strings.Contains(imageName, ":") {
		imageName += ":latest"
	}

	var registry string
	var image string
	parts := strings.Split(imageName, "/")

	if len(parts) == 3 {
		if Contains(dockerImageWithoutPrefix, strings.Split(parts[2], ":")[0]) {
			registry = parts[0] + "/" + parts[1]
			image = parts[2]
		} else {
			registry = parts[0]
			image = parts[1] + "/" + parts[2]
		}
	} else if len(parts) == 4 {
		registry = parts[0] + "/" + parts[1]
		image = parts[2] + "/" + parts[3]
	}

	result := &ImageDetails{
		Registry: registry,
		Image:    strings.Split(image, ":")[0],
		Version:  strings.Split(image, ":")[1],
	}

	return result, nil
}
