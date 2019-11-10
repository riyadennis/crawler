package internal

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const regExpDomain = `^([a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,}$`

type Crawler interface {
	Crawl(source string, depth int, ch chan map[int]map[int]string)
}

// webCrawler holds data that we need to parse a web page
type webCrawler struct {
	Fetcher func(source string) (io.ReadCloser, error)
	Parser  func(source string, reader io.ReadCloser) map[int]string
}

//Crawl does the scrapping of links and sub links
func (c *webCrawler) Crawl(source string, depth int, ch chan map[int]map[int]string) {
	if depth <= 0 {
		return
	}
	links := make(map[int]map[int]string)

	links[depth] = linksFrmURL(source)
	ch <- links
	for _, l := range links {
		for _, li := range l {
			c.Crawl(li, depth-1, ch)
		}
	}
}

func linksFrmURL(url string) map[int]string {
	c, err := NewWebCrawler(url)
	if err != nil {
		fmt.Printf("failed to create crawler :: %v", err)
	}
	r, err := c.Fetcher(url)
	if err != nil {
		fmt.Printf("failed to fetch:: %v", err)
	}
	defer r.Close()
	return c.Parser(url, r)
}

// NewWebCrawler initialises the webCrawler to search for links in a webpage
func NewWebCrawler(url string) (*webCrawler, error) {
	c := &webCrawler{}
	err := validateURL(url)
	if err != nil {
		return nil, err
	}
	c.Fetcher = c.fetcher
	c.Parser = c.parser
	return c, nil
}

func validateURL(rootURL string) error {
	url, err := url.Parse(rootURL)
	if err != nil {
		return err
	}
	reg, err := regexp.Compile(regExpDomain)
	if err != nil {
		return err
	}
	if !reg.MatchString(url.Host) {
		if url.Host == "" {
			return errors.New("empty host name")
		}
		return fmt.Errorf("invalid host name %s", url.Host)
	}
	return nil
}

func (c *webCrawler) fetcher(source string) (io.ReadCloser, error) {
	resp, err := http.Get(source)
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
		return nil, fmt.Errorf("response content type was %s not text/html\n", ctype)
	}
	return body, nil
}
