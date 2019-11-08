package internal

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"strings"
	"testing"
)

func TestSiteMap(t *testing.T) {
	var scenarios = []struct {
		name        string
		url         string
		expectedErr error
	}{
		//{
		//	name: "invalid host name",
		//	url:  "invalid",
		//	expectedErr: errors.New("invalid host name invalid"),
		//},
		{
			name: "valid host",
			url:  "google.co.uk",
			expectedErr: nil,
		},
		//{
		//	name: "valid host",
		//	url:  "mail.google.com/mail/u/0/#inbox",
		//	expectedErr: nil,
		//},
	}
	for _, sc := range scenarios{
		t.Run(sc.name, func(t *testing.T){
			err := SiteMap(sc.url)
			checkErr(t, err, sc.expectedErr)
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

func testHTMLParsing(){
	s := `<p>Links:</p><ul><li><a href="foo">Foo</a><li><a href="/bar/baz">BarBaz</a></ul>`
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					fmt.Println(a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}