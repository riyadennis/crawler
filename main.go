package main

import (
	"flag"
	"github.com/riyadennis/crawler/internal"
)

func main() {
	rootURL := flag.String("root", "https://monzo.com", "root ur2l")
	flag.Parse()
	depth := 2

	ch := make(chan map[int]map[int]string, depth)
	go internal.Crawl(*rootURL, depth, 0, ch)
	internal.Display(ch)
	close(ch)
}
