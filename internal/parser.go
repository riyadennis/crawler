package internal

import (
	"golang.org/x/net/html"
	"io"
)

func tokenize(reader io.Reader) map[int]string {
	links := make(map[int]string)
	token := html.NewTokenizer(reader)
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
			links = searchLinks(t)
		}
	}
	return links
}

func searchLinks(t html.Token) map[int]string {
	links := make(map[int]string)
	i := 0
	if t.Data == "a" {
		for _, att := range t.Attr {
			if att.Key == "href" {
				if att.Val != "#" {
					links[i] = att.Val
					i++
				}
			}
		}
	}
	return links
}
