package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(os.Args) > 4 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	firstArg := os.Args[1]

	maxConcurrency := 1
	maxPages := 10

	if len(os.Args) >= 3 {
		val, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid maxConcurrency: %v\n", err)
			os.Exit(1)
		}
		maxConcurrency = val
	}

	if len(os.Args) == 4 {
		val, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid maxPages: %v\n", err)
			os.Exit(1)
		}
		maxPages = val
	}

	parsedURL, err := url.Parse(firstArg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't parse input URL: %v\n", err)
		os.Exit(1)
	}
	cfg := config{
		pages:              map[string]PageData{},
		baseURL:            parsedURL,
		mu:                 &sync.RWMutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	cfg.wg.Add(1) // you must add 1 before run goroutine , avoid negative WaitGroup crash if you hit "done" early
	go cfg.crawlPage(firstArg)
	cfg.wg.Wait() // Task.Whenall() will unblock until WaitGroup = 0

	//wait for page , if do before cfg.wg.Wait() give same effect like you forget await keyword
	for normalizedURL := range cfg.pages {
		fmt.Printf("found: %s\n", normalizedURL)
	}

	writeJSONReport(cfg.pages, "report.json")

}
