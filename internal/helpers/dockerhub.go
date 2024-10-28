package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-logr/logr"
)

type AuthResponse struct {
	Token string `json:"token"`
}

func GetDockerhubLimit(setupLog logr.Logger, resultChan chan<- int) {
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

	token, _, err := new(jwt.Parser).ParseUnverified(authResp.Token, jwt.MapClaims{})
	if err != nil {
		log.Fatalf("Error parsing token: %s\n", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		var pull_limit = claims["access"].([]interface{})[0].(map[string]interface{})["parameters"].(map[string]interface{})["pull_limit"]
		i, err := strconv.Atoi(pull_limit.(string))
		if err != nil {
			fmt.Printf("Error converting string to int: %s\n", err)
			resultChan <- 0
		} else {
			resultChan <- i
		}
	} else {
		log.Fatalf("Invalid token claims")
	}

	resultChan <- 0
}
