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

func TestRedirectHttpClientSendCookiesWithAttributes(t *testing.T) {
	for _, x := range []struct {
		name, source, cookie, target string
		cookies                      []string
	}{
		{"path", "http://foo.com/foo/bar", "foo=bar; Path=/", "http://foo.com/baz", []string{"foo=bar"}},
		{"other path", "http://foo.com", "foo=bar; Path=/foo", "http://foo.com/bar", nil},
		{"domain", "http://foo.com", "foo=bar; Domain=foo.com", "http://bar.foo.com", []string{"foo=bar"}},
		{"public suffix domain", "http://foo.co.uk", "foo=bar; Domain=co.uk", "http://bar.co.uk", nil},
		{"secure", "https://foo.com", "foo=bar; Secure", "https://foo.com/foo", []string{"foo=bar"}},
		{"secure over HTTP", "https://foo.com", "foo=bar; Secure", "http://foo.com", nil},
		{"expired", "http://foo.com", "foo=bar; Max-Age=0", "http://foo.com/foo", nil},
		{"invalid", "http://foo.com", "foo", "http://foo.com/foo", nil},
	} {
		t.Run(x.name, func(t *testing.T) {
			u, err := url.Parse(x.source)

			assert.Nil(t, err)

			c := &headerRecordingHttpClient{
				func(u *url.URL) *fakeHttpResponse {
					if u.String() == x.source {
						return newFakeHttpResponse(
							300,
							"",
							nil,
							map[string]string{"location": x.target, "set-cookie": x.cookie},
						)
					}

					return newFakeHtmlResponse("", "")
				},
				map[string]http.Header{},
			}

			r, err := newRedirectHttpClient(c, 42).Get(u, nil)

			assert.Nil(t, err)
			assert.Equal(t, 200, r.StatusCode())
			assert.Equal(t, x.cookies, c.headers[x.target].Values("cookie"))
		})
	}
}
