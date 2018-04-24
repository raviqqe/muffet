package main

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestNewPage(t *testing.T) {
	newPage("https://foo.com", nil)
}

func TestNewPageError(t *testing.T) {
	assert.Panics(t, func() {
		newPage(":", nil)
	})
}

func TestPageURL(t *testing.T) {
	s := "https://foo.com"
	u, err := url.Parse(s)

	assert.Nil(t, err)
	assert.Equal(t, u, newPage(s, nil).URL())
}

func TestPageBody(t *testing.T) {
	n, err := html.Parse(strings.NewReader("I'm Body."))

	assert.Nil(t, err)
	assert.Equal(t, n, newPage("https://foo.com", n).Body())
}
