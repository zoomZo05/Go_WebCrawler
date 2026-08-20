package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	URL            string
	Heading        string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
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


func getHTML(rawURL string) (string, error){
		
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


func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int){
	base, err := url.Parse(rawBaseURL)
	if err != nil {
		return
	}

	current, err := url.Parse(rawCurrentURL)
	if err != nil {
		return
	}

	if base.Hostname() != current.Hostname(){
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	if _, ok := pages[normalizedURL]; ok {
		pages[normalizedURL]++
		return
	}

	pages[normalizedURL] = 1

	fmt.Printf("Crawling: %s\n", rawCurrentURL)

	result ,err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}

	listOfUrl,err := getURLsFromHTML(result,current)
	if err != nil {
		return
	}

	for _, href := range listOfUrl{
		crawlPage(rawBaseURL, href, pages)
	}
}

