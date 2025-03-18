package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {

	//make sure a website was provided
	args := os.Args
	if len(os.Args) != 4 {
		fmt.Println("Usage: WC_GO <website> <maxConcurrency> <maxPages>")
		os.Exit(1)
	}

	maxConcurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("effor parsing maxConcurrency: %v\n", err)
		os.Exit(1)
	}

	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Printf("error parsing maxPages: %v\n", err)
	}

	//parse the input to start crawling
	rawURL := os.Args[1]
	parsedBase, err := url.Parse(rawURL)
	if err != nil {
		fmt.Printf("Error parsing base URL %q: %v\n", rawURL, err)
	}

	fmt.Printf("starting crawl of: %s\n", args[1])

	//create the config struct
	cfg := &config{
		pages:              make(map[string]int),
		baseURL:            parsedBase,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	//context to monitor the amount of runtime
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//start crawling
	cfg.crawlPage(ctx, rawURL)
	cfg.wg.Wait() //wait for all requests to complete

	fmt.Printf("\nCrawling pages:\n")
	for page, count := range cfg.pages {
		fmt.Printf(" %d -> %s\n", count, page)
	}

}
