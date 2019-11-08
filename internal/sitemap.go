package internal

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

const regExpDomain  = `^([a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,}$`

func SiteMap(rootURL string) error{
	u, err := validateURL(rootURL)
	if err != nil{
		return err
	}
	body, err := fetchURL(u)
	if err!=nil{
		return err
	}
	if body == nil{
		return errors.New("unable to load data from url")
	}
	return nil
}

func fetchURL(u *url.URL) ([]byte, error){
	resp, err := http.Get(u.String())
	defer resp.Body.Close()
	if err != nil{
		return nil, err
	}
	if resp.StatusCode != http.StatusOK{
		return nil, errors.New("unable to load the url")
	}
	return  ioutil.ReadAll(resp.Body)
}

func validateURL(rootURL string) (*url.URL, error){
	rootURL = fmt.Sprintf("%s://%s", "http", rootURL)
	url, err := url.Parse(rootURL)
	if err != nil{
		return nil, err
	}
	reg, err := regexp.Compile(regExpDomain)
	if err != nil{
		return nil,err
	}
	if !reg.MatchString(url.Host){
		return nil, fmt.Errorf("invalid host name %s", url.Host)
	}
	return url, nil
}