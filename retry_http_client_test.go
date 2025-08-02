package main

import (
	"errors"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeNetError struct{}

func (fakeNetError) Error() string   { return "network error" }
func (fakeNetError) Timeout() bool   { return true }
func (fakeNetError) Temporary() bool { return true }

func TestRetryHttpClientRetry(t *testing.T) {
	const maxRetries = 3

	u, err := url.Parse("http://foo.com/foo")
	assert.Nil(t, err)

	for _, tt := range []struct {
		errorCount int
		success    bool
	}{
		{errorCount: 0, success: true},
		{errorCount: 1, success: true},
		{errorCount: 2, success: true},
		{errorCount: 3, success: true},
		{errorCount: 4, success: false},
	} {
		t.Run(
			fmt.Sprintf("%d errors", tt.errorCount),
			func(t *testing.T) {
				count := 0

				c := newRetryHttpClient(
					newFakeHttpClient(
						func(u *url.URL) (*fakeHttpResponse, error) {
							count++

							if u.String() != "http://foo.com/foo" {
								return newFakeHtmlResponse("http://foo.com/", ""), nil
							} else if count <= tt.errorCount {
								return nil, &fakeNetError{}
							}

							return newFakeHtmlResponse("http://foo.com/foo", ""), nil
						},
					),
					maxRetries,
					0,
				)

				r, err := c.Get(u, nil)

				assert.Equal(t, tt.success, r != nil)
				assert.Equal(t, tt.success, err == nil)
				assert.Equal(t, min(tt.errorCount+1, maxRetries+1), count)
			},
		)
	}
}

func TestRetryHttpClientNoRetry(t *testing.T) {
	u, err := url.Parse("http://foo.com/foo")
	assert.Nil(t, err)

	count := 0

	c := newRetryHttpClient(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				count++
				return nil, errors.New("foo")
			},
		),
		42,
		0,
	)

	r, err := c.Get(u, nil)

	assert.Nil(t, r)
	assert.Equal(t, err, errors.New("foo"))
	assert.Equal(t, 1, count)
}
