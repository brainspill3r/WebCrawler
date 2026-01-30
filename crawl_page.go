package main

import (
	"fmt"
	"net/url"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	// Parse the base and current URLs
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("error parsing base URL: %v\n", err)
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("error parsing current URL: %v\n", err)
		return
	}

	// Check if current URL is on the same domain as base URL
	if baseURL.Host != currentURL.Host {
		return
	}

	// Normalize the current URL
	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("error normalizing URL: %v\n", err)
		return
	}

	// Check if we've already visited this page
	if _, visited := pages[normalizedURL]; visited {
		pages[normalizedURL]++
		return
	}

	// Mark this page as visited
	pages[normalizedURL] = 1

	// Print what we're crawling (so you can see progress)
	fmt.Printf("crawling: %s\n", rawCurrentURL)

	// Fetch the HTML from the current URL
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error fetching HTML: %v\n", err)
		return
	}

	// Extract all URLs from the HTML
	urls, err := getURLsFromHTML(html, currentURL)
	if err != nil {
		fmt.Printf("error getting URLs from HTML: %v\n", err)
		return
	}

	// Recursively crawl each URL found on the page
	for _, nextURL := range urls {
		crawlPage(rawBaseURL, nextURL, pages)
	}
}
