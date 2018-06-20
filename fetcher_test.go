package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFetcher(t *testing.T) {
	newFetcher(fetcherOptions{})
}

func TestFetcherFetch(t *testing.T) {
	f := newFetcher(fetcherOptions{})

	for _, s := range []string{rootURL, existentURL, escapedWSExistentURL, fragmentURL, erroneousURL} {
		r, err := f.Fetch(s)
		_, ok := r.Page()

		assert.Equal(t, 200, r.StatusCode())
		assert.True(t, ok)
		assert.Nil(t, err)
	}
}

func TestFetcherFetchCache(t *testing.T) {
	f := newFetcher(fetcherOptions{})

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

func TestFetcherFetchWithEscapedWS(t *testing.T) {
	f := newFetcher(fetcherOptions{})

	r, err := f.Fetch(escapedWSExistentURL)
	_, ok := r.Page()

	assert.Equal(t, 200, r.StatusCode())
	assert.True(t, ok)
	assert.Nil(t, err)

	_, err = f.Fetch(escapedWSNonexistentURL)
	assert.NotNil(t, err)
}

func TestFetcherFetchWithFragments(t *testing.T) {
	f := newFetcher(fetcherOptions{})

	r, err := f.Fetch(existentIDURL)
	_, ok := r.Page()

	assert.Equal(t, 200, r.StatusCode())
	assert.True(t, ok)
	assert.Nil(t, err)

	_, err = f.Fetch(nonExistentIDURL)
	assert.Equal(t, "id #bar not found", err.Error())
}

func TestFetcherFetchIgnoreFragments(t *testing.T) {
	_, err := newFetcher(fetcherOptions{}).Fetch(nonExistentIDURL)

	assert.NotNil(t, err)

	r, err := newFetcher(fetcherOptions{IgnoreFragments: true}).Fetch(nonExistentIDURL)

	assert.NotEqual(t, fetchResult{}, r)
	assert.Nil(t, err)
}

func TestFetcherFetchWithTLSVerification(t *testing.T) {
	_, err := newFetcher(fetcherOptions{}).Fetch(selfCertificateURL)
	assert.NotNil(t, err)

	p, err := newFetcher(fetcherOptions{SkipTLSVerification: true}).Fetch(selfCertificateURL)
	assert.NotEqual(t, page{}, p)
	assert.Nil(t, err)
}

func TestFetcherFetchWithInfiniteRedirections(t *testing.T) {
	_, err := newFetcher(fetcherOptions{}).Fetch(infiniteRedirectURL)
	assert.NotNil(t, err)
}

func TestFetcherFetchError(t *testing.T) {
	f := newFetcher(fetcherOptions{})

	for _, s := range []string{nonExistentURL, ":"} {
		_, err := f.Fetch(s)

		assert.NotNil(t, err)
	}
}

func TestFetcherSendRequest(t *testing.T) {
	f := newFetcher(fetcherOptions{})

	for _, s := range []string{rootURL, existentURL, fragmentURL, erroneousURL, redirectURL} {
		r, err := f.sendRequest(s)
		_, ok := r.Page()

		assert.Equal(t, 200, r.StatusCode())
		assert.True(t, ok)
		assert.Nil(t, err)
	}
}

func TestFetcherSendRequestWithMissingLocationHeader(t *testing.T) {
	_, err := newFetcher(fetcherOptions{}).sendRequest(invalidRedirectURL)
	assert.NotNil(t, err)
}

func TestFetcherSendRequestWithInvalidMIMEType(t *testing.T) {
	_, err := newFetcher(fetcherOptions{}).sendRequest(invalidMIMETypeURL)
	assert.Equal(t, "mime: no media type", err.Error())
}

func TestFetcherSendRequestWithTimeout(t *testing.T) {
	_, err := newFetcher(fetcherOptions{Timeout: 1 * time.Second}).sendRequest(timeoutURL)
	assert.NotNil(t, err)

	r, err := newFetcher(fetcherOptions{Timeout: 60 * time.Second}).sendRequest(timeoutURL)
	assert.Equal(t, 200, r.StatusCode())
	assert.Nil(t, err)
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

func TestSeparateFragmentError(t *testing.T) {
	_, _, err := separateFragment(":")

	assert.NotNil(t, err)
}
