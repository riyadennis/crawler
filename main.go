package main

import (
	"flag"
	"fmt"
)

func main(){
	url := flag.String("root", "https://google.co.uk", "root url")
	flag.Parse()
	fmt.Printf("RRR %v", url)
}
