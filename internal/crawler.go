package internal

import (
	"fmt"
	"github.com/disiqueira/gotree"
)

type Crawler interface {
	Crawl(source string, depth int, ch chan map[int]map[int]string)
	Display(source string,ch chan map[int]map[int]string)
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
func (c *webCrawler) Crawl(source string, depth int, ch chan map[int]map[int]string) {
	if depth <= 0 {
		return
	}
	links := make(map[int]map[int]string)

	links[depth] = c.extractLinks(source)
	ch <- links
	for _, l := range links {
		for _, li := range l {
			c.Crawl(li, depth-1, ch)
		}
	}
}

//Display will listen to the channel and print results into  console
func (c *webCrawler) Display(source string, ch chan map[int]map[int]string) {
	artist := gotree.New(source)
	for {
		select {
		case dlinks := <-ch:
			go func(){
				for _, dl := range dlinks {
					child := artist.Add(dl[0])
					for _, l := range dl {
						child.Add(l)
						//fmt.Printf("\n page %d: %s \n", i, l)
					}
				}
				fmt.Println(artist.Print())
			}()
		}
	}

}
