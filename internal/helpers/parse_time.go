package helpers

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func convertFormat(input string) string {
	// Handle days
	reDays := regexp.MustCompile(`(\d+)d`)
	matchesDays := reDays.FindStringSubmatch(input)
	if matchesDays != nil {
		days, err := strconv.Atoi(matchesDays[1])
		if err != nil {
			fmt.Println(err)
		}

		hours := days * 24
		return fmt.Sprintf("%dh", hours)
	}

	// Handle weeks
	reWeeks := regexp.MustCompile(`(\d+)w`)
	matchesWeeks := reWeeks.FindStringSubmatch(input)
	if matchesWeeks != nil {
		weeks, err := strconv.Atoi(matchesWeeks[1])
		if err != nil {
			fmt.Println(err)
		}

		hours := weeks * 7 * 24
		return fmt.Sprintf("%dh", hours)
	}

	return input
}

func ParseTime(timeStr string) time.Duration {
	parsedFrequency, parsedFrequencyError := time.ParseDuration(convertFormat(timeStr))

	if parsedFrequencyError != nil {
		// By default, we return 5 minutes in case of error.
		return 5 * time.Minute
	}

	return parsedFrequency
}
