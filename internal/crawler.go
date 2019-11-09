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
	"sync"
)

const regExpDomain = `^([a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,}$`

// Crawler holds data that we need to parse a web page
type Crawler struct {
	RootURL *url.URL
	Fetcher func() (io.ReadCloser, error)
	Parser  func(reader io.ReadCloser) map[int]string
}

func Crawl(url string, ch chan map[int]string, wg *sync.WaitGroup) {
	c, err := NewCrawler(url)
	if err != nil {
		fmt.Printf("failed to create crawler :: %v", err)
	}
	r, err := c.Fetcher()
	if err != nil {
		fmt.Printf("failed to fetch:: %v", err)
	}
	defer r.Close()
	ch <- c.Parser(r)
	wg.Done()
}

// NewCrawler initialises the Crawler
func NewCrawler(url string) (*Crawler, error) {
	c := &Crawler{}
	u, err := validateURL(url)
	if err != nil {
		return nil, err
	}
	c.RootURL = u
	c.Fetcher = c.fetchData
	c.Parser = c.tokenize
	return c, nil
}

func validateURL(rootURL string) (*url.URL, error) {
	url, err := url.Parse(rootURL)
	if err != nil {
		return nil, err
	}
	reg, err := regexp.Compile(regExpDomain)
	if err != nil {
		return nil, err
	}
	if !reg.MatchString(url.Host) {
		if url.Host == "" {
			return nil, errors.New("empty host name")
		}
		return nil, fmt.Errorf("invalid host name %s", url.Host)
	}
	return url, nil
}

func (c *Crawler) fetchData() (io.ReadCloser, error) {
	resp, err := http.Get(c.RootURL.String())
	if err != nil {
		return nil, err
	}
	body := resp.Body
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