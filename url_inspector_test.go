package main

import (
	"crypto/tls"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestNewURLInspector(t *testing.T) {
	_, err := newURLInspector(&fasthttp.Client{}, rootURL, false, false)
	assert.Nil(t, err)
}

func TestNewURLInspectorError(t *testing.T) {
	_, err := newURLInspector(&fasthttp.Client{}, ":", false, false)
	assert.NotNil(t, err)
}

func TestNewURLInspectorWithSitemapXML(t *testing.T) {
	_, err := newURLInspector(&fasthttp.Client{}, rootURL, false, true)
	assert.Nil(t, err)
}

func TestNewURLInspectorErrorWithRobotsTxt(t *testing.T) {
	for _, s := range []string{missingMetadataURL, invalidRobotsTxtURL, noResponseURL} {
		_, err := newURLInspector(&fasthttp.Client{}, s, true, false)
		assert.NotNil(t, err)
	}
}

func TestNewURLInspectorWithMissingSitemapXML(t *testing.T) {
	for _, s := range []string{missingMetadataURL, noResponseURL} {
		_, err := newURLInspector(&fasthttp.Client{}, s, false, true)
		assert.NotNil(t, err)
	}
}

func TestNewURLInspectorWithSelfCertifiedServer(t *testing.T) {
	for _, bs := range [][2]bool{{true, false}, {false, true}, {true, true}} {
		_, err := newURLInspector(&fasthttp.Client{}, selfCertificateURL, bs[0], bs[1])
		assert.NotNil(t, err)

		_, err = newURLInspector(
			&fasthttp.Client{TLSConfig: &tls.Config{InsecureSkipVerify: true}},
			selfCertificateURL, bs[0], bs[1])
		assert.Nil(t, err)
	}
}

func TestURLInspectorInspectWithSitemapXML(t *testing.T) {
	i, err := newURLInspector(&fasthttp.Client{}, rootURL, false, true)
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
	i, err := newURLInspector(&fasthttp.Client{}, rootURL, true, false)
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
