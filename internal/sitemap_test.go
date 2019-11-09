package internal

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestNewCrawler(t *testing.T) {
	var scenarios = []struct {
		name        string
		url         string
		expectedErr error
	}{
		{
			name:        "invalid host name",
			url:         "invalid",
			expectedErr: errors.New("invalid host name invalid"),
		},
		{
			name:        "valid host",
			url:         "monzo.com",
			expectedErr: nil,
		},
		{
			name:        "valid host",
			url:         "mail.google.com/mail/u/0/#inbox",
			expectedErr: nil,
		},
	}
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			_, err := NewCrawler(sc.url)
			checkErr(t, err, sc.expectedErr)
		})
	}
}

func TestCrawl(t *testing.T) {
	scenarios := []struct {
		name string
		url  string
	}{
		{
			name: "google",
			url:  "google.co.uk",
		},
	}
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			c, _ := NewCrawler(sc.url)
			links, err := c.Crawl()
			if err != nil {
				t.Error(err)
			}
			fmt.Printf("links %v", links)
		})
	}
}

func TestTokenize(t *testing.T) {
	scenarios := []struct {
		name        string
		html        string
		expectedErr error
	}{
		{
			name:        "not html",
			html:        "not html",
			expectedErr: nil,
		},
		{
			name:        "no links",
			html:        htmlStr("<p>"),
			expectedErr: nil,
		},
		{
			name:        "dead link",
			html:        htmlStr(`<a>Test</a>`),
			expectedErr: nil,
		},
		{
			name:        "invalid href link",
			html:        htmlStr(`<a href="#">Test</a>`),
			expectedErr: nil,
		},
		{
			name:        "valid link",
			html:        htmlStr(`<a href="/test">Test</a>`),
			expectedErr: nil,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			links, _ := tokenize(strings.NewReader(sc.html))
			fmt.Printf("links %v", links)
		})
	}
}

func checkErr(t *testing.T, actualErr, expectedErr error) {
	t.Helper()
	if actualErr != nil && expectedErr == nil {
		t.Fatalf("unexpected error = %v", actualErr)
	}
	if actualErr == nil && expectedErr != nil {
		t.Fatalf("want error = %v, but there was none", expectedErr)
	}
	if actualErr != nil && expectedErr != nil {
		if actualErr.Error() != expectedErr.Error() {
			t.Fatalf("error == %s, want %s",
				actualErr.Error(), expectedErr.Error())
		}
	}
}

func htmlStr(tag string) string {
	return `<!DOCTYPE html>
<html>
<head>
<title>Page Title</title>
</head>
<body>

	<h1>My First Heading</h1>
	<p>My first paragraph.</p>
	` + tag + `
</body>
</html>`
}
