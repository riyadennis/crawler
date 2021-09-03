package internal

import (
	"errors"
	"fmt"
	"sync"

	"github.com/disiqueira/gotree"
	"golang.org/x/net/context"
)

type Crawler interface {
	Crawl(context.Context, string, int) (<-chan map[int]map[int]string, error)
	Display(context.Context, string, int, <-chan map[int]map[int]string)
}

// NewCrawler initialises the Crawler to search for links in a web page
func NewCrawler(url, topic string) (Crawler, error) {
	err := validateURL(url)
	if err != nil {
		return nil, err
	}

	return &webCrawler{
		Content: content,
		SiteMap: siteMap,
		Topic:   topic,
	}, nil
}

// Crawl does the scrapping of links and sub links
func (c *webCrawler) Crawl(ctx context.Context, source string, depth int) (<-chan map[int]map[int]string, error) {
	if depth <= 0 {
		return nil, errors.New("invalid depth argument")
	}

	link, err := c.extractLinks(source)
	if link == nil {
		return nil, errors.New("no links in the page")
	}

	ch := make(chan map[int]map[int]string, depth)
	links := make(map[int]map[int]string)

	go func() {
		for i, li := range link {
			if len(li) > 0 {
				links[i], err = c.extractLinks(li)
			}
		}
		ch <- links
	}()

	go func() {
		<-ctx.Done()
		close(ch)
	}()

	return ch, nil
}

// Display will listen to the channel and print results into  console
func (c *webCrawler) Display(ctx context.Context, source string, depth int, ch <-chan map[int]map[int]string) {
	tree := gotree.New(source)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		select {
		case links := <-ch:
			for i, dl := range links {
				if i < depth {
					child := tree.Add(dl[i])
					for _, dl := range dl {
						if len(dl) > 0 {
							child.Add(dl)
						}
					}
				}
			}
			fmt.Println(tree.Print())
		case <-ctx.Done():
			wg.Done()
			return
		}
	}()

	wg.Wait()
}
