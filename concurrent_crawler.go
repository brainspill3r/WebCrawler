package main

import (
	"fmt"
	"net/url"
	"sync"
)

// config struct for concurrent crawling
type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int // Maximum number of pages to crawl
}

// addPageVisit helper method
func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, exists := cfg.pages[normalizedURL]; exists {
		return false // Already visited
	}

	cfg.pages[normalizedURL] = PageData{} // Mark as visited
	return true                           // First visit
}

// crawlPage method - updated signature per assignment
func (cfg *config) crawlPage(rawCurrentURL string) {
	// Check if we've hit maxPages limit at the very start
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		cfg.wg.Done() // Still need to decrement wait group
		return
	}
	cfg.mu.Unlock()

	// Send to concurrency control channel
	cfg.concurrencyControl <- struct{}{}

	// Defer cleanup as recommended in assignment
	defer func() {
		<-cfg.concurrencyControl // Remove from channel
		cfg.wg.Done()            // Decrement wait group
	}()

	// Parse current URL
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("error parsing URL %s: %v\n", rawCurrentURL, err)
		return
	}

	// Check if same domain
	if currentURL.Host != cfg.baseURL.Host {
		return
	}

	// Normalize URL
	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("error normalizing URL %s: %v\n", rawCurrentURL, err)
		return
	}

	// Check if already visited using helper method
	if !cfg.addPageVisit(normalizedURL) {
		return // Already crawled this page
	}

	// Check again if we've hit maxPages after adding this page
	cfg.mu.Lock()
	if len(cfg.pages) > cfg.maxPages {
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()

	fmt.Printf("crawling: %s\n", rawCurrentURL)

	// Get HTML
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error fetching HTML from %s: %v\n", rawCurrentURL, err)
		return
	}

	// Extract and store page data
	pageData := extractPageData(html, rawCurrentURL)
	cfg.mu.Lock()
	cfg.pages[normalizedURL] = pageData
	cfg.mu.Unlock()

	// Get URLs from this page
	urls, err := getURLsFromHTML(html, currentURL)
	if err != nil {
		fmt.Printf("error getting URLs from HTML for %s: %v\n", rawCurrentURL, err)
		return
	}

	// Spawn goroutines for each URL (wg.Add before spawning as per tips)
	for _, nextURL := range urls {
		cfg.wg.Add(1)
		go cfg.crawlPage(nextURL)
	}
}
