package main

import (
	"fmt"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeNetError struct{}

func (fakeNetError) Error() string   { return "network error" }
func (fakeNetError) Timeout() bool   { return true }
func (fakeNetError) Temporary() bool { return true }

func TestRetryHttpClientGet(t *testing.T) {
	u, err := url.Parse("http://foo.com/foo")
	assert.Nil(t, err)

	for _, tt := range []struct {
		errorCount       int
		expectedRequests int
		success          bool
	}{
		{errorCount: 0, expectedRequests: 1, success: true},
		{errorCount: 1, expectedRequests: 2, success: true},
		{errorCount: 2, expectedRequests: 3, success: true},
		{errorCount: 3, expectedRequests: 3, success: false},
	} {
		t.Run(
			fmt.Sprintf("%d retries", tt.errorCount),
			func(t *testing.T) {
				var count atomic.Int32
				c := newRetryHttpClient(
					newFakeHttpClient(
						func(u *url.URL) (*fakeHttpResponse, error) {
							if u.String() == "http://foo.com/foo" {
								if count.Add(1) <= int32(tt.errorCount) {
									return nil, &fakeNetError{}
								}
								return newFakeHtmlResponse("http://foo.com/foo", ""), nil
							}
							return newFakeHtmlResponse("http://foo.com/", ""), nil
						},
					),
					3,
				)

				r, err := c.Get(u, nil)
				assert.Equal(t, tt.success, r != nil)
				assert.Equal(t, tt.success, err == nil)
				assert.Equal(t, int32(tt.expectedRequests), count.Load())
			},
		)
	}
}
