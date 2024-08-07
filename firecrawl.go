package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func scrapeURL(url string) (*FirecrawlResponse, error) {
	apiURL := "https://api.firecrawl.dev/v0/scrape"
	authToken := os.Getenv("FC_AUTH_TOKEN")

	if authToken == "" {
		return nil, fmt.Errorf("environment variable FC_AUTH_TOKEN is not set")
	}

	// Create the JSON payload
	payload := map[string]string{"url": url}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the JSON response into the FirecrawlResponse structure
	var firecrawlResponse FirecrawlResponse
	if err := json.Unmarshal(body, &firecrawlResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %v", err)
	}

	return &firecrawlResponse, nil
}
