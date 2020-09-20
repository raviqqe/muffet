package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLValidatorValidateWithSitemapXML(t *testing.T) {
	i := newURLValidator("", nil, nil)

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
	i := newURLValidator("", nil, nil)

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
