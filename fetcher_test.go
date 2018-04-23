package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFetcher(t *testing.T) {
	newFetcher(1)
}

func TestFetcherCache(t *testing.T) {
	f := newFetcher(1)

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

func TestFetcherFetchError(t *testing.T) {
	f := newFetcher(1)

	for _, s := range []string{nonExistentURL, ":"} {
		p, err := f.Fetch(s)

		assert.Equal(t, (*page)(nil), p)
		assert.NotNil(t, err)
	}
}
