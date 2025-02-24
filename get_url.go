package main

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	urls := make([]string, 0)

	//start with an io-reader for html body
	reader := strings.NewReader(htmlBody)
	doc, err := html.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	//use recursion to travers html tree
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		// check if this is an anchor
		if n.Type == html.ElementNode && n.Data == "a" {
			//look for links
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link := attr.Val
					absLink, err := resolveURL(link, rawBaseURL)
					if err == nil {
						urls = append(urls, absLink)
					}
					// once we have a linke break out of the attr loop
					break
				}
			}
		}
		//recursively travers child nodes
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}
	traverse(doc)

	return urls, nil
}

// take in a linke and a base url and return absolute urls
func resolveURL(link, base string) (string, error) {
	parsedLink, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	// if link is absolute return
	if parsedLink.IsAbs() {
		return parsedLink.String(), nil
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	resolvedURL := baseURL.ResolveReference(parsedLink)
	return resolvedURL.String(), nil
}
