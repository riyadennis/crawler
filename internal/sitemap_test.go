package internal

import (
	"errors"
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
	node := &html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			{
				Key: "href",
				Val: "test",
			},
		},
	}
	l := links(node)
	t.Logf("links %v", l)
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
