package main

import (
	"flag"

	"github.com/riyadennis/crawler/internal"
)

func main(){
	url := flag.String("root", "google.co.uk", "root url")
	flag.Parse()
	internal.SiteMap(*url)
}
