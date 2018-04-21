package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFetcher(t *testing.T) {
	newFetcher(1)
}

func TestFetcherFetchError(t *testing.T) {
	err := newFetcher(1).Fetch("https://google.com/non/existent/path", nil)

	assert.NotNil(t, err)
}
