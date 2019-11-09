package main

import (
	"flag"
	"fmt"

	"github.com/riyadennis/crawler/internal"
)

func main() {
	url := flag.String("root", "google.co.uk", "root url")
	flag.Parse()
	c, err := internal.NewCrawler(*url)
	if err != nil {
		panic(err)
	}
	r, err := c.Fetcher()
	if err != nil {
		panic(err)
	}
	links := c.Parser(r)
	fmt.Printf("got links %v", links)
}
