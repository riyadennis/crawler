package internal

import (
	"errors"
	"reflect"
	"testing"

	"golang.org/x/net/html"
)

func TestSiteMap(t *testing.T) {
	var scenarios = []struct {
		name        string
		url         string
		expectedErr error
	}{
		{
			name:        "invalid host name",
			url:         "invalid",
			expectedErr: errors.New("invalid host name invalid"),
		},
		{
			name:        "valid host",
			url:         "google.co.uk",
			expectedErr: nil,
		},
		{
			name:        "valid host",
			url:         "mail.google.com/mail/u/0/#inbox",
			expectedErr: nil,
		},
	}
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			c := NewCrawler(sc.url)
			err := c.Map()
			checkErr(t, err, sc.expectedErr)
		})
	}
}
func TestLinks(t *testing.T) {
	scenarios := []struct{
		name string
		node *html.Node
		linksExp map[int]string
	}{
		{
			name: "empty node",
			node: nil,
			linksExp: map[int]string{},
		},
		{
			name: "empty node",
			node: &html.Node{
				Type: html.ElementNode,
				Data: "a",
				Attr: []html.Attribute{
					{
						Key: "href",
						Val: "test",
					},
				},
			},
			linksExp: map[int]string{ 1: "test"},
		},
	}
	for _, sc := range scenarios{
		t.Run(sc.name, func(t *testing.T){
			l := links(sc.node)
			if reflect.DeepEqual(l, sc.linksExp){
				t.Errorf("links got %v, want %v", l, sc.linksExp)
			}
		})
	}
}

func checkErr(t *testing.T, actualErr, expectedErr error) {
	t.Helper()
	if actualErr != nil && expectedErr == nil {
		t.Fatalf("unexpected error = %v", actualErr)
	}
	if actualErr == nil && expectedErr != nil {
		t.Fatalf("want error = %v, but there was none", expectedErr)
	}
	if actualErr != nil && expectedErr != nil {
		if actualErr.Error() != expectedErr.Error() {
			t.Fatalf("error == %s, want %s",
				actualErr.Error(), expectedErr.Error())
		}
	}
}
