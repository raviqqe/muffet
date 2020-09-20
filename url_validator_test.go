package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/temoto/robotstxt"
)

func TestURLValidatorReturnTrueForSameHostname(t *testing.T) {
	i := newURLValidator("foo.com", nil, nil)

	for _, s := range []string{
		"http://foo.com",
		"http://foo.com/bar",
		"https://foo.com",
	} {
		u, err := url.Parse(s)
		assert.Nil(t, err)
		assert.True(t, i.Validate(u))
	}
}

func TestURLValidatorReturnFalseForDifferentHostname(t *testing.T) {
	i := newURLValidator("foo.com", nil, nil)

	u, err := url.Parse("http://bar.com")
	assert.Nil(t, err)
	assert.False(t, i.Validate(u))
}

func TestURLValidatorValidateWithSitemap(t *testing.T) {
	i := newURLValidator(
		"foo.com",
		nil,
		map[string]struct{}{"http://foo.com/foo": {}},
	)

	u, err := url.Parse("http://foo.com/foo")
	assert.Nil(t, err)
	assert.True(t, i.Validate(u))

	u, err = url.Parse("http://foo.com/bar")
	assert.Nil(t, err)
	assert.False(t, i.Validate(u))
}

func TestURLValidatorValidateWithRobotsTxt(t *testing.T) {
	r, err := robotstxt.FromString(`
		User-Agent: *
		Disallow: /bar
	`)
	assert.Nil(t, err)

	i := newURLValidator("foo.com", r, nil)

	u, err := url.Parse("http://foo.com/foo")
	assert.Nil(t, err)
	assert.True(t, i.Validate(u))

	u, err = url.Parse("http://foo.com/bar")
	assert.Nil(t, err)
	assert.False(t, i.Validate(u))
}
