package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestLinkFilterer() linkFilterer {
	return newLinkFilterer(nil, nil)
}

func TestLinkFiltererIsLinkExcluded(t *testing.T) {
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
		rs, err := compileRegexps(x.regexps)
		assert.Nil(t, err)

		assert.Equal(t, x.answer, newLinkFilterer(rs, nil).isLinkExcluded("http://foo.com"))
	}
}

func TestLinkFiltererIsLinkIncluded(t *testing.T) {
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
		rs, err := compileRegexps(x.regexps)
		assert.Nil(t, err)

		assert.Equal(t, x.answer, newLinkFilterer(nil, rs).isLinkIncluded("http://foo.com"))
	}
}

func TestLinkFiltererExcludeEntireUrl(t *testing.T) {
	b, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	rs, err := compileRegexps([]string{"foo"})
	assert.Nil(t, err)

	assert.False(t, newLinkFilterer(rs, nil).Filter(b))
}

func TestLinkFiltererIncludeEntireUrl(t *testing.T) {
	b, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	rs, err := compileRegexps([]string{"foo"})
	assert.Nil(t, err)

	assert.True(t, newLinkFilterer(nil, rs).Filter(b))
}
