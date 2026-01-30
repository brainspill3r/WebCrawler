package main

import (
	"fmt"
	"net/url"
	"sync"
	"testing"
)

func TestAddPageVisit(t *testing.T) {
	// Create a test config
	cfg := &config{
		pages:    make(map[string]PageData),
		mu:       &sync.Mutex{},
		maxPages: 10,
	}

	// Test first visit
	isFirst := cfg.addPageVisit("example.com")
	if !isFirst {
		t.Errorf("expected first visit to return true, got false")
	}

	// Check it was added to pages
	if len(cfg.pages) != 1 {
		t.Errorf("expected 1 page in map, got %d", len(cfg.pages))
	}

	// Test second visit to same page
	isFirst = cfg.addPageVisit("example.com")
	if isFirst {
		t.Errorf("expected second visit to return false, got true")
	}

	// Should still only have 1 page
	if len(cfg.pages) != 1 {
		t.Errorf("expected 1 page in map after duplicate, got %d", len(cfg.pages))
	}
}

func TestConfigWithMaxPages(t *testing.T) {
	baseURL, err := url.Parse("https://example.com")
	if err != nil {
		t.Fatalf("error parsing base URL: %v", err)
	}

	maxPages := 5
	cfg := &config{
		pages:              make(map[string]PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, 3),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	// Test that maxPages is set correctly
	if cfg.maxPages != maxPages {
		t.Errorf("expected maxPages to be %d, got %d", maxPages, cfg.maxPages)
	}

	// Test adding pages up to the limit
	for i := 0; i < maxPages; i++ {
		url := fmt.Sprintf("example.com/page%d", i)
		isFirst := cfg.addPageVisit(url)
		if !isFirst {
			t.Errorf("expected first visit to page %d to return true", i)
		}
	}

	if len(cfg.pages) != maxPages {
		t.Errorf("expected %d pages in map, got %d", maxPages, len(cfg.pages))
	}
}

func TestMaxPagesCheck(t *testing.T) {
	cfg := &config{
		pages:    make(map[string]PageData),
		mu:       &sync.Mutex{},
		maxPages: 2,
	}

	// Add pages up to the limit
	cfg.addPageVisit("page1.com")
	cfg.addPageVisit("page2.com")

	// Check if we've hit the limit
	cfg.mu.Lock()
	hitLimit := len(cfg.pages) >= cfg.maxPages
	cfg.mu.Unlock()

	if !hitLimit {
		t.Error("expected to hit maxPages limit with 2 pages")
	}

	// Try to add another page
	cfg.addPageVisit("page3.com")

	// Should still only have 2 pages if we properly respect the limit
	// Note: addPageVisit doesn't check the limit, that's done in crawlPage
	if len(cfg.pages) != 3 {
		t.Errorf("addPageVisit should still add pages (limit checking is in crawlPage), got %d pages", len(cfg.pages))
	}
}

func TestConcurrentAddPageVisit(t *testing.T) {
	cfg := &config{
		pages:    make(map[string]PageData),
		mu:       &sync.Mutex{},
		maxPages: 100, // Set high so we don't hit limit in test
	}

	// Test concurrent access to addPageVisit
	var wg sync.WaitGroup
	numGoroutines := 10

	// All goroutines try to add the same URL
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cfg.addPageVisit("example.com")
		}()
	}

	wg.Wait()

	// Should only have 1 page (first one wins, others return false)
	if len(cfg.pages) != 1 {
		t.Errorf("expected 1 page after concurrent access, got %d", len(cfg.pages))
	}
}
