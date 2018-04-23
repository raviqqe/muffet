package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPage(t *testing.T) {
	newPage("https://foo.com", []byte(""))
}

func TestNewPageError(t *testing.T) {
	assert.Panics(t, func() {
		newPage(":", []byte(""))
	})
}

func TestPageURL(t *testing.T) {
	s := "https://foo.com"
	u, err := url.Parse(s)

	assert.Nil(t, err)
	assert.Equal(t, u, newPage(s, []byte("")).URL())
}

func TestPageBody(t *testing.T) {
	b := []byte("I'm Body.")
	assert.Equal(t, b, newPage("https://foo.com", b).Body())
}
