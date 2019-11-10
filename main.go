package main

import (
	"flag"
	"github.com/riyadennis/crawler/internal"
)

func main() {
	rootURL := flag.String("root", "https://monzo.com", "root ur2l")
	depth := flag.Int("depth", 2, "depth for crawling")
	flag.Parse()
	webCrawler, err := internal.NewWebCrawler(*rootURL)
	if err != nil {
		panic(err)
	}
	ch := make(chan map[int]map[int]string, *depth)
	go webCrawler.Crawl(*rootURL, *depth, ch)
	internal.Display(ch)
	close(ch)
}
