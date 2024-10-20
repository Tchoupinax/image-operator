package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

const (
	// Quay.io URL for the Argo CD image
	quayURL = "https://quay.io/api/v1/repository/argoproj/argocd/tag?limit=100"
)

func listVersion() {
	// Make a request to Quay.io API
	resp, err := http.Get(quayURL)
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

	// Regular expression to match version pattern 2.12.x
	re := regexp.MustCompile(`^v2\.13\.\d+(-rc\d)?$`)

	// Print available versions
	fmt.Println("Available versions starting from 2.13.x:")
	for _, tag := range result.Tags {
		if re.MatchString(tag.Name) {
			fmt.Println(tag.Name)
		}
	}
}
