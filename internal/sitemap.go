package internal

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/url"
	"strings"
)

func (c *Crawler) siteMap(reader io.ReadCloser) map[int]string {
	links := make(map[int]string)
	token := html.NewTokenizer(reader)
	defer reader.Close()
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
			link := searchLinks(t, c.Source.Host)
			if link != "" {
				if checkDomain(c.Source.Host, link) {
					links[i] = link
					i++
				}
			}
		}
	}
	return links
}

func searchLinks(t html.Token, hostname string) string {
	if t.Data == "a" {
		for _, att := range t.Attr {
			if att.Key == "href" {
				if att.Val != "#" {
					// if its an internal link we need to append full path
					if strings.HasPrefix(att.Val, "/") {
						return fmt.Sprintf("https://%s%s", hostname, att.Val)
					}
					return att.Val
				}
			}
		}
	}
	return ""
}

func checkDomain(hostname, link string) bool {
	l, err := url.Parse(link)
	if err != nil {
		return false
	}
	parentDomain := strings.TrimPrefix(hostname, "www.")
	linkDomain := strings.TrimPrefix(l.Host, "www.")
	if parentDomain == linkDomain {
		return true
	}
	return false
}
