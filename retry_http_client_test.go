package main

import (
	"errors"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeNetError struct{}

func (fakeNetError) Error() string   { return "my network error" }
func (fakeNetError) Timeout() bool   { return true }
func (fakeNetError) Temporary() bool { return true }

func TestRetryHttpClientRetry(t *testing.T) {
	const maxRetries = 3

	u, err := url.Parse("http://foo.com/")
	assert.Nil(t, err)

	for errorCount, success := range map[int]bool{
		0: true,
		1: true,
		2: true,
		3: true,
		4: false,
	} {
		t.Run(
			fmt.Sprintf("%d errors", errorCount),
			func(t *testing.T) {
				count := 0

				c := newRetryHttpClient(
					newFakeHttpClient(
						func(*url.URL) (*fakeHttpResponse, error) {
							count++

							if count <= errorCount {
								return nil, &fakeNetError{}
							}

							return newFakeHtmlResponse("http://foo.com/", ""), nil
						},
					),
					maxRetries,
					0,
				)

				r, err := c.Get(u, nil)

				assert.Equal(t, success, r != nil)
				assert.Equal(t, success, err == nil)
				assert.Equal(t, min(errorCount+1, maxRetries+1), count)
			},
		)
	}
}

func TestRetryHttpClientNoRetry(t *testing.T) {
	u, err := url.Parse("http://foo.com/")
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

func TestRetryHttpClientRetryExceededError(t *testing.T) {
	u, err := url.Parse("http://foo.com/")
	assert.Nil(t, err)

	count := 0

	c := newRetryHttpClient(
		newFakeHttpClient(
			func(*url.URL) (*fakeHttpResponse, error) {
				count++
				return nil, fakeNetError{}
			},
		),
		42,
		0,
	)

	r, err := c.Get(u, nil)

	assert.Nil(t, r)
	assert.ErrorContains(t, err, "max retry count 42 exceeded: my network error")
	assert.Equal(t, 43, count)
}
