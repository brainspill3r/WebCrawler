package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	// Check if the correct number of arguments was provided
	if len(os.Args) < 4 {
		fmt.Println("Usage: ./crawler URL maxConcurrency maxPages")
		fmt.Println("Example: ./crawler https://example.com 3 10")
		os.Exit(1)
	}

	if len(os.Args) > 4 {
		fmt.Println("too many arguments provided")
		fmt.Println("Usage: ./crawler URL maxConcurrency maxPages")
		os.Exit(1)
	}

	// Get the arguments from command line
	rawURL := os.Args[1]
	maxConcurrencyStr := os.Args[2]
	maxPagesStr := os.Args[3]

	// Parse maxConcurrency
	maxConcurrency, err := strconv.Atoi(maxConcurrencyStr)
	if err != nil {
		fmt.Printf("error parsing maxConcurrency '%s': %v\n", maxConcurrencyStr, err)
		os.Exit(1)
	}
	if maxConcurrency < 1 {
		fmt.Println("maxConcurrency must be at least 1")
		os.Exit(1)
	}

	// Parse maxPages
	maxPages, err := strconv.Atoi(maxPagesStr)
	if err != nil {
		fmt.Printf("error parsing maxPages '%s': %v\n", maxPagesStr, err)
		os.Exit(1)
	}
	if maxPages < 1 {
		fmt.Println("maxPages must be at least 1")
		os.Exit(1)
	}

	// Parse the base URL
	baseURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Printf("error parsing base URL: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("starting crawl of: %s\n", rawURL)
	fmt.Printf("maxConcurrency: %d\n", maxConcurrency)
	fmt.Printf("maxPages: %d\n", maxPages)

	// Create the config struct per assignment
	cfg := &config{
		pages:              make(map[string]PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	// wg.Add before spawning goroutine (as per assignment tips)
	cfg.wg.Add(1)
	go cfg.crawlPage(rawURL)

	// Single wg.Wait() in main (as per assignment tips)
	cfg.wg.Wait()

	// Generate CSV report
	filename := "report.csv"
	fmt.Printf("\nGenerating CSV report: %s\n", filename)

	err = writeCSVReport(cfg.pages, filename)
	if err != nil {
		fmt.Printf("error writing CSV report: %v\n", err)
		os.Exit(1)
	}

	// Print basic summary
	fmt.Printf("Crawl completed: %d pages found (max: %d)\n", len(cfg.pages), maxPages)
	fmt.Printf("Report saved to %s\n", filename)
}
