package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewURLInspector(t *testing.T) {
	_, err := newURLInspector(newFakeHTTPClient(), rootURL, false, false)
	assert.Nil(t, err)
}

func TestNewURLInspectorError(t *testing.T) {
	_, err := newURLInspector(newFakeHTTPClient(), ":", false, false)
	assert.NotNil(t, err)
}

func TestNewURLInspectorWithSitemapXML(t *testing.T) {
	_, err := newURLInspector(newFakeHTTPClient(), rootURL, false, true)
	assert.Nil(t, err)
}

func TestNewURLInspectorErrorWithRobotsTxt(t *testing.T) {
	for _, s := range []string{missingMetadataURL, invalidRobotsTxtURL, noResponseURL} {
		_, err := newURLInspector(newFakeHTTPClient(), s, true, false)
		assert.NotNil(t, err)
	}
}

func TestNewURLInspectorWithMissingSitemapXML(t *testing.T) {
	for _, s := range []string{missingMetadataURL, noResponseURL} {
		_, err := newURLInspector(newFakeHTTPClient(), s, false, true)
		assert.NotNil(t, err)
	}
}

func TestURLInspectorInspectWithSitemapXML(t *testing.T) {
	i, err := newURLInspector(newFakeHTTPClient(), rootURL, false, true)
	assert.Nil(t, err)

	for _, s := range []string{rootURL, existentURL} {
		u, err := url.Parse(s)
		assert.Nil(t, err)
		assert.True(t, i.Inspect(u))
	}

	for _, s := range []string{nonExistentURL, erroneousURL, fragmentURL} {
		u, err := url.Parse(s)
		assert.Nil(t, err)
		assert.False(t, i.Inspect(u))
	}
}

func TestURLInspectorInspectWithRobotsTxt(t *testing.T) {
	i, err := newURLInspector(newFakeHTTPClient(), rootURL, true, false)
	assert.Nil(t, err)

	for _, s := range []string{rootURL, existentURL} {
		u, err := url.Parse(s)
		assert.Nil(t, err)
		assert.True(t, i.Inspect(u))
	}

	for _, s := range []string{erroneousURL, fragmentURL} {
		u, err := url.Parse(s)
		assert.Nil(t, err)
		assert.False(t, i.Inspect(u))
	}
}
