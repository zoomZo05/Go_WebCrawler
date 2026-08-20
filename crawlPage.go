package main

import (
	"fmt"
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.RWMutex //share read
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func (c *config) crawlPage(rawCurrentURL string) {

	// Acquire semaphore slot & register goroutine , and this will be recursive call.
	defer c.wg.Done()            // Decrement waitgroup "-1"
	c.concurrencyControl <- struct{}{} // this will be added until 5
	defer func() {
		<-c.concurrencyControl // Release slot
	}()

	current, err := url.Parse(rawCurrentURL)
	if err != nil {
		return
	}

	if c.baseURL.Hostname() != current.Hostname() {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	if !c.addPageIfNew(normalizedURL) {
		return
	}
	//end of base case

	fmt.Printf("Crawling: %s\n", rawCurrentURL)

	result, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}

	currentPageData := extractPageData(result, rawCurrentURL)

	c.mu.Lock()
	c.pages[normalizedURL] = currentPageData
	c.mu.Unlock()

	listOfUrl, err := getURLsFromHTML(result, current)
	if err != nil {
		return
	}

	for _, href := range listOfUrl {
		//avoid deadlock run fire and forget , if no "go" here, it will never reach "defer"
		c.wg.Add(1) // +1
		go c.crawlPage(href)
	}

	
	//defer is working here channel is free 1 slot
}

// check and claim technique
// addPageIfNew atomically checks if a URL is visited.
// Returns true if this goroutine claimed it, false if already visited.

func (c *config) addPageIfNew(normalizedURL string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.pages[normalizedURL]; exists {
		return false
	}

	// Claim ownership ---- so in the case of same URL calling this will not allow same map "key"
	//hint we lock both check key and write key
	c.pages[normalizedURL] = PageData{}
	return true
}
