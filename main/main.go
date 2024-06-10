package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	url := os.Getenv("INPUT_URL")
	responseCode := os.Getenv("INPUT_RESPONSECODE")
	timeout, err := strconv.Atoi(os.Getenv("INPUT_TIMEOUT"))
	if err != nil {
		fmt.Println("Timeout must be an integer")
		os.Exit(1)
	}
	interval, err := strconv.Atoi(os.Getenv("INPUT_INTERVAL"))
	if err != nil {
		fmt.Println("Interval must be an integer")
		os.Exit(1)
	}
	localhost := os.Getenv("INPUT_LOCALHOST")

	username := os.Getenv("INPUT_USERNAME")
	password := os.Getenv("INPUT_PASSWORD")

	fmt.Printf("Polling URL `%s` as `%s` for response code %s for up to %d ms at %d ms intervals\n", url, username, responseCode, timeout, interval)
	startTime := time.Now()
	timeoutDuration := time.Duration(timeout) * time.Millisecond
	sleepDuration := time.Duration(interval) * time.Millisecond

	if localhost != "" && strings.Contains(url, "localhost") {
		url = strings.ReplaceAll(url, "localhost", localhost)
	}

	client := &http.Client{}
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		fmt.Printf("Request Construction: %v", err)
	}

	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	codes := strings.Split(responseCode, ",")

	for {
		res, err := client.Do(req)
		if err == nil && containsStr(codes, strconv.Itoa(res.StatusCode)) {
			fmt.Printf("Response header: %v", res)
			os.Exit(0)
		}
		time.Sleep(sleepDuration)
		elapsed := time.Now().Sub(startTime)
		if elapsed > timeoutDuration {
			fmt.Printf("Timed out\n")
			os.Exit(1)
		}
	}
}

func containsStr(haystack []string, needle string) bool {
	for _, straw := range haystack {
		if needle == straw {
			return true
		}
	}
	return false
}
