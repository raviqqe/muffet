package main

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestScrapePage(t *testing.T) {
	b, err := url.Parse("https://localhost")
	assert.Nil(t, err)

	for _, c := range []struct {
		html  string
		links int
	}{
		{``, 0},
		{`<a href="/" />`, 1},
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

		for _, err := range newScraper(nil, false).Scrape(n, b) {
			if err == nil {
				s++
			} else {
				e++
			}
		}

		assert.Equal(t, c.links, s)
		assert.Equal(t, 0, e)
	}
}

func TestScrapePageError(t *testing.T) {
	b, err := url.Parse("https://localhost")
	assert.Nil(t, err)

	n, err := html.Parse(strings.NewReader(htmlWithBody(`<a href=":" />`)))
	assert.Nil(t, err)

	s, e := 0, 0

	for _, err := range newScraper(nil, false).Scrape(n, b) {
		if err == nil {
			s++
		} else {
			e++
		}
	}

	assert.Equal(t, 0, s)
	assert.Equal(t, 1, e)
}

func TestScraperIsURLExcluded(t *testing.T) {
	for _, x := range []struct {
		url     string
		regexps []string
		answer  bool
	}{
		{
			rootURL,
			[]string{"localhost"},
			true,
		},
		{
			rootURL,
			[]string{"localhost", "foo"},
			true,
		},
		{
			rootURL,
			[]string{"foo", "localhost"},
			true,
		},
		{
			rootURL,
			[]string{"foo"},
			false,
		},
	} {
		rs, err := compileRegexps(x.regexps)
		assert.Nil(t, err)

		assert.Equal(t, x.answer, newScraper(rs, false).isURLExcluded(x.url))
	}
}
