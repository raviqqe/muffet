// +build !v2

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFetchResult(t *testing.T) {
	newFetchResult(200, nil)
}

func TestNewFetchResultWithPage(t *testing.T) {
	p, err := newPage("", dummyHTML(t), newScraper(nil, false))
	assert.Nil(t, err)

	newFetchResult(200, p)
}

func TestFetchResultStatusCode(t *testing.T) {
	assert.Equal(t, 200, newFetchResult(200, nil).StatusCode())
}

func TestFetchResultPage(t *testing.T) {
	p, ok := newFetchResult(200, nil).Page()

	assert.False(t, ok)
	assert.Equal(t, (*page)(nil), p)

	q, err := newPage("", dummyHTML(t), newScraper(nil, false))
	assert.Nil(t, err)

	p, ok = newFetchResult(200, q).Page()

	assert.True(t, ok)
	assert.Equal(t, q, p)
}
