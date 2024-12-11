package helpers

import (
	"fmt"
	"strings"
)

func GenerateSkopeoJobName(
	name string,
	version string,
) string {
	return truncateString(fmt.Sprintf(
		"image-operator-copy-%s-%s",
		strings.ReplaceAll(strings.ReplaceAll(name, ".", "-"), "/", "-"),
		strings.ReplaceAll(version, ".", "-"),
	), 63)
}

func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen]
	}
	return s
}
