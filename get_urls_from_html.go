package main

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	urls := []string{}
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		// Parse the href to handle relative URLs
		parsedURL, err := url.Parse(href)
		if err != nil {
			return
		}

		// Resolve relative URLs against the base URL
		absoluteURL := baseURL.ResolveReference(parsedURL)
		urls = append(urls, absoluteURL.String())
	})

	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	images := []string{}
	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists || src == "" {
			return
		}

		// Parse the src to handle relative URLs
		parsedURL, err := url.Parse(src)
		if err != nil {
			return
		}

		// Resolve relative URLs against the base URL
		absoluteURL := baseURL.ResolveReference(parsedURL)
		images = append(images, absoluteURL.String())
	})

	return images, nil
}
