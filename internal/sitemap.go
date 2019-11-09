package internal

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"regexp"
)

const regExpDomain = `^([a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,}$`

// Crawler holds data that we need to parse a web page
type Crawler struct {
	RootURL string
}

// NewCrawler initialises the Crawler
func NewCrawler(url string) (*Crawler, error) {
	c := &Crawler{}
	u, err := validateURL(url)
	if err != nil {
		return nil, err
	}
	c.RootURL = u.String()
	return c, nil
}

// Map will try to create a site map of links
func (c *Crawler) Crawl() (map[int]string, error) {
	resp, err := http.Get(c.RootURL)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unable to load the url")
	}
	links := make(map[int]string)
	token := html.NewTokenizer(resp.Body)
	i := 0
	for {
		tokenType := token.Next()
		switch tokenType {
		case html.ErrorToken:
		case html.StartTagToken:
			t := token.Token()
			if t.Data == "a" {
				for _, att := range t.Attr {
					if att.Key == "href" {
						links[i] = att.Val
					}
				}

			}

		}
	}
	return links, nil
}

func validateURL(rootURL string) (*url.URL, error) {
	rootURL = fmt.Sprintf("%s://%s", "http", rootURL)
	url, err := url.Parse(rootURL)
	if err != nil {
		return nil, err
	}
	reg, err := regexp.Compile(regExpDomain)
	if err != nil {
		return nil, err
	}
	if !reg.MatchString(url.Host) {
		return nil, fmt.Errorf("invalid host name %s", url.Host)
	}
	return url, nil
}
