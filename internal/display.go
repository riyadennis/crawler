package internal

import "fmt"

//Display will listen to the channel and print results into  console
func Display(ch chan map[int]map[int]string) {
	for {
		select {
		case dlinks := <-ch:
			for _, dl := range dlinks {
				for _, l := range dl {
					fmt.Printf("\n from channel %s \n", l)
				}
			}
		}
	}
}
