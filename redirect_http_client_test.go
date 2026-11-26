package main

import (
	"errors"
	"net/http"
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

type headerRecordingHttpClient struct {
	handler func(*url.URL) *fakeHttpResponse
	headers map[string]http.Header
}

func (c *headerRecordingHttpClient) Get(u *url.URL, h http.Header) (httpResponse, error) {
	c.headers[u.String()] = h.Clone()

	return c.handler(u), nil
}

func TestRedirectHttpClientSendCookies(t *testing.T) {
	u, err := url.Parse("http://foo.com/a")

	assert.Nil(t, err)

	c := &headerRecordingHttpClient{
		func(u *url.URL) *fakeHttpResponse {
			switch u.String() {
			case "http://foo.com/a":
				return newFakeHttpResponse(
					300,
					"",
					nil,
					map[string]string{"location": "/b", "set-cookie": "foo=bar"},
				)
			case "http://foo.com/b":
				return newFakeHttpResponse(
					300,
					"",
					nil,
					map[string]string{"location": "/c"},
				)
			}

			return newFakeHtmlResponse("", "")
		},
		map[string]http.Header{},
	}

	r, err := newRedirectHttpClient(c, 42).Get(u, nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode())
	assert.Equal(t, []string(nil), c.headers["http://foo.com/a"].Values("cookie"))
	assert.Equal(t, []string{"foo=bar"}, c.headers["http://foo.com/b"].Values("cookie"))
	assert.Equal(t, []string{"foo=bar"}, c.headers["http://foo.com/c"].Values("cookie"))
}

func TestRedirectHttpClientDoNotSendCookiesToOtherHosts(t *testing.T) {
	u, err := url.Parse("http://foo.com")

	assert.Nil(t, err)

	c := &headerRecordingHttpClient{
		func(u *url.URL) *fakeHttpResponse {
			if u.Host == "foo.com" {
				return newFakeHttpResponse(
					300,
					"",
					nil,
					map[string]string{"location": "http://bar.com", "set-cookie": "foo=bar"},
				)
			}

			return newFakeHtmlResponse("", "")
		},
		map[string]http.Header{},
	}

	r, err := newRedirectHttpClient(c, 42).Get(u, nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode())
	assert.Equal(t, []string(nil), c.headers["http://bar.com"].Values("cookie"))
}
