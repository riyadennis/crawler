package main

import (
	"flag"
	"fmt"
	"github.com/riyadennis/crawler/internal"
	"sync"
)

func main() {
	var wt sync.WaitGroup
	rootURL := flag.String("root", "https://monzo.com", "root ur2l")
	flag.Parse()

	wt.Add(1)
	ch := make(chan map[int]string)
	go internal.Crawl(*rootURL, ch, &wt)
	fmt.Printf("%v", <-ch)
	wt.Wait()

	defer close(ch)
}
