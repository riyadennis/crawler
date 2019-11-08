package internal

import (
	"errors"
	"fmt"
	"net/http"
)

func SiteMap(url string) error{
	resp, err := http.Get(url)
	if err != nil{
		return err
	}
	if resp.StatusCode != http.StatusOK{
		return errors.New("expected")
	}
	fmt.Printf("url from site map %s", url)
	return nil
}
