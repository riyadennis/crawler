package internal

import (
	"fmt"

	"github.com/disiqueira/gotree"
	"golang.org/x/net/context"
)

type Crawler interface {
	Crawl(ctx context.Context, source string, depth, index int, ch chan map[int]map[int]string)
	Display(ctx context.Context, source string, depth int, ch <-chan map[int]map[int]string)
}

// NewCrawler initialises the Crawler to search for links in a web page
func NewCrawler(url string) (Crawler, error) {
	c := &webCrawler{}
	err := validateURL(url)
	if err != nil {
		return nil, err
	}

	c.Content = content
	c.SiteMap = siteMap
	return c, nil
}

//Crawl does the scrapping of links and sub links
func (c *webCrawler) Crawl(ctx context.Context, source string,
	depth, index int, ch chan map[int]map[int]string) {
	if depth <= 0 {
		return
	}
	links := make(map[int]map[int]string)
	link := c.extractLinks(source)
	if link == nil {
		return
	}
	links[index] = link
	ch <- links
	index++
	for _, li := range link {
		c.Crawl(ctx, li, depth, index, ch)
	}
	close(ch)
	<-ctx.Done()
}

//Display will listen to the channel and print results into  console
func (c *webCrawler) Display(ctx context.Context, source string,
	depth int, ch <-chan map[int]map[int]string) {
	ch1 := make(chan map[int]map[int]string)
	go func() {
		for dl := range ch {
			ch1 <- dl
		}
		close(ch1)
	}()
	tree := gotree.New(source)
	for i := 0; i < depth; i++ {
		for i, dl := range <-ch {
			child := tree.Add(dl[i])
			for _, dl := range dl {
				child.Add(dl)
			}
		}
	}
	fmt.Println(tree.Print())
}
