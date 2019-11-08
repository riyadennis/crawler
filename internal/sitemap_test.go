package internal

import (
	"fmt"
	"testing"
)

func TestSiteMap(t *testing.T) {
	scenarios := []struct{
		name string
		url string
		expectedErr error
	}{
		{
			name: "invalid url",
			url: "invalid",
			expectedErr: fmt.Errorf(
				"Get http://invalid: dial tcp: lookup invalid" +
					" on 10.32.160.171:53: no such host"),
		},
		{
			name: "invalid host",
			url: "mmm",
			expectedErr:fmt.Errorf(
				"Get invalid: unsupported protocol scheme %q",
				""),
		},
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
