package main

import (
	"flag"
	"fmt"
	"github.com/riyadennis/crawler/internal"
)

func main() {
	ur2l := flag.String("root", "http://monzo.com", "root ur2l")
	flag.Parse()
	var f func(string) map[int]string
	f = func(url string) map[int]string {
		links := make(map[int]string)
		c, err := internal.NewCrawler(url)
		if err != nil {
			fmt.Printf("failed to create crawler :: %v", err)
		}
		r, err := c.Fetcher()
		if err != nil {
			fmt.Printf("failed to fetch:: %v", err)
		}
		links = c.Parser(r)
		return links
	}
	fmt.Printf("inside %v", f(*ur2l))
}
