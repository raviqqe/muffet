package main

import (
	"errors"
	"net/url"
	"strconv"
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

				return newFakeHtmlResponse(testUrl, ""), nil
			},
		),
		42,
	).Get(u, nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode())
}

func TestRedirectHttpClientGetWithoutRedirect(t *testing.T) {
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
		0,
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
	).Get(u, nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode())
	assert.Equal(t, redirected, true)
}

func TestRedirectHttpClientFailWithTooManyRedirects(t *testing.T) {
	for _, n := range []int{0, 1, 42} {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
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
				n,
			).Get(u, nil)

			assert.Nil(t, r)
			assert.Equal(t, err.Error(), "too many redirections")
			assert.Equal(t, n+1, i)
		})
	}
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
	).Get(u, nil)

	assert.Nil(t, r)
	assert.Contains(t, err.Error(), "parse")
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

				return nil, errors.New("foo")
			},
		),
		42,
	).Get(u, nil)

	assert.Nil(t, r)
	assert.Contains(t, err.Error(), "following redirect http://foo.com/foo")
}
