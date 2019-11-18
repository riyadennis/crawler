package internal

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLinksFromURL(t *testing.T) {
	scenarios := []struct {
		name    string
		crawler *webCrawler
		url     string
		links   map[int]string
	}{
		{
			name:    "empty crawler",
			crawler: nil,
			url:     "test",
			links:   nil,
		},
		{
			name:    "empty crawler",
			crawler: &webCrawler{},
			url:     "test",
			links:   nil,
		},
		{
			name:    "empty crawler",
			crawler: &webCrawler{Content: nil},
			url:     "test",
			links:   nil,
		},
		{
			name:    "invalid url",
			crawler: &webCrawler{Content: content},
			url:     "test",
			links:   nil,
		},
		{
			name:    "empty site map",
			crawler: &webCrawler{Content: content},
			url:     "http://google.co.uk",
			links:   nil,
		},
		{
			name:    "empty site map",
			crawler: &webCrawler{Content: fileContent, SiteMap: siteMap},
			url:     "testdata/sample.html",
			links: map[int]string{
				0: "test",
			},
		},
	}
	ctx := context.TODO()
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			links := sc.crawler.extractLinks(ctx, sc.url)
			if !cmp.Equal(links, sc.links) {
				t.Errorf("unexpected links,got %v, want %v", links, sc.links)
			}
		})
	}
}
