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

func TestRedirectHttpClientGetWithRedirect(t *testing.T) {
	u, err := url.Parse(testUrl)

	assert.Nil(t, err)

	redirected := false
	r, err := newRedirectHttpClient(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				if u.String() != testUrl {
					return nil, errors.New("")
				} else if !redirected {
					redirected = true

					return newFakeHttpResponse(
						300,
						"http://foo.com",
						nil,
						map[string]string{"Location": "http://foo.com"},
					), nil
				}

				return newFakeHtmlResponse("http://foo.com", ""), nil
			},
		),
		42,
	).Get(u)

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode())
}

func TestRedirectHttpClientGetWithRedirects(t *testing.T) {
	const maxRedirections = 42

	u, err := url.Parse(testUrl)

	assert.Nil(t, err)

	i := 0
	r, err := newRedirectHttpClient(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				if u.String() != testUrl {
					return nil, errors.New("")
				} else if i < maxRedirections {
					i++

					return newFakeHttpResponse(
						300,
						"http://foo.com",
						nil,
						map[string]string{"Location": "http://foo.com"},
					), nil
				}

				return newFakeHtmlResponse("http://foo.com", ""), nil
			},
		),
		maxRedirections,
	).Get(u)

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode())
	assert.Equal(t, maxRedirections, i)
}
