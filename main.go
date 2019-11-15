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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	webCrawler.Display(ctx, *rootURL, *depth,
		webCrawler.Crawl(ctx, *rootURL, *depth, 0))
}
