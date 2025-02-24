package main

import (
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:     "multiple relative URLs",
			inputURL: "https://example.com",
			inputBody: `
<html>
	<body>
		<a href="/page1">Page 1</a>
		<a href="/page2">Page 2</a>
		<a href="/page3">Page 3</a>
	</body>
</html>
`,
			expected: []string{
				"https://example.com/page1",
				"https://example.com/page2",
				"https://example.com/page3",
			},
		},
		{
			name:     "no links in HTML",
			inputURL: "https://example.com",
			inputBody: `
<html>
	<body>
		<p>No links here!</p>
	</body>
</html>
`,
			expected: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %#v, got %#v", tc.expected, actual)
			}
		})
	}
}
