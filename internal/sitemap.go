package internal

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
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
	r, err := fetchData(c.RootURL)
	if err != nil {
		return nil, err
	}
	return tokenize(r)
}

func tokenize(reader io.Reader) (map[int]string, error) {
	links := make(map[int]string)
	token := html.NewTokenizer(reader)
	for {
		if token.Err() == io.EOF {
			break
		}
		if token.Err() != nil {
			return nil, token.Err()
		}
		tokenType := token.Next()
		switch tokenType {
		case html.StartTagToken:
			t := token.Token()
			links = searchLinks(t)
		}
	}
	return links, nil
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

func fetchData(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	body := resp.Body
	if err != nil {
		return nil, err
	}
	// check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unable to load the url")
	}
	//check response content type
	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		log.Fatalf("response content type was %s not text/html\n", ctype)
	}
	return body, nil
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
