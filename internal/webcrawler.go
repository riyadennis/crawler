package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const regExpDomain = `^([a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,}$`

// webCrawler holds data that we need to parse a web page
type webCrawler struct {
	Content func(url string) (io.ReadCloser, error)
	SiteMap func(url,topic string, reader io.ReadCloser) (map[int]string, error)
	Topic string
}

func (c *webCrawler) extractLinks(url string) (map[int]string, error) {
	if c.Content == nil {
		return nil, errors.New("empty content")
	}

	r, err := c.Content(url)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	if c.SiteMap == nil {
		return nil, errors.New("no links in the page")
	}

	return c.SiteMap(url,c.Topic,  r)
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

	if url.Host == "" {
		return errors.New("empty host name")
	}

	if !reg.MatchString(url.Host) {
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

	// check response content type
	cType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(cType, "text/html") {
		return nil, fmt.Errorf("response content type was %s not text/html\n", cType)
	}

	return body, nil
}

func fileContent(name string) (io.ReadCloser, error) {
	cnt, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	return ioutil.NopCloser(bytes.NewReader(cnt)), nil
}
