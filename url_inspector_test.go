package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewURLInspector(t *testing.T) {
	_, err := newURLInspector(rootURL, false)
	assert.Nil(t, err)
}

func TestNewURLInspectorError(t *testing.T) {
	_, err := newURLInspector(":", false)
	assert.NotNil(t, err)
}

func TestNewURLInspectorWithSitemap(t *testing.T) {
	_, err := newURLInspector(rootURL, true)
	assert.Nil(t, err)
}

func TestNewURLInspectorWithMissingSitemap(t *testing.T) {
	_, err := newURLInspector(missingSitemapURL, true)
	assert.NotNil(t, err)
}

func TestURLInspectorInspect(t *testing.T) {
	i, err := newURLInspector(rootURL, true)
	assert.Nil(t, err)

	for _, s := range []string{rootURL, existentURL} {
		u, err := url.Parse(s)
		assert.Nil(t, err)
		assert.True(t, i.Inspect(u))
	}

	for _, s := range []string{nonExistentURL, erroneousURL} {
		u, err := url.Parse(s)
		assert.Nil(t, err)
		assert.False(t, i.Inspect(u))
	}
}
