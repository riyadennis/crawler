package internal

import "fmt"

type Crawler interface {
	Crawl(source string, depth int, ch chan map[int]map[int]string)
	Display(ch chan map[int]map[int]string)
}

//Crawl does the scrapping of links and sub links
func (c *webCrawler) Crawl(source string, depth int, ch chan map[int]map[int]string) {
	if depth <= 0 {
		return
	}
	links := make(map[int]map[int]string)

	links[depth] = c.linksFrmURL(source)
	ch <- links
	for _, l := range links {
		for _, li := range l {
			c.Crawl(li, depth-1, ch)
		}
	}
}

//Display will listen to the channel and print results into  console
func (c *webCrawler) Display(ch chan map[int]map[int]string) {
	for {
		select {
		case dlinks := <-ch:
			for _, dl := range dlinks {
				for _, l := range dl {
					fmt.Printf("\n from channel %s \n", l)
				}
			}
		}
	}
}
