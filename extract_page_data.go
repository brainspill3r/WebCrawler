package main

import "net/url"

type PageData struct {
	URL            string
	H1             string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func extractPageData(html, pageURL string) PageData {
	// Parse the page URL
	baseURL, err := url.Parse(pageURL)
	if err != nil {
		// If URL parsing fails, return empty PageData with the original URL
		return PageData{
			URL:            pageURL,
			H1:             "",
			FirstParagraph: "",
			OutgoingLinks:  []string{},
			ImageURLs:      []string{},
		}
	}

	// Extract H1
	h1 := getH1FromHTML(html)

	// Extract first paragraph
	firstParagraph := getFirstParagraphFromHTML(html)

	// Extract outgoing links
	outgoingLinks, err := getURLsFromHTML(html, baseURL)
	if err != nil {
		outgoingLinks = []string{}
	}

	// Extract image URLs
	imageURLs, err := getImagesFromHTML(html, baseURL)
	if err != nil {
		imageURLs = []string{}
	}

	return PageData{
		URL:            pageURL,
		H1:             h1,
		FirstParagraph: firstParagraph,
		OutgoingLinks:  outgoingLinks,
		ImageURLs:      imageURLs,
	}
}
