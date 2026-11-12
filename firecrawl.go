package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type FirecrawlRequest struct {
	URL     string   `json:"url"`
	Formats []string `json:"formats"`
}

type FirecrawlResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Data    struct {
		Markdown string `json:"markdown"`
		HTML     string `json:"html,omitempty"`
		Metadata struct {
			Title       string      `json:"title"`
			Description string      `json:"description"`
			Language    interface{} `json:"language"`
			SourceURL   string      `json:"sourceURL"`
			StatusCode  int         `json:"statusCode,omitempty"`
		} `json:"metadata"`
	} `json:"data"`
}

func scrapeURL(url string) (*FirecrawlResponse, error) {
	apiURL := "https://api.firecrawl.dev/v2/scrape"
	authToken := os.Getenv("FC_AUTH_TOKEN")

	if authToken == "" {
		return nil, fmt.Errorf("environment variable FC_AUTH_TOKEN is not set")
	}

	// Create the JSON payload
	payload := FirecrawlRequest{
		URL:     url,
		Formats: []string{"markdown"},
	}
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
