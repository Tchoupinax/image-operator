package helpers

import (
	"fmt"
	"strings"
)

func GenerateSkopeoJobName(
	name string,
	version string,
) string {
	return fmt.Sprintf(
		"skopeo-job-copy-%s-%s",
		strings.ReplaceAll(strings.ReplaceAll(name, ".", "-"), "/", "-"),
		strings.ReplaceAll(version, ".", "-"),
	)
}
