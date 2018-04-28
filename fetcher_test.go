package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFetcher(t *testing.T) {
	newFetcher(1, false)
}

func TestFetcherCache(t *testing.T) {
	f := newFetcher(1, false)

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
	assert.False(t, ok)

	_, err = f.Fetch(nonExistentURL)
	assert.NotNil(t, err)
}

func TestFetcherFetchIgnoreFragments(t *testing.T) {
	_, err := newFetcher(1, false).Fetch(nonExistentIDURL)

	assert.NotNil(t, err)

	r, err := newFetcher(1, true).Fetch(nonExistentIDURL)

	assert.NotEqual(t, fetchResult{}, r)
	assert.Nil(t, err)
}

func TestFetcherFetchError(t *testing.T) {
	f := newFetcher(1, false)

	for _, s := range []string{nonExistentURL, ":"} {
		_, err := f.Fetch(s)

		assert.NotNil(t, err)
	}
}

func TestFetcherSendRequestWithFragment(t *testing.T) {
	f := newFetcher(1, false)

	for _, s := range []string{rootURL, existentURL, fragmentURL, erroneousURL} {
		r, err := f.sendRequestWithFragment(s, "")

		assert.NotEqual(t, fetchResult{}, r)
		assert.Nil(t, err)
	}
}

func TestFetcherSendRequestWithFragmentWithFragment(t *testing.T) {
	f := newFetcher(1, false)

	r, err := f.sendRequestWithFragment(fragmentURL, "foo")
	assert.NotEqual(t, fetchResult{}, r)
	assert.Nil(t, err)

	_, err = f.sendRequestWithFragment(fragmentURL, "bar")
	assert.NotNil(t, err)
}

func TestFetcherSendRequestWithFragmentError(t *testing.T) {
	f := newFetcher(1, false)

	for _, s := range []string{":", nonExistentURL} {
		_, err := f.sendRequestWithFragment(s, "")

		assert.NotNil(t, err)
	}
}

func TestFetcherSendRequest(t *testing.T) {
	f := newFetcher(1, false)

	for _, s := range []string{rootURL, existentURL, fragmentURL, erroneousURL, redirectURL} {
		r, err := f.sendRequest(s)

		assert.NotEqual(t, fetchResult{}, r)
		assert.Nil(t, err)
	}
}

func TestFetcherSendRequestWithMissingLocationHeader(t *testing.T) {
	_, err := newFetcher(1, false).sendRequest(invalidRedirectURL)

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

func TestSeparateFragmentError(t *testing.T) {
	_, _, err := separateFragment(":")

	assert.NotNil(t, err)
}
