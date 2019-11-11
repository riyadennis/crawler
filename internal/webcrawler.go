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

// webCrawler holds data that we need to parse a web page
type webCrawler struct {
	Content func(url string) (io.ReadCloser, error)
	SiteMap func(url string, reader io.ReadCloser) map[int]string
}

func (c *webCrawler) linksFrmURL(url string) map[int]string {
	r, err := c.Content(url)
	if err != nil {
		fmt.Printf("failed to fetch:: %v", err)
	}
	defer r.Close()
	return c.SiteMap(url, r)
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

func content(source string) (io.ReadCloser, error) {
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
