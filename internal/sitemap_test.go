package internal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSiteMap(t *testing.T) {
	var scenarios = []struct {
		name        string
		url         string
		expectedErr error
	}{
		//{
		//	name:        "invalid host name",
		//	url:         "invalid",
		//	expectedErr: errors.New("invalid host name invalid"),
		//},
		{
			name:        "valid host",
			url:         "monzo.com",
			expectedErr: nil,
		},
		//{
		//	name:        "valid host",
		//	url:         "mail.google.com/mail/u/0/#inbox",
		//	expectedErr: nil,
		//},
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
		html []byte
		linksExp map[int]string
	}{
		//{
		//	name: "empty node",
		//	html: nil,
		//	linksExp: nil,
		//},
		//{
		//	name: "invalid node",
		//	html: []byte("invalid"),
		//	linksExp: map[int]string{},
		//},
		//{
		//	name: "invalid node",
		//	html: []byte(`<!DOCTYPE html>
		//	<html>
		//	<body>
		//
		//	<h1>My First Heading</h1>
		//	<p>My first paragraph.</p>
		//
		//	</body>
		//	</html>`),
		//	linksExp: map[int]string{},
		//},
		{
			name: "invalid node",
			html: []byte(`
			<html>
			<body>

			<h1>My First Heading</h1>
			<p>My first paragraph.</p>
			<a href="#">Hello</a>
			</body>
			</html>`),
			linksExp: map[int]string{},
		},
	}
	for _, sc := range scenarios{
		t.Run(sc.name, func(t *testing.T){
			l := links(sc.html)
			if !cmp.Equal(l, sc.linksExp){
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
