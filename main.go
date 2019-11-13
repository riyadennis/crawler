package main

import (
	"context"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go webCrawler.Crawl(ctx, *rootURL, *depth,0, ch)
	webCrawler.Display(ctx, *rootURL, ch)

	close(ch)
}
