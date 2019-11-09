package internal

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTokenize(t *testing.T) {
	scenarios := []struct {
		name  string
		html  string
		links map[int]string
	}{
		{
			name:  "not html",
			html:  "not html",
			links: map[int]string{},
		},
		{
			name:  "no links",
			html:  htmlStr("<p>hello</p>"),
			links: map[int]string{},
		},
		{
			name:  "dead link",
			html:  htmlStr(`<a>Test</a>`),
			links: map[int]string{},
		},
		{
			name:  "invalid href link",
			html:  htmlStr(`<a href="#">Test</a>`),
			links: map[int]string{},
		},
		{
			name: "valid link",
			html: htmlStr(`<a href="/test">Test</a>`),
			links: map[int]string{
				0: "https://www.google.co.uk/test",
			},
		},

		{
			name: "valid two links",
			html: htmlStr(`<a href="/test1">Test</a>
		<a href="/test2">Test</a>`),
			links: map[int]string{
				0: "https://www.google.co.uk/test1",
				1: "https://www.google.co.uk/test2"},
		},
		{
			name: "valid two links",
			html: htmlStr(`<a href="http://www.google.co.uk/imghp?hl=en&tab=w">Test</a>
<a href="http://www.maps.co.uk/imghp?hl=en&tab=w">Test</a>`),
			links: map[int]string{0: "http://www.google.co.uk/imghp?hl=en&tab=w"},
		},
	}
	c, err := NewCrawler("http://www.google.co.uk")
	if err != nil {
		t.Fatal(err)
	}
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			links := c.tokenize(ioutil.NopCloser(strings.NewReader(sc.html)))
			if !cmp.Equal(links, sc.links) {
				t.Errorf("got %v, want %v", links, sc.links)
			}
		})
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
