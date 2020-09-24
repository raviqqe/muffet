package main

import (
	"errors"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createTestLinkFetcher(c *fakeHTTPClient) *linkFetcher {
	return createTestLinkFetcherWithOptions(c, linkFetcherOptions{})
}

func createTestLinkFetcherWithOptions(c *fakeHTTPClient, o linkFetcherOptions) *linkFetcher {
	return newLinkFetcher(c, newPageParser(newLinkFinder(nil)), o)
}

func TestNewFetcher(t *testing.T) {
	createTestLinkFetcher(newFakeHTTPClient(nil))
}

func TestLinkFetcherFetch(t *testing.T) {
	f := createTestLinkFetcher(
		newFakeHTTPClient(
			func(u *url.URL) (*fakeHTTPResponse, error) {
				if u.String() != "http://foo.com" {
					return nil, errors.New("")
				}

				return newFakeHTTPResponse(
					200,
					"http://foo.com",
					"text/html",
					nil,
				), nil
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

	f := createTestLinkFetcher(
		newFakeHTTPClient(
			func(u *url.URL) (*fakeHTTPResponse, error) {
				if !ok {
					return nil, errors.New("")
				}

				ok = false

				return newFakeHTTPResponse(
					200,
					s,
					"text/html",
					nil,
				), nil
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

	f := createTestLinkFetcher(
		newFakeHTTPClient(
			func(u *url.URL) (*fakeHTTPResponse, error) {
				c++

				return newFakeHTTPResponse(200, "http://foo.com", "text/html", nil), nil
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
	f := createTestLinkFetcher(
		newFakeHTTPClient(
			func(u *url.URL) (*fakeHTTPResponse, error) {
				if u.String() != s {
					return nil, errors.New("")
				}

				return newFakeHTTPResponse(200, s, "text/html", []byte(`<p id="foo" />`)), nil
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
	f := createTestLinkFetcherWithOptions(
		newFakeHTTPClient(
			func(u *url.URL) (*fakeHTTPResponse, error) {
				if u.String() != s {
					return nil, errors.New("")
				}

				return newFakeHTTPResponse(200, s, "text/html", nil), nil
			},
		),
		linkFetcherOptions{IgnoreFragments: true},
	)

	_, _, err := f.Fetch(s + "#bar")
	assert.Nil(t, err)
}

func TestLinkFetcherFailToFetch(t *testing.T) {
	f := createTestLinkFetcher(
		newFakeHTTPClient(func(*url.URL) (*fakeHTTPResponse, error) {
			return nil, errors.New("")
		}))

	_, _, err := f.Fetch("http://foo.com")

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
