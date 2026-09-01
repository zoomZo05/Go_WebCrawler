package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	URL            string   `json:"url"`
	Heading        string   `json:"heading"`
	FirstParagraph string   `json:"first_paragraph"`
	OutgoingLinks  []string `json:"outgoing_links"`
	ImageURLs      []string `json:"image_urls"`
}

func getHeadingFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(doc.Find("h1, h2").First().Text())
}

func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	// Prioritize paragraph inside <main> if present
	if mainP := doc.Find("main p").First(); mainP.Length() > 0 {
		return strings.TrimSpace(mainP.Text())
	}

	// Fallback to the first general paragraph
	return strings.TrimSpace(doc.Find("p").First().Text())
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	result := []string{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		parsedURL, err := url.Parse(href)
		if err != nil {
			return
		}
		result = append(result, baseURL.ResolveReference(parsedURL).String())
	})

	return result, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	result := []string{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	//square bracket image must have source
	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		parsedURL, err := url.Parse(src)
		if err != nil {
			return
		}
		result = append(result, baseURL.ResolveReference(parsedURL).String())
	})

	return result, nil
}

func extractPageData(html, pageURL string) PageData {
	url, err := url.Parse(pageURL)
	if err != nil {
		return PageData{URL: pageURL}
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return PageData{URL: pageURL}
	}

	listOfImage, err := getImagesFromHTML(html, url)
	if err != nil {
		return PageData{URL: pageURL}
	}

	listOfLink, err := getURLsFromHTML(html, url)
	if err != nil {
		return PageData{URL: pageURL}
	}

	data := PageData{
		URL:            pageURL,
		Heading:        strings.TrimSpace(doc.Find("h1").First().Text()),
		FirstParagraph: strings.TrimSpace(doc.Find("p").First().Text()),
		OutgoingLinks:  listOfLink,
		ImageURLs:      listOfImage,
	}

	return data
}

func getHTML(rawURL string) (string, error) {

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "BootCrawler/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", err
	}

	contentType := resp.Header.Get("content-type")
	if !strings.Contains(contentType, "text/html") {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
