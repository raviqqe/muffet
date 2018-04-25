package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestNewFetcher(t *testing.T) {
	newFetcher(1, false)
}

func TestFetcherCache(t *testing.T) {
	f := newFetcher(1, false)

	p, err := f.Fetch(rootURL)

	assert.NotEqual(t, (*page)(nil), p)
	assert.Nil(t, err)

	p, err = f.Fetch(nonExistentURL)

	assert.Equal(t, (*page)(nil), p)
	assert.NotNil(t, err)

	p, err = f.Fetch(rootURL)

	assert.Equal(t, (*page)(nil), p)
	assert.Nil(t, err)

	p, err = f.Fetch(nonExistentURL)

	assert.Equal(t, (*page)(nil), p)
	assert.NotNil(t, err)
}

func TestFetcherFetchIgnoreFragments(t *testing.T) {
	p, err := newFetcher(1, false).Fetch(nonExistentIDURL)

	assert.Equal(t, (*page)(nil), p)
	assert.NotNil(t, err)

	p, err = newFetcher(1, true).Fetch(nonExistentIDURL)

	assert.NotEqual(t, (*page)(nil), p)
	assert.Nil(t, err)
}

func TestFetcherFetchError(t *testing.T) {
	f := newFetcher(1, false)

	for _, s := range []string{nonExistentURL, ":"} {
		p, err := f.Fetch(s)

		assert.Equal(t, (*page)(nil), p)
		assert.NotNil(t, err)
	}
}

func TestFetcherFetchHTML(t *testing.T) {
	f := newFetcher(1, false)

	for _, s := range []string{rootURL, existentURL, fragmentURL, erroneousURL} {
		n, err := f.fetchHTML(s, "")

		assert.NotEqual(t, (*html.Node)(nil), n)
		assert.Nil(t, err)
	}
}

func TestFetcherFetchHTMLWithFragment(t *testing.T) {
	f := newFetcher(1, false)

	n, err := f.fetchHTML(fragmentURL, "foo")
	assert.NotEqual(t, (*html.Node)(nil), n)
	assert.Nil(t, err)

	n, err = f.fetchHTML(fragmentURL, "bar")
	assert.Equal(t, (*html.Node)(nil), n)
	assert.NotNil(t, err)
}

func TestFetcherFetchHTMLError(t *testing.T) {
	f := newFetcher(1, false)

	for _, s := range []string{":", nonExistentURL} {
		n, err := f.fetchHTML(s, "")

		assert.Equal(t, (*html.Node)(nil), n)
		assert.NotNil(t, err)
	}
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
