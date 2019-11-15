package internal

import (
	"fmt"
	"github.com/disiqueira/gotree"
	"golang.org/x/net/context"
	"sync"
)

type Crawler interface {
	Crawl(ctx context.Context, source string, depth, index int)<-chan map[int]map[int]string
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
	ch := make(chan map[int]map[int]string, depth)
	if depth <= 0 {
		return ch
	}
	links := make(map[int]map[int]string)
	link := c.extractLinks(ctx, source)
	if link == nil {
		return ch
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func(){
		for i, li := range link {
			links[i] = c.extractLinks(ctx, li)
			if i > depth{
				break
			}
		}
		ch <- links
		wg.Done()
	}()

	go func(){
		wg.Wait()
		close(ch)
	}()
	return ch
}

//Display will listen to the channel and print results into  console
func (c *webCrawler) Display(ctx context.Context, source string,
	depth int, ch <-chan map[int]map[int]string) {
	ch1 := make(chan map[int]map[int]string, depth)
	go func() {
		defer close(ch1)
		for dl := range ch {
			select{
				case ch1 <- dl:
				case <-ctx.Done():
			}
		}
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
	<-ctx.Done()
}
