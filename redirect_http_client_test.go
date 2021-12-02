package main

import (
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testUrl = "http://foo.com"

func TestNewRedirectHttpClient(t *testing.T) {
	newRedirectHttpClient(newFakeHttpClient(nil), 42)
}

func TestRedirectHttpClientGet(t *testing.T) {
	u, err := url.Parse(testUrl)

	assert.Nil(t, err)

	r, err := newRedirectHttpClient(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				if u.String() != testUrl {
					return nil, errors.New("")
				}

				return newFakeHttpResponse(
					200,
					"http://foo.com",
					nil,
					nil,
				), nil
			},
		),
		42,
	).Get(u)

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode())
}
