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
	topic := flag.String("topic", "crawler", "topic to add links to")

	flag.Parse()

	var statsB, statsE *internal.MemStats
	if *stats {
		statsB = internal.GetMemStats()
	}

	webCrawler, err := internal.NewWebCrawler(*rootURL, *topic)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	links, err := webCrawler.Crawl(ctx, *rootURL, *depth)
	if err != nil {
		panic(err)
	}

	webCrawler.Display(ctx, *rootURL, *depth, links)


	if *stats {
		statsE = internal.GetMemStats()
		fmt.Printf("difference in alloc %d\n", statsB.Alloc-statsE.Alloc)
		fmt.Printf("difference in TotalAlloc %d\n", statsB.TotalAlloc-statsE.TotalAlloc)
	}

	go func(){
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil{
				panic(err)
			}
		case err := <-webCrawler.ErrChan:
			if err != nil{
				panic(err)
			}
		}
	}()
}
