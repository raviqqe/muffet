package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLinkResult(t *testing.T) {
	newLinkResult(200)
}

func TestNewLinkResultWithPage(t *testing.T) {
	newLinkResultWithPage(200, newPage("", nil))
}

func TestLinkResultStatusCode(t *testing.T) {
	assert.Equal(t, 200, newLinkResult(200).StatusCode())
}

func TestLinkResultPage(t *testing.T) {
	p, ok := newLinkResult(200).Page()

	assert.False(t, ok)
	assert.Equal(t, page{}, p)

	q := newPage("", nil)
	p, ok = newLinkResultWithPage(200, q).Page()

	assert.True(t, ok)
	assert.Equal(t, q, p)
}
