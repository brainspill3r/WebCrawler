package main

import (
	"encoding/csv"
	"os"
	"strings"
)

// writeCSVReport exports crawl results to a CSV file
func writeCSVReport(pages map[string]PageData, filename string) error {
	// Create the CSV file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	headers := []string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write data rows
	for _, pageData := range pages {
		// Join slices with semicolons as specified
		outgoingLinks := strings.Join(pageData.OutgoingLinks, ";")
		imageURLs := strings.Join(pageData.ImageURLs, ";")

		// Create row with all the data
		row := []string{
			pageData.URL,
			pageData.H1,
			pageData.FirstParagraph,
			outgoingLinks,
			imageURLs,
		}

		// Write the row
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
