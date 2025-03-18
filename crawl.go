package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

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

func (cfg *config) crawlPage(ctx context.Context, rawCurrentURL string) {

	//check if we have reached max pages
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()

	//check if I ran out of time...
	select {
	case <-ctx.Done():
		return
	default:
		//continue
	}

	//Block if we reached max concurrency and add to wait group
	select {
	case cfg.concurrencyControl <- struct{}{}:
	case <-ctx.Done():
		return
	}
	cfg.wg.Add(1)
	defer func() {
		<-cfg.concurrencyControl //release the slot
		cfg.wg.Done()
	}()

	//slow down to be be polite
	if rand.Float32() > 0.667 {
		time.Sleep(1 * time.Second)
	}

	//parse the current url
	parsedCurrent, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error parsing URL %q: %v\n", rawCurrentURL, err)
	}

	//ensure current is on the same host as the base url
	if parsedCurrent.Host != cfg.baseURL.Host {
		return
	}

	//normalize
	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error normalizing URL %q: %v\n", rawCurrentURL, err)
	}

	//add/update the visits
	if !cfg.addPageVisit(normalizedURL) {
		fmt.Printf("Already visited: %s\n", normalizedURL)
	}

	fmt.Printf("Visiting: %s\n", normalizedURL)

	htmlBody, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error fetching HTML from %q: %v\n", rawCurrentURL, err)
	}

	//extract urls from html
	links, err := getURLsFromHTML(htmlBody, rawCurrentURL)
	if err != nil {
		fmt.Printf("Error extracting links from %q: %v\n", rawCurrentURL, err)
	}

	//recursively call each link
	for _, link := range links {
		select {
		case <-ctx.Done():
			return
		default:
		}

		//before spawning, check maxPages
		cfg.mu.Lock()
		if len(cfg.pages) >= cfg.maxPages {
			cfg.mu.Unlock()
			return
		}
		cfg.mu.Unlock()

		go cfg.crawlPage(ctx, link)
	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {

	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	if _, ok := cfg.pages[normalizedURL]; ok {
		cfg.pages[normalizedURL]++
		return false
	}
	cfg.pages[normalizedURL] = 1
	return true
}
