package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
)

type AuthResponse struct {
	Token string `json:"token"`
}

type DockerHubQuota struct {
	Ip            string
	Limit         int
	LimitWait     int
	Remaining     int
	RemainingWait int
	Succeeded     bool
}

func GetDockerhubLimit(setupLog logr.Logger) DockerHubQuota {
	url := "https://auth.docker.io/token?service=registry.docker.io&scope=repository:ratelimitpreview/test:pull"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching token: %s\n", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s\n", err)
	}

	var authResp AuthResponse
	err = json.Unmarshal(body, &authResp)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %s\n", err)
	}

	return callWithToken(authResp.Token)
}

func callWithToken(token string) DockerHubQuota {
	var url = "https://registry-1.docker.io/v2/ratelimitpreview/test/manifests/latest"

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error fetching token: %s\n", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 429 || resp.Header.Get("Ratelimit-Limit") == "" {
		return DockerHubQuota{
			Succeeded: false,
		}
	}

	limitSplited := strings.Split(resp.Header.Get("Ratelimit-Limit"), ";")
	limit, err := strconv.Atoi(limitSplited[0])
	if err != nil {
		log.Fatalf("Error fetching token5: %s\n", err)
	}
	limitWait, err := strconv.Atoi(strings.Replace(limitSplited[1], "w=", "", 1))
	if err != nil {
		log.Fatalf("Error fetching token4: %s\n", err)
	}

	remainingSplited := strings.Split(resp.Header.Get("Ratelimit-Remaining"), ";")
	remaining, err := strconv.Atoi(remainingSplited[0])
	if err != nil {
		log.Fatalf("Error fetching token3: %s\n", err)
	}
	remainingWait, err := strconv.Atoi(strings.Replace(remainingSplited[1], "w=", "", 1))
	if err != nil {
		log.Fatalf("Error fetching token2: %s\n", err)
	}

	return DockerHubQuota{
		Limit:         limit,
		LimitWait:     limitWait,
		Remaining:     remaining,
		RemainingWait: remainingWait,
		Ip:            resp.Header.Get("Docker-Ratelimit-Source"),
		Succeeded:     true,
	}
}
