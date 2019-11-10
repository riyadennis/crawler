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
	Source  *url.URL
	Fetcher func() (io.ReadCloser, error)
	Parser  func(reader io.ReadCloser) map[int]string
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
	r, err := c.Fetcher()
	if err != nil {
		fmt.Printf("failed to fetch:: %v", err)
	}
	defer r.Close()
	return c.Parser(r)
}

// NewWebCrawler initialises the webCrawler to search for links in a webpage
func NewWebCrawler(url string) (*webCrawler, error) {
	c := &webCrawler{}
	u, err := validateURL(url)
	if err != nil {
		return nil, err
	}
	c.Source = u
	c.Fetcher = c.readURL
	c.Parser = c.siteMap
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

func (c *webCrawler) readURL() (io.ReadCloser, error) {
	resp, err := http.Get(c.Source.String())
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
