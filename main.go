package main

import (
	"flag"
	"github.com/riyadennis/crawler/internal"
)

func main() {
	rootURL := flag.String("root", "https://google.com", "root ur2l")
	depth := flag.Int("depth", 3, "depth for crawling")
	flag.Parse()

	webCrawler, err := internal.NewCrawler(*rootURL)
	if err != nil {
		panic(err)
	}
	ch := make(chan map[int]map[int]string, *depth)
	go webCrawler.Crawl(*rootURL, *depth,0, ch)
	webCrawler.Display(*rootURL, ch)

	close(ch)
}
