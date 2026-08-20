package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetHeadingFromHTMLBasic(t *testing.T) {
	inputBody := "<html><body><h1>Test Title</h1></body></html>"
	actual := getHeadingFromHTML(inputBody)
	expected := "Test Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	inputBody := `<html><body>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "Main paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetURLsFromHTMLAbsolute(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:      "absolute URL",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="https://crawler-test.com/path">Boot.dev</a></body></html>`,
			expected:  []string{"https://crawler-test.com/path"},
		},
		{
			name:      "relative URL with root slash",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="/path/one">Link</a></body></html>`,
			expected:  []string{"https://crawler-test.com/path/one"},
		},
		{
			name:      "relative URL without root slash and nested path base",
			inputURL:  "https://crawler-test.com/sub/dir/",
			inputBody: `<html><body><a href="page.html">Link</a></body></html>`,
			expected:  []string{"https://crawler-test.com/sub/dir/page.html"},
		},
		{
			name:     "multiple links including nested elements",
			inputURL: "https://crawler-test.com",
			inputBody: `<html><body>
				<a href="/first">First</a>
				<div><a href="https://other-domain.com/second"><span>Second</span></a></div>
			</body></html>`,
			expected: []string{"https://crawler-test.com/first", "https://other-domain.com/second"},
		},
		{
			name:      "anchor tag without href",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a>No Href</a><a href="/valid">Valid</a></body></html>`,
			expected:  []string{"https://crawler-test.com/valid"},
		},
		{
			name:      "no links in body",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><p>Plain text</p></body></html>`,
			expected:  []string{},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}

			actual, err := getURLsFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Fatalf("%v --- unexpected error: %v", i, err)
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("%v --- expected %v, got %v ", i, tc.expected, actual)
			}
		})
	}
}

func TestGetImagesFromHTMLRelative(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:      "relative image source",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><img src="/logo.png" alt="Logo"></body></html>`,
			expected:  []string{"https://crawler-test.com/logo.png"},
		},
		{
			name:      "absolute image source",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><img src="https://cdn.example.com/pic.jpg"></body></html>`,
			expected:  []string{"https://cdn.example.com/pic.jpg"},
		},
		{
			name:     "multiple images",
			inputURL: "https://crawler-test.com",
			inputBody: `<html><body>
				<img src="/img1.png">
				<img src="/img2.png">
			</body></html>`,
			expected: []string{"https://crawler-test.com/img1.png", "https://crawler-test.com/img2.png"},
		},
		{
			name:      "img tag without src",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><img alt="placeholder"><img src="/valid.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/valid.png"},
		},
		{
			name:      "no images in body",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><p>No images</p></body></html>`,
			expected:  []string{},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}

			actual, err := getImagesFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Fatalf("%v --- unexpected error: %v", i, err)
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("%v --- expected %v, got %v ", i, tc.expected, actual)
			}
		})
	}
}

func TestExtractPageData(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body>
        <h1>Test Title</h1>
        <p>This is the first paragraph.</p>
        <a href="/link1">Link 1</a>
        <img src="/image1.jpg" alt="Image 1">
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://crawler-test.com",
		Heading:        "Test Title",
		FirstParagraph: "This is the first paragraph.",
		OutgoingLinks:  []string{"https://crawler-test.com/link1"},
		ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

