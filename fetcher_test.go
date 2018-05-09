package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFetcher(t *testing.T) {
	newFetcher(fetcherOptions{})
}

func TestFetcherFetchPage(t *testing.T) {
	p, err := newFetcher(fetcherOptions{}).FetchPage(rootURL)

	assert.NotEqual(t, page{}, p)
	assert.Nil(t, err)
}

func TestFetcherFetchPageWithNonHTMLPage(t *testing.T) {
	_, err := newFetcher(fetcherOptions{}).FetchPage(robotsTxtURL)
	assert.Equal(t, "non-HTML page", err.Error())
}

func TestFetcherFetchPageWithTLSVerification(t *testing.T) {
	_, err := newFetcher(fetcherOptions{}).FetchPage(selfCertificateURL)
	assert.NotNil(t, err)

	p, err := newFetcher(fetcherOptions{SkipTLSVerification: true}).FetchPage(selfCertificateURL)
	assert.NotEqual(t, page{}, p)
	assert.Nil(t, err)
}

func TestFetcherFetchPageWithInfiniteRedirections(t *testing.T) {
	_, err := newFetcher(fetcherOptions{}).FetchPage(infiniteRedirectURL)
	assert.NotNil(t, err)
}

func TestFetcherFetchLinkCache(t *testing.T) {
	f := newFetcher(fetcherOptions{})

	r, err := f.FetchLink(rootURL)
	assert.NotEqual(t, fetchResult{}, r)
	assert.Nil(t, err)
	_, ok := r.Page()
	assert.True(t, ok)

	_, err = f.FetchLink(nonExistentURL)
	assert.NotNil(t, err)

	r, err = f.FetchLink(rootURL)
	assert.NotEqual(t, fetchResult{}, r)
	assert.Nil(t, err)
	_, ok = r.Page()
	assert.False(t, ok)

	_, err = f.FetchLink(nonExistentURL)
	assert.NotNil(t, err)
}

func TestFetcherFetchLinkIgnoreFragments(t *testing.T) {
	_, err := newFetcher(fetcherOptions{}).FetchLink(nonExistentIDURL)

	assert.NotNil(t, err)

	r, err := newFetcher(fetcherOptions{IgnoreFragments: true}).FetchLink(nonExistentIDURL)

	assert.NotEqual(t, fetchResult{}, r)
	assert.Nil(t, err)
}

func TestFetcherFetchLinkError(t *testing.T) {
	f := newFetcher(fetcherOptions{})

	for _, s := range []string{nonExistentURL, ":"} {
		_, err := f.FetchLink(s)

		assert.NotNil(t, err)
	}
}

func TestFetcherSendRequestWithFragment(t *testing.T) {
	f := newFetcher(fetcherOptions{})

	for _, s := range []string{rootURL, existentURL, fragmentURL, erroneousURL} {
		r, err := f.sendRequestWithFragment(s, "")
		_, ok := r.Page()

		assert.Equal(t, 200, r.StatusCode())
		assert.True(t, ok)
		assert.Nil(t, err)
	}
}

func TestFetcherSendRequestWithFragmentWithFragment(t *testing.T) {
	f := newFetcher(fetcherOptions{})

	r, err := f.sendRequestWithFragment(fragmentURL, "foo")
	_, ok := r.Page()

	assert.Equal(t, 200, r.StatusCode())
	assert.True(t, ok)
	assert.Nil(t, err)

	_, err = f.sendRequestWithFragment(fragmentURL, "bar")
	assert.NotNil(t, err)
}

func TestFetcherSendRequestWithFragmentError(t *testing.T) {
	f := newFetcher(fetcherOptions{})

	for _, s := range []string{":", nonExistentURL} {
		_, err := f.sendRequestWithFragment(s, "")

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
