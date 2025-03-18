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

	path := parsedURL.Path
	if path == "" {
		//there was no provided path
		path = "/"
	} else {
		path = strings.TrimSuffix(path, "/")
	}
	return parsedURL.Host + path, nil
}
