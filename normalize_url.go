package main

import (
	"net/url"
	"strings"
)

func normalizeURL(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	normalizedPath := strings.TrimSuffix(parsedURL.Path, "/")
	normalizedURL := parsedURL.Host + normalizedPath
	normalizedURL = strings.ToLower(normalizedURL)

	return normalizedURL, nil
}
