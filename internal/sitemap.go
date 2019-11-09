package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"golang.org/x/net/html"
)

const regExpDomain = `^([a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,}$`

// Crawler holds data that we need to parse a web page
type Crawler struct {
	RootURL string
	Links   func(*html.Node) []string
}

// NewCrawler initialises the Crawler
func NewCrawler(url string) *Crawler {
	c := &Crawler{}
	c.RootURL = url
	c.Links = links
	return c
}

// Map will try to create a site map of links
func (c *Crawler) Map() error {
	u, err := validateURL(c.RootURL)
	if err != nil {
		return err
	}
	re, err := fetchURL(u)
	if err != nil {
		return err
	}
	if re == nil {
		return errors.New("unable to load data from url")
	}
	doc, err := html.Parse(bytes.NewReader(re))
	if err != nil {
		log.Fatal(err)
	}
	l := c.Links(doc)
	fmt.Printf("%v", l)
	return nil
}

func links(node *html.Node) []string {
	var l []string
	if node.Type != html.ElementNode && node.Data != "a" {
		return l
	}
	for _, a := range node.Attr {
		if a.Key == "href" {
			l = append(l, a.Val)
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		l = links(child)
	}
	return l
}

func fetchURL(u *url.URL) ([]byte, error) {
	resp, err := http.Get(u.String())
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unable to load the url")
	}
	return ioutil.ReadAll(resp.Body)
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
