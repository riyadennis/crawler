package internal

import (
	"golang.org/x/net/html"
	"io"
)

func tokenize(reader io.Reader) map[int]string {
	links := make(map[int]string)
	token := html.NewTokenizer(reader)
	i := 0
	for {
		if token.Err() == io.EOF {
			break
		}
		if token.Err() != nil {
			return nil
		}
		tokenType := token.Next()
		switch tokenType {
		case html.StartTagToken:
			t := token.Token()
			link := searchLinks(t)
			if link != "" {
				links[i] = link
				i++
			}
		}
	}
	return links
}

func searchLinks(t html.Token) string {
	if t.Data == "a" {
		for _, att := range t.Attr {
			if att.Key == "href" {
				if att.Val != "#" {
					return att.Val
				}
			}
		}
	}
	return ""
}
