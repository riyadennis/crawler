package internal

import (
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
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
			name:  "valid link",
			html:  htmlStr(`<a href="/test">Test</a>`),
			links: map[int]string{0: "/test"},
		},

		{
			name: "valid two links",
			html: htmlStr(`<a href="/test1">Test</a>
<a href="/test2">Test</a>`),
			links: map[int]string{0: "/test1", 1: "/test2"},
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			links := tokenize(strings.NewReader(sc.html))
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
