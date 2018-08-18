package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFetchResult(t *testing.T) {
	newFetchResult(200)
}

func TestNewFetchResultWithPage(t *testing.T) {
	p, err := newPage("", dummyHTML(t), newScraper(nil, false))
	assert.Nil(t, err)

	newFetchResultWithPage(200, p)
}

func TestFetchResultStatusCode(t *testing.T) {
	assert.Equal(t, 200, newFetchResult(200).StatusCode())
}

func TestFetchResultPage(t *testing.T) {
	p, ok := newFetchResult(200).Page()

	assert.False(t, ok)
	assert.Equal(t, (*page)(nil), p)

	q, err := newPage("", dummyHTML(t), newScraper(nil, false))
	assert.Nil(t, err)

	p, ok = newFetchResultWithPage(200, q).Page()

	assert.True(t, ok)
	assert.Equal(t, q, p)
}
