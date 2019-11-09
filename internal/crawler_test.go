package internal

import (
	"errors"
	"testing"
)

func TestNewCrawler(t *testing.T) {
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
			url:         "monzo.com",
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
			_, err := NewCrawler(sc.url)
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
