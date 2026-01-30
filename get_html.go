package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set User-Agent header
	req.Header.Set("User-Agent", "BootCrawler/1.0")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP error status codes (400+)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: status code %d", resp.StatusCode)
	}

	// Check content-type header
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("invalid content type: %s (expected text/html)", contentType)
	}

	// Read the response body
	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(htmlBytes), nil
}
