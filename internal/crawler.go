package internal

import (
	"fmt"
	"sync"

	"github.com/disiqueira/gotree"
	"golang.org/x/net/context"
)

type Crawler interface {
	Crawl(ctx context.Context, source string, depth, index int) <-chan map[int]map[int]string
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
func (c *webCrawler) Crawl(ctx context.Context,
	source string, depth, index int) <-chan map[int]map[int]string {
	if depth <= 0 {
		return nil
	}
	link := c.extractLinks(ctx, source)
	if link == nil {
		return nil
	}
	ch := make(chan map[int]map[int]string, depth)
	links := make(map[int]map[int]string)
	go func() {
		for i, li := range link {
			if len(li) > 0 {
				links[i] = c.extractLinks(ctx, li)
			}
		}
		ch <- links
	}()

	go func() {
		<-ctx.Done()
		close(ch)
	}()

	return ch
}

//Display will listen to the channel and print results into  console
func (c *webCrawler) Display(ctx context.Context, source string,
	depth int, ch <-chan map[int]map[int]string) {
	tree := gotree.New(source)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(){
		select {
		case links := <-ch:
			for i, dl := range links {
				if i < depth{
					child := tree.Add(dl[i])
					for _, dl := range dl {
						if len(dl) > 0{
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
