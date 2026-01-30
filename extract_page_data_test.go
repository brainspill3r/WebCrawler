package main

import (
	"reflect"
	"testing"
)

func TestExtractPageData(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body>
        <h1>Test Title</h1>
        <p>This is the first paragraph.</p>
        <a href="/link1">Link 1</a>
        <img src="/image1.jpg" alt="Image 1">
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://blog.boot.dev",
		H1:             "Test Title",
		FirstParagraph: "This is the first paragraph.",
		OutgoingLinks:  []string{"https://blog.boot.dev/link1"},
		ImageURLs:      []string{"https://blog.boot.dev/image1.jpg"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestExtractPageDataWithMain(t *testing.T) {
	inputURL := "https://example.com"
	inputBody := `<html><body>
		<h1>Page Title</h1>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
		<a href="https://other.com">External</a>
		<a href="/internal">Internal</a>
		<img src="/logo.png" alt="Logo">
	</body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://example.com",
		H1:             "Page Title",
		FirstParagraph: "Main paragraph.",
		OutgoingLinks:  []string{"https://other.com", "https://example.com/internal"},
		ImageURLs:      []string{"https://example.com/logo.png"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestExtractPageDataEmpty(t *testing.T) {
	inputURL := "https://empty.com"
	inputBody := `<html><body></body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://empty.com",
		H1:             "",
		FirstParagraph: "",
		OutgoingLinks:  []string{},
		ImageURLs:      []string{},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}
