package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewURLValidator(t *testing.T) {
	_, err := newURLValidator(newFakeHTTPClient(), rootURL, false, false)
	assert.Nil(t, err)
}

func TestNewURLValidatorError(t *testing.T) {
	_, err := newURLValidator(newFakeHTTPClient(), ":", false, false)
	assert.NotNil(t, err)
}

func TestNewURLValidatorWithSitemapXML(t *testing.T) {
	_, err := newURLValidator(newFakeHTTPClient(), rootURL, false, true)
	assert.Nil(t, err)
}

func TestNewURLValidatorErrorWithRobotsTxt(t *testing.T) {
	for _, s := range []string{missingMetadataURL, invalidRobotsTxtURL, noResponseURL} {
		_, err := newURLValidator(newFakeHTTPClient(), s, true, false)
		assert.NotNil(t, err)
	}
}

func TestNewURLValidatorWithMissingSitemapXML(t *testing.T) {
	for _, s := range []string{missingMetadataURL, noResponseURL} {
		_, err := newURLValidator(newFakeHTTPClient(), s, false, true)
		assert.NotNil(t, err)
	}
}

func TestURLValidatorValidateWithSitemapXML(t *testing.T) {
	i, err := newURLValidator(newFakeHTTPClient(), rootURL, false, true)
	assert.Nil(t, err)

	for _, s := range []string{rootURL, existentURL} {
		u, err := url.Parse(s)
		assert.Nil(t, err)
		assert.True(t, i.Validate(u))
	}

	for _, s := range []string{nonExistentURL, erroneousURL, fragmentURL} {
		u, err := url.Parse(s)
		assert.Nil(t, err)
		assert.False(t, i.Validate(u))
	}
}

func TestURLValidatorValidateWithRobotsTxt(t *testing.T) {
	i, err := newURLValidator(newFakeHTTPClient(), rootURL, true, false)
	assert.Nil(t, err)

	for _, s := range []string{rootURL, existentURL} {
		u, err := url.Parse(s)
		assert.Nil(t, err)
		assert.True(t, i.Validate(u))
	}

	for _, s := range []string{erroneousURL, fragmentURL} {
		u, err := url.Parse(s)
		assert.Nil(t, err)
		assert.False(t, i.Validate(u))
	}
}
