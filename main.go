package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	// map of the available cli commands
	args := os.Args

	if len(args) == 1 {
		fmt.Printf("no website provided\n")
		os.Exit(1)
	}
	if len(args) > 2 {
		fmt.Printf("too many arguments provided\n")
		os.Exit(1)
	}
	fmt.Printf("starting crawl of: %s\n", args[1])
	// raw_HTML, err := getHTML(args[1])
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("%s\n", raw_HTML)

	pages := make(map[string]int)

	//start crawling
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	crawlPage(ctx, args[1], args[1], pages)

	fmt.Printf("\nCrawling pages:\n")
	for page, count := range pages {
		fmt.Printf(" %d -> %s\n", count, page)
	}

}

// getHTML fetches the HTML from rawURL, ensuring the response is both valid
func getHTML(rawURL string) (string, error) {

	resp, err := http.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch the webpage: %w", err)
	}
	defer resp.Body.Close()

	// Check for 4xx or 5xx status codes
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("received an error status code: %d", resp.StatusCode)
	}

	// Verify the response is text/html
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("expected 'text/html' content type but got '%s'", contentType)
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(bodyBytes), nil
}
