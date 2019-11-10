package main

import (
	"flag"
	"github.com/riyadennis/crawler/internal"
)

func main() {
	rootURL := flag.String("root", "https://monzo.com", "root ur2l")
	flag.Parse()
	depth := 3

	ch := make(chan map[int]map[int]string, depth)
	go internal.Crawl(*rootURL, depth, ch)
	internal.Display(ch)
	close(ch)
}
