package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/riyadennis/crawler/internal"
)

func main() {
	rootURL := flag.String("root", "https://google.co.uk", "root ur2l")
	depth := flag.Int("depth", 3, "depth for crawling")
	stats := flag.Bool("stats", false, "show memory stats")

	flag.Parse()

	var statsB, statsE *internal.MemStats
	if *stats {
		statsB = internal.GetMemStats()
	}

	webCrawler, err := internal.NewCrawler(*rootURL)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	webCrawler.Display(ctx, *rootURL, *depth, webCrawler.Crawl(ctx, *rootURL, *depth))

	if *stats {
		statsE = internal.GetMemStats()
		fmt.Printf("difference in alloc %d\n", statsB.Alloc-statsE.Alloc)
		fmt.Printf("difference in TotalAlloc %d\n", statsB.TotalAlloc-statsE.TotalAlloc)
	}
}
