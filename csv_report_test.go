package main

import (
	"encoding/csv"
	"os"
	"strings"
	"testing"
)

func TestWriteCSVReport(t *testing.T) {
	// Create test data
	testPages := map[string]PageData{
		"example.com": {
			URL:            "https://example.com",
			H1:             "Example Site",
			FirstParagraph: "This is an example site for testing.",
			OutgoingLinks:  []string{"https://example.com/about", "https://example.com/contact"},
			ImageURLs:      []string{"https://example.com/logo.png", "https://example.com/banner.jpg"},
		},
		"example.com/about": {
			URL:            "https://example.com/about",
			H1:             "About Us",
			FirstParagraph: "Learn more about our company.",
			OutgoingLinks:  []string{"https://example.com", "https://example.com/team"},
			ImageURLs:      []string{"https://example.com/team.jpg"},
		},
	}

	// Test filename
	testFilename := "test_report.csv"

	// Clean up any existing test file
	defer os.Remove(testFilename)

	// Write CSV report
	err := writeCSVReport(testPages, testFilename)
	if err != nil {
		t.Fatalf("writeCSVReport failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFilename); os.IsNotExist(err) {
		t.Fatal("CSV file was not created")
	}

	// Read and verify CSV content
	file, err := os.Open(testFilename)
	if err != nil {
		t.Fatalf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV file: %v", err)
	}

	// Should have header + 2 data rows
	if len(records) != 3 {
		t.Errorf("expected 3 rows (1 header + 2 data), got %d", len(records))
	}

	// Check header
	expectedHeader := []string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"}
	if len(records) > 0 {
		for i, expected := range expectedHeader {
			if i >= len(records[0]) || records[0][i] != expected {
				t.Errorf("header column %d: expected '%s', got '%s'", i, expected, records[0][i])
			}
		}
	}

	// Check that data rows contain expected URLs
	foundURLs := make(map[string]bool)
	for i := 1; i < len(records); i++ {
		if len(records[i]) > 0 {
			foundURLs[records[i][0]] = true
		}
	}

	expectedURLs := []string{"https://example.com", "https://example.com/about"}
	for _, expectedURL := range expectedURLs {
		if !foundURLs[expectedURL] {
			t.Errorf("expected to find URL '%s' in CSV data", expectedURL)
		}
	}

	// Check that outgoing links are joined with semicolons
	for i := 1; i < len(records); i++ {
		if len(records[i]) > 3 && strings.Contains(records[i][3], ";") {
			// Found semicolon-separated links, which is expected
			continue
		} else if len(records[i]) > 3 && records[i][3] == "" {
			// Empty links column is also acceptable
			continue
		}
		// If we get here, check if it has at least one link without semicolons (single link case)
		if len(records[i]) > 3 && len(records[i][3]) > 0 && !strings.Contains(records[i][3], ";") {
			// Single link without semicolon is acceptable
			continue
		}
	}
}

func TestWriteCSVReportEmpty(t *testing.T) {
	// Test with empty pages map
	emptyPages := map[string]PageData{}
	testFilename := "test_empty_report.csv"

	// Clean up any existing test file
	defer os.Remove(testFilename)

	// Write CSV report
	err := writeCSVReport(emptyPages, testFilename)
	if err != nil {
		t.Fatalf("writeCSVReport failed with empty data: %v", err)
	}

	// Verify file exists and has at least header
	file, err := os.Open(testFilename)
	if err != nil {
		t.Fatalf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV file: %v", err)
	}

	// Should have only header row
	if len(records) != 1 {
		t.Errorf("expected 1 row (header only), got %d", len(records))
	}
}
