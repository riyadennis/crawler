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

	"github.com/sirupsen/logrus"
)

const regExpDomain = `^([a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,}$`

// webCrawler holds data that we need to parse a web page
type webCrawler struct {
	Content func(url string) (io.ReadCloser, error)
	SiteMap func(url string, reader io.ReadCloser) map[int]string
}

func (c *webCrawler) extractLinks(url string) map[int]string {
	if c == nil {
		logrus.Errorf("failed to extract links, empty crawler")
		return nil
	}
	if c.Content == nil {
		logrus.Errorf("method to extract content not set")
		return nil
	}
	r, err := c.Content(url)
	if err != nil {
		logrus.Errorf("failed to fetch:: %v", err)
		return nil
	}
	defer r.Close()
	if c.SiteMap == nil {
		logrus.Errorf("failed to create site map:: %v", err)
		return nil
	}
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

func fileContent(name string) (io.ReadCloser, error) {
	cnt, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(bytes.NewReader(cnt)), nil
}
