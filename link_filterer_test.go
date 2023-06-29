package main

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestLinkFilterer() linkFilterer {
	return newLinkFilterer(nil, nil)
}

func TestLinkFiltererIsLinkExcluded(t *testing.T) {
	u, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	for _, x := range []struct {
		regexps []string
		answer  bool
	}{
		{
			[]string{"foo\\.com"},
			false,
		},
		{
			[]string{"foo"},
			false,
		},
		{
			[]string{"bar", "foo"},
			false,
		},
		{
			[]string{"bar"},
			true,
		},
	} {
		t.Run(fmt.Sprint(x.regexps), func(t *testing.T) {
			rs, err := compileRegexps(x.regexps)
			assert.Nil(t, err)

			assert.Equal(t, x.answer, newLinkFilterer(rs, nil).IsValid(u))
		})
	}
}

func TestLinkFiltererIsLinkIncluded(t *testing.T) {
	u, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	for _, x := range []struct {
		regexps []string
		answer  bool
	}{
		{
			[]string{"foo\\.com"},
			true,
		},
		{
			[]string{"foo"},
			true,
		},
		{
			[]string{"bar", "foo"},
			true,
		},
		{
			[]string{"bar"},
			false,
		},
	} {
		t.Run(fmt.Sprint(x.regexps), func(t *testing.T) {
			rs, err := compileRegexps(x.regexps)
			assert.Nil(t, err)

			assert.Equal(t, x.answer, newLinkFilterer(nil, rs).IsValid(u))
		})
	}
}

func TestLinkFiltererExcludeEntireUrl(t *testing.T) {
	b, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	rs, err := compileRegexps([]string{"foo"})
	assert.Nil(t, err)

	assert.False(t, newLinkFilterer(rs, nil).IsValid(b))
}

func TestLinkFiltererIncludeEntireUrl(t *testing.T) {
	b, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	rs, err := compileRegexps([]string{"foo"})
	assert.Nil(t, err)

	assert.True(t, newLinkFilterer(nil, rs).IsValid(b))
}

func TestLinkFiltererExcludeInvalidScheme(t *testing.T) {
	b, err := url.Parse("mailto:foo@bar.baz")
	assert.Nil(t, err)

	assert.False(t, newLinkFilterer(nil, nil).IsValid(b))
}
