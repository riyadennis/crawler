package internal

import (
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
	for {
		if token.Err() == io.EOF {
			break
		}
		if token.Err() != nil {
			//TODO handle error properly
			return nil
		}
		tokenType := token.Next()
		switch tokenType {
		case html.StartTagToken:
			t := token.Token()
			select {
			case link := <-searchLinks(t, u.Host):
				if link != "" {
					if checkDomain(u.Host, link) {
						links[i] = link
						i++
					}
				}
			default:
				break
			}
		}
	}
	return links
}

func searchLinks(t html.Token, hostname string) <-chan string {
	ch := make(chan string)
	var wg sync.WaitGroup
	if t.Data == "a" {
		wg.Add(1)
		go func() {
			for _, att := range t.Attr {
				if att.Key == "href" {
					if att.Val != "#" {
						// if its an internal link we need to append full path
						if strings.HasPrefix(att.Val, "/") {
							ch <- fmt.Sprintf("https://%s%s", hostname, att.Val)
						}
					}
				}
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
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
