package main

import (
	"errors"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newTestLinkFetcher(c *fakeHttpClient) *linkFetcher {
	return newTestLinkFetcherWithOptions(c, linkFetcherOptions{})
}

func newTestLinkFetcherWithOptions(c *fakeHttpClient, o linkFetcherOptions) *linkFetcher {
	return newLinkFetcher(c, newPageParser(newLinkFinder(nil, nil)), o)
}

func TestNewFetcher(t *testing.T) {
	newTestLinkFetcher(newFakeHttpClient(nil))
}

func TestLinkFetcherFetch(t *testing.T) {
	f := newTestLinkFetcher(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				if u.String() != "http://foo.com" {
					return nil, errors.New("")
				}

				return newFakeHtmlResponse("http://foo.com", ""), nil
			}),
	)

	s, p, err := f.Fetch("http://foo.com")

	assert.Equal(t, 200, s)
	assert.NotNil(t, p)
	assert.Nil(t, err)
}

func TestLinkFetcherFetchFromCache(t *testing.T) {
	ok := true
	s := "http://foo.com"

	f := newTestLinkFetcher(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				if !ok {
					return nil, errors.New("")
				}

				ok = false

				return newFakeHtmlResponse(s, ""), nil
			}),
	)

	sc, p, err := f.Fetch(s)
	assert.Equal(t, 200, sc)
	assert.NotNil(t, p)
	assert.Nil(t, err)

	sc, p, err = f.Fetch(s)
	assert.Equal(t, 200, sc)
	assert.NotNil(t, p)
	assert.Nil(t, err)
}

func TestLinkFetcherFetchCacheConcurrency(t *testing.T) {
	c := 0

	f := newTestLinkFetcher(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				c++

				return newFakeHtmlResponse("http://foo.com", ""), nil
			}),
	)

	g := &sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		g.Add(1)
		go func() {
			defer g.Done()

			time.Sleep(time.Millisecond)

			_, _, err := f.Fetch("http://foo.com")
			assert.Nil(t, err)
		}()
	}

	g.Wait()

	assert.Equal(t, 1, c)
}

func TestLinkFetcherFetchWithFragments(t *testing.T) {
	s := "http://foo.com"
	f := newTestLinkFetcher(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				if u.String() != s {
					return nil, errors.New("")
				}

				return newFakeHtmlResponse(s, `<p id="foo" />`), nil
			},
		),
	)

	sc, p, err := f.Fetch(s + "#foo")

	assert.Equal(t, 200, sc)
	assert.NotNil(t, p)
	assert.Nil(t, err)

	_, _, err = f.Fetch(s + "#bar")

	assert.Equal(t, "id #bar not found", err.Error())
}

func TestLinkFetcherFetchIgnoringFragments(t *testing.T) {
	s := "http://foo.com"
	f := newTestLinkFetcherWithOptions(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				if u.String() != s {
					return nil, errors.New("")
				}

				return newFakeHtmlResponse(s, ""), nil
			},
		),
		linkFetcherOptions{IgnoreFragments: true},
	)

	_, _, err := f.Fetch(s + "#bar")
	assert.Nil(t, err)
}

func TestLinkFetcherFetchSkippingTextFragment(t *testing.T) {
	s := "http://foo.com"
	f := newTestLinkFetcher(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				if u.String() != s {
					return nil, errors.New("")
				}

				return newFakeHtmlResponse(s, ""), nil
			},
		),
	)

	_, _, err := f.Fetch(s + "#:~:text=foo")
	assert.Nil(t, err)
}

func TestLinkFetcherFailToFetch(t *testing.T) {
	f := newTestLinkFetcher(
		newFakeHttpClient(func(*url.URL) (*fakeHttpResponse, error) {
			return nil, errors.New("")
		}))

	_, _, err := f.Fetch("http://foo.com")

	assert.NotNil(t, err)
}

func TestLinkFetcherFailToParseURL(t *testing.T) {
	f := newTestLinkFetcher(
		newFakeHttpClient(func(*url.URL) (*fakeHttpResponse, error) {
			return newFakeHtmlResponse("", ""), nil
		}))

	_, _, err := f.Fetch(":")

	assert.NotNil(t, err)
}

func TestSeparateFragment(t *testing.T) {
	for _, ss := range [][3]string{
		{"http://foo.com#bar", "http://foo.com", "bar"},
		{"#bar", "", "bar"},
	} {
		u, id, err := separateFragment(ss[0])

		assert.Nil(t, err)
		assert.Equal(t, ss[1], u)
		assert.Equal(t, ss[2], id)
	}
}

func TestFailToSeparateFragment(t *testing.T) {
	_, _, err := separateFragment(":")

	assert.NotNil(t, err)
}
