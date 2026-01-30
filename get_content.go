package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getH1FromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	h1 := doc.Find("h1").First()
	return strings.TrimSpace(h1.Text())
}

func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	// First, try to find a <p> tag within <main>
	mainSection := doc.Find("main")
	if mainSection.Length() > 0 {
		paragraph := mainSection.Find("p").First()
		if paragraph.Length() > 0 {
			return strings.TrimSpace(paragraph.Text())
		}
	}

	// Fallback: find first <p> tag in the entire document
	paragraph := doc.Find("p").First()
	return strings.TrimSpace(paragraph.Text())
}
