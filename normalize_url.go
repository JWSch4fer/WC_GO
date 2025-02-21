package main

import (
	"net/url"
	"strings"
)

// normalize a URL string and return the string
func normalizeURL(input string) (string, error) {
	// basic implementation
	parsedURL, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	normalizedPath := strings.TrimSuffix(parsedURL.Path, "/")
	return parsedURL.Host + normalizedPath, nil
}
