package internal

import (
	"fmt"
	"github.com/disiqueira/gotree"
	"golang.org/x/net/context"
)

type Crawler interface {
	Crawl(ctx context.Context, source string, depth, index int, ch chan map[int]map[int]string)
	Display(ctx context.Context, source string,ch chan map[int]map[int]string)
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
	depth,index int, ch chan map[int]map[int]string) {
	if depth <= 0 {
		return
	}
	links := make(map[int]map[int]string)
	link := c.extractLinks(source)
	if link == nil{
		return
	}
	links[index] = link
	ch<-links
	index++
	for _, li := range link {
		c.Crawl(ctx, li, depth-1, index, ch)
	}

	<-ctx.Done()
}

//Display will listen to the channel and print results into  console
func (c *webCrawler) Display(ctx context.Context, source string, ch chan map[int]map[int]string) {
	tree := gotree.New(source)
	for {
		select {
		case dlinks := <-ch:
			for i, dl := range dlinks {
				child := tree.Add(dl[i])
				for _, l := range dl {
					child.Add(l)
				}
				fmt.Println(tree.Print())
			}
		case <-ctx.Done():
			return
		}
	}

}
