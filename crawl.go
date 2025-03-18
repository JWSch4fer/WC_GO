package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

func crawlPage(ctx context.Context, raw_base_URL, raw_current_URL string, pages map[string]int) {

	//check if I ran out of time...
	select {
	case <-ctx.Done():
		return
	default:
		//continue
	}

	//don't constantly bombard the site
	be_polite := rand.Float32()
	if be_polite > 0.667 {
		time.Sleep(1 * time.Second)

	}

	//ensure current URL is on the same domain as raw_base_URL
	parsed_base, err := url.Parse(raw_base_URL)
	if err != nil {
		fmt.Printf("Error parsing base URL %q: %v\n", raw_base_URL, err)
		return
	}
	parsed_current, err := url.Parse(raw_current_URL)
	if err != nil {
		fmt.Printf("Error parsing URL %q: %v\n", raw_current_URL, err)
		return
	}
	if parsed_base.Host != parsed_current.Host {
		//dirrent Domain, skip
		return
	}

	//normalize the current URL
	normalized_URL, err := normalizeURL(raw_current_URL)
	if err != nil {
		fmt.Printf("Error nomalizing URL %q: %v\n", raw_current_URL, err)
	}

	if _, ok := pages[normalized_URL]; ok {
		pages[normalized_URL]++
		fmt.Printf("count: %d  | page: %s \n", pages[normalized_URL], normalized_URL)
		return
	}

	//if this is a new page
	pages[normalized_URL] = 1

	//get html for the current page
	htmlBody, err := getHTML(raw_current_URL)
	if err != nil {
		fmt.Printf("Error fetching HTML from %q: %v\n", raw_current_URL, err)
	}

	links, err := getURLsFromHTML(htmlBody, raw_current_URL)
	if err != nil {
		fmt.Printf("Error fetching links from %q: %v\n", raw_current_URL, err)
	}

	//recursively crawl each discovered links
	for _, link := range links {
		select {
		case <-ctx.Done():
			return
		default:
		}
		crawlPage(ctx, raw_base_URL, link, pages)
	}

}
