package helpers

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func converFormat(input string) string {
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

func ParseTime(timeStr string) (time.Duration, error) {
	parsedFrequency, parsedFrequencyError := time.ParseDuration(converFormat(timeStr))

	if parsedFrequencyError != nil {
		return time.Duration(time.Now().Day()), parsedFrequencyError
	}

	return parsedFrequency, nil
}
