package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestCrawlPage(t *testing.T) {
	base_URL := []string{"https://webscraper.io/test-sites/e-commerce/allinone", "https://demo.cyotek.com/html/index.php", "https://quotes.toscrape.com"}

	for _, url := range base_URL {
		pages := make(map[string]int)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		crawlPage(ctx, url, url, pages)

		if len(pages) == 0 {
			t.Fatalf("expected to discover at least one page, got 0\n")
		}

		fmt.Println("Discovered Pages:")
		for page, count := range pages {
			fmt.Printf(" %d -> %s\n", count, page)

		}

	}
}
