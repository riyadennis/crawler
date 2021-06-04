package internal

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func siteMap(rootURL string, reader io.ReadCloser) map[int]string {
	token := html.NewTokenizer(reader)
	defer reader.Close()

	i := 0
	u, err := url.Parse(rootURL)
	if err != nil {
		return nil
	}

	links := make(map[int]string)
	if token.Err() != nil {
		//TODO handle error properly
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for {
			if errors.Is(token.Err(), io.EOF) {
				break
			}
			tokenType := token.Next()
			switch tokenType {
			case html.StartTagToken:
				t := token.Token()
				link := searchLinks(t, u.Host)
				if link != "" {
					if checkDomain(u.Host, link) {
						links[i] = link
						i++
					}
				}
			}
		}
		wg.Done()
	}()

	wg.Wait()
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

	// compare parent and link domains
	return strings.TrimPrefix(hostname, "www.") == strings.TrimPrefix(l.Host, "www.")
}
