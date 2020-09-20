package main

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestLinkFinderFindLinks(t *testing.T) {
	b, err := url.Parse("https://localhost")
	assert.Nil(t, err)

	for _, c := range []struct {
		html      string
		linkCount int
	}{
		{``, 0},
		{`<a href="/" />`, 1},
		{`<a href="/f
		o 	o" />`, 1},
		// TODO: Test <frame> tag.
		{`<iframe src="/iframe"></iframe>`, 1},
		{`<img src="/foo.jpg" />`, 1},
		{`<link href="/link" />`, 1},
		{`<script src="/foo.js"></script>`, 1},
		{`<source src="/foo.png" />`, 1},
		{`<source srcset="/foo.png" />`, 1},
		{`<source src="/foo.png" srcset="/bar.png" />`, 2},
		{`<track src="/foo.vtt" />`, 1},
		{`<a href="/"><img src="/foo.png" /></a>`, 2},
		{`<a href="/" /><a href="/" />`, 1},
	} {
		n, err := html.Parse(strings.NewReader(htmlWithBody(c.html)))
		assert.Nil(t, err)

		s, e := 0, 0

		for _, err := range newLinkFinder(nil).Find(n, b) {
			if err == nil {
				s++
			} else {
				e++
			}
		}

		assert.Equal(t, c.linkCount, s)
		assert.Equal(t, 0, e)
	}
}

func TestLinkFinderScrapePageError(t *testing.T) {
	b, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	n, err := html.Parse(strings.NewReader(htmlWithBody(`<a href=":" />`)))
	assert.Nil(t, err)

	ls := newLinkFinder(nil).Find(n, b)

	assert.Equal(t, 1, len(ls))
	assert.NotNil(t, ls[":"])
}

func TestLinkFinderIsLinkExcluded(t *testing.T) {
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

		assert.Equal(t, x.answer, newLinkFinder(rs).isLinkExcluded("http://foo.com"))
	}
}
