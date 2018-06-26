package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewURLInspector(t *testing.T) {
	_, err := newURLInspector(rootURL, false, false, false)
	assert.Nil(t, err)
}

func TestNewURLInspectorError(t *testing.T) {
	_, err := newURLInspector(":", false, false, false)
	assert.NotNil(t, err)
}

func TestNewURLInspectorWithSitemapXML(t *testing.T) {
	_, err := newURLInspector(rootURL, false, true, false)
	assert.Nil(t, err)
}

func TestNewURLInspectorErrorWithRobotsTxt(t *testing.T) {
	for _, s := range []string{missingMetadataURL, invalidRobotsTxtURL, noResponseURL} {
		_, err := newURLInspector(s, true, false, false)
		assert.NotNil(t, err)
	}
}

func TestNewURLInspectorWithMissingSitemapXML(t *testing.T) {
	_, err := newURLInspector(missingMetadataURL, false, true, false)
	assert.NotNil(t, err)
}

func TestURLInspectorInspectWithSitemapXML(t *testing.T) {
	i, err := newURLInspector(rootURL, false, true, false)
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
	i, err := newURLInspector(rootURL, true, false, false)
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
