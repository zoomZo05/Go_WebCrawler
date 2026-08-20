package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	firstArg := os.Args[1]
	parsedURL, err := url.Parse(firstArg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't parse input URL: %v\n", err)
		os.Exit(1)
	}

	const maxConcurrency = 5
	cfg := config{
		pages:              map[string]PageData{},
		baseURL:            parsedURL,
		mu:                 &sync.RWMutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
	}

	cfg.wg.Add(1) // you must add 1 before run goroutine , avoid negative WaitGroup crash if you hit "done" early
	go cfg.crawlPage(firstArg)
	cfg.wg.Wait() // Task.Whenall() will unblock until WaitGroup = 0

	//wait for page , if do before cfg.wg.Wait() give same effect like you forget await keyword
	for url, page := range cfg.pages {
		fmt.Printf("Found %v internal links to %s\n", url, page.Heading)
	}

}
