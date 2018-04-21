package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFetcher(t *testing.T) {
	newFetcher(1)
}

func TestFetcherFetchError(t *testing.T) {
	_, err := newFetcher(1).Fetch("https://google.com/non/existent/path")

	assert.NotNil(t, err)
}
