package main

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"testing"
	"time"
)

func TestCrawlPage(t *testing.T) {
	base_URL := []string{
		"https://webscraper.io/test-sites/e-commerce/allinone",
		"https://demo.cyotek.com/html/index.php",
		"https://quotes.toscrape.com"}

	for _, rawURL := range base_URL {

		//parse the base url
		parsedBase, err := url.Parse(rawURL)
		if err != nil {
			t.Fatalf("Error parsing base URl %q: %v", rawURL, err)
		}

		//create the config struct
		//hard code maxConcurrency and maxPages
		cfg := &config{
			pages:              make(map[string]int),
			baseURL:            parsedBase,
			mu:                 &sync.Mutex{},
			concurrencyControl: make(chan struct{}, 5),
			wg:                 &sync.WaitGroup{},
			maxPages:           20,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cfg.crawlPage(ctx, rawURL)
		cfg.wg.Wait() // wait until all crawling go routines finish

		if len(cfg.pages) == 0 {
			t.Fatalf("expected to discover at least one page, got 0\n")
		}

		fmt.Println("Discovered Pages: %w", rawURL)
		for page, count := range cfg.pages {
			fmt.Printf(" %d -> %s\n", count, page)
		}

	}
}
