package main

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFetcher(t *testing.T) {
	newFetcher(newFakeHTTPClient(nil), fetcherOptions{})
}

func TestFetcherFetch(t *testing.T) {
	f := newFetcher(newFakeHTTPClient(nil), fetcherOptions{})

	for _, s := range []string{rootURL, existentURL, fragmentURL, erroneousURL} {
		r, err := f.Fetch(s)
		_, ok := r.Page()

		assert.Equal(t, 200, r.StatusCode())
		assert.True(t, ok)
		assert.Nil(t, err)
	}
}

func TestFetcherFetchCache(t *testing.T) {
	f := newFetcher(newFakeHTTPClient(nil), fetcherOptions{})

	r, err := f.Fetch(rootURL)
	assert.NotEqual(t, fetchResult{}, r)
	assert.Nil(t, err)
	_, ok := r.Page()
	assert.True(t, ok)

	_, err = f.Fetch(nonExistentURL)
	assert.NotNil(t, err)

	r, err = f.Fetch(rootURL)
	assert.NotEqual(t, fetchResult{}, r)
	assert.Nil(t, err)
	_, ok = r.Page()
	assert.True(t, ok)

	_, err = f.Fetch(nonExistentURL)
	assert.NotNil(t, err)
}

func TestFetcherFetchCacheConcurrency(t *testing.T) {
	g := &sync.WaitGroup{}
	f := newFetcher(newFakeHTTPClient(nil), fetcherOptions{})

	for i := 0; i < 1000; i++ {
		g.Add(1)
		go func() {
			_, err := f.Fetch(countingURL)
			assert.Nil(t, err)
			g.Done()
		}()
	}

	g.Wait()

	assert.Equal(t, 1, testCountingHandler.Count())
}

func TestFetcherFetchWithFragments(t *testing.T) {
	f := newFetcher(newFakeHTTPClient(nil), fetcherOptions{})

	r, err := f.Fetch(existentIDURL)
	_, ok := r.Page()

	assert.Equal(t, 200, r.StatusCode())
	assert.True(t, ok)
	assert.Nil(t, err)

	_, err = f.Fetch(nonExistentIDURL)
	assert.Equal(t, "id #bar not found", err.Error())
}

func TestFetcherFetchIgnoreFragments(t *testing.T) {
	_, err := newFetcher(newFakeHTTPClient(nil), fetcherOptions{}).Fetch(nonExistentIDURL)

	assert.NotNil(t, err)

	r, err := newFetcher(newFakeHTTPClient(nil), fetcherOptions{IgnoreFragments: true}).Fetch(nonExistentIDURL)

	assert.NotEqual(t, fetchResult{}, r)
	assert.Nil(t, err)
}

func TestFetcherFetchWithInfiniteRedirections(t *testing.T) {
	_, err := newFetcher(newFakeHTTPClient(nil), fetcherOptions{}).Fetch(infiniteRedirectURL)
	assert.NotNil(t, err)
}

func TestFetcherFetchError(t *testing.T) {
	f := newFetcher(newFakeHTTPClient(nil), fetcherOptions{})

	for _, s := range []string{nonExistentURL, ":"} {
		_, err := f.Fetch(s)

		assert.NotNil(t, err)
	}
}

func TestFetcherSendRequest(t *testing.T) {
	f := newFetcher(newFakeHTTPClient(nil), fetcherOptions{})

	for _, s := range []string{rootURL, existentURL, fragmentURL, erroneousURL, redirectURL} {
		r, err := f.sendRequest(s)
		_, ok := r.Page()

		assert.Equal(t, 200, r.StatusCode())
		assert.True(t, ok)
		assert.Nil(t, err)
	}
}

func TestFetcherSendRequestWithMissingLocationHeader(t *testing.T) {
	_, err := newFetcher(newFakeHTTPClient(nil), fetcherOptions{}).sendRequest(invalidRedirectURL)
	assert.NotNil(t, err)
}

func TestFetcherSendRequestWithInvalidMIMEType(t *testing.T) {
	_, err := newFetcher(newFakeHTTPClient(nil), fetcherOptions{}).sendRequest(invalidMIMETypeURL)
	assert.Equal(t, "mime: no media type", err.Error())
}

func TestFetcherSendRequestWithTimeout(t *testing.T) {
	_, err := newFetcher(newFakeHTTPClient(nil), fetcherOptions{Timeout: 1 * time.Second}).sendRequest(timeoutURL)
	assert.NotNil(t, err)

	r, err := newFetcher(newFakeHTTPClient(nil), fetcherOptions{Timeout: 60 * time.Second}).sendRequest(timeoutURL)
	assert.Equal(t, 200, r.StatusCode())
	assert.Nil(t, err)
}

// TODO Eanble it fixing its flakiness.
// func TestFetcherSendRequestConcurrency(t *testing.T) {
// 	c := 900
// 	f := newFetcher(&fasthttp.Client{MaxConnsPerHost: c}, fetcherOptions{Concurrency: c})

// 	g := sync.WaitGroup{}

// 	for i := 0; i < 10000; i++ {
// 		g.Add(1)
// 		go func() {
// 			_, err := f.sendRequest("http://httpbin.org/get")
// 			assert.Nil(t, err)
// 			g.Done()
// 		}()
// 	}

// 	g.Wait()
// }

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

func TestSeparateFragmentError(t *testing.T) {
	_, _, err := separateFragment(":")

	assert.NotNil(t, err)
}
