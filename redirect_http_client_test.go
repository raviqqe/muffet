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
						map[string]string{"location": "http://foo.com"},
					), nil
				}

				return newFakeHtmlResponse("", ""), nil
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
						map[string]string{"location": "http://foo.com"},
					), nil
				}

				return newFakeHtmlResponse("", ""), nil
			},
		),
		maxRedirections,
	).Get(u)

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode())
	assert.Equal(t, maxRedirections, i)
}

func TestRedirectHttpClientGetWithRelativeRedirect(t *testing.T) {
	const maxRedirections = 42

	u, err := url.Parse(testUrl)

	assert.Nil(t, err)

	redirected := false
	r, err := newRedirectHttpClient(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				switch u.String() {
				case "http://foo.com/foo":
					return newFakeHtmlResponse("", ""), nil
				case testUrl:
					if !redirected {
						redirected = true

						return newFakeHttpResponse(
							300,
							"http://foo.com",
							nil,
							map[string]string{"location": "/foo"},
						), nil
					}
				}

				return nil, errors.New("")
			},
		),
		maxRedirections,
	).Get(u)

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode())
	assert.Equal(t, redirected, true)
}

func TestRedirectHttpClientFailWithTooManyRedirects(t *testing.T) {
	const maxRedirections = 42

	u, err := url.Parse(testUrl)

	assert.Nil(t, err)

	i := 0
	r, err := newRedirectHttpClient(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				i++

				return newFakeHttpResponse(
					300,
					"http://foo.com",
					nil,
					map[string]string{"location": "http://foo.com"},
				), nil
			},
		),
		maxRedirections,
	).Get(u)

	assert.Nil(t, r)
	assert.Equal(t, err.Error(), "too many redirections")
	assert.Equal(t, maxRedirections+1, i)
}

func TestRedirectHttpClientFailWithUnsetLocationHeader(t *testing.T) {
	u, err := url.Parse(testUrl)

	assert.Nil(t, err)

	r, err := newRedirectHttpClient(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				return newFakeHttpResponse(300, "http://foo.com", nil, nil), nil
			},
		),
		42,
	).Get(u)

	assert.Nil(t, r)
	assert.Equal(t, err.Error(), "location header not set")
}

func TestRedirectHttpClientFailWithInvalidLocationURL(t *testing.T) {
	u, err := url.Parse(testUrl)

	assert.Nil(t, err)

	r, err := newRedirectHttpClient(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				return newFakeHttpResponse(
					300,
					"http://foo.com",
					nil,
					map[string]string{"location": ":"},
				), nil
			},
		),
		42,
	).Get(u)

	assert.Nil(t, r)
	assert.Contains(t, err.Error(), "parse")
}

func TestRedirectHttpClientFailWithInvalidStatusCode(t *testing.T) {
	u, err := url.Parse(testUrl)

	assert.Nil(t, err)

	r, err := newRedirectHttpClient(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				return newFakeHttpResponse(404, "http://foo.com", nil, nil), nil
			},
		),
		42,
	).Get(u)

	assert.Nil(t, r)
	assert.Equal(t, err.Error(), "404")
}
