package main

import (
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testUrl = "http://foo.com"

var acceptedStatusCodes = statusCodeCollection{[]statusCodeRange{{200, 300}}}

func TestNewRedirectHttpClient(t *testing.T) {
	newRedirectHttpClient(newFakeHttpClient(nil), 42, acceptedStatusCodes)
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

				return newFakeHtmlResponse(testUrl, ""), nil
			},
		),
		42,
		acceptedStatusCodes,
	).Get(u, nil)

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
						testUrl,
						nil,
						map[string]string{"location": testUrl},
					), nil
				}

				return newFakeHtmlResponse("", ""), nil
			},
		),
		42,
		acceptedStatusCodes,
	).Get(u, nil)

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
						testUrl,
						nil,
						map[string]string{"location": testUrl},
					), nil
				}

				return newFakeHtmlResponse("", ""), nil
			},
		),
		maxRedirections,
		acceptedStatusCodes,
	).Get(u, nil)

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
							testUrl,
							nil,
							map[string]string{"location": "/foo"},
						), nil
					}
				}

				return nil, errors.New("")
			},
		),
		maxRedirections,
		acceptedStatusCodes,
	).Get(u, nil)

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
					testUrl,
					nil,
					map[string]string{"location": testUrl},
				), nil
			},
		),
		maxRedirections,
		acceptedStatusCodes,
	).Get(u, nil)

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
				return newFakeHttpResponse(300, testUrl, nil, nil), nil
			},
		),
		42,
		acceptedStatusCodes,
	).Get(u, nil)

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
					testUrl,
					nil,
					map[string]string{"location": ":"},
				), nil
			},
		),
		42,
		acceptedStatusCodes,
	).Get(u, nil)

	assert.Nil(t, r)
	assert.Contains(t, err.Error(), "parse")
}

func TestRedirectHttpClientFailWithInvalidStatusCode(t *testing.T) {
	u, err := url.Parse(testUrl)

	assert.Nil(t, err)

	r, err := newRedirectHttpClient(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				return newFakeHttpResponse(404, testUrl, nil, nil), nil
			},
		),
		42,
		acceptedStatusCodes,
	).Get(u, nil)

	assert.Nil(t, r)
	assert.Equal(t, err.Error(), "404")
}

func TestRedirectHttpClientFailAfterRedirect(t *testing.T) {
	u, err := url.Parse(testUrl)

	assert.Nil(t, err)

	redirected := false
	r, err := newRedirectHttpClient(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				if !redirected {
					redirected = true

					return newFakeHttpResponse(
						300,
						testUrl,
						nil,
						map[string]string{"location": "/foo"},
					), nil
				}

				return newFakeHttpResponse(404, "", nil, nil), nil
			},
		),
		42,
		acceptedStatusCodes,
	).Get(u, nil)

	assert.Nil(t, r)
	assert.Contains(t, err.Error(), "following redirect http://foo.com/foo")
}
