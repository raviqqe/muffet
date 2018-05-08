package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetcherOptionsInitialize(t *testing.T) {
	o := fetcherOptions{}
	o.Initialize()

	assert.Equal(t, defaultConcurrency, o.Concurrency)
	assert.Equal(t, defaultMaxRedirections, o.MaxRedirections)
	assert.Equal(t, defaultTimeout, o.Timeout)
}
