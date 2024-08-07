package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type FirecrawlResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Content  string `json:"content"`
		Markdown string `json:"markdown"`
		Metadata struct {
			Title       string      `json:"title"`
			Description string      `json:"description"`
			Language    interface{} `json:"language"`
			SourceURL   string      `json:"sourceURL"`
		} `json:"metadata"`
	} `json:"data"`
}

func scrapeURL(url string) error {
	apiURL := "https://api.firecrawl.dev/v0/scrape"
	authToken := os.Getenv("FC_AUTH_TOKEN")

	if authToken == "" {
		return fmt.Errorf("environment variable FC_AUTH_TOKEN is not set")
	}

	// Create the JSON payload
	payload := map[string]string{"url": url}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	fmt.Println("Request successful")
	return nil
}
