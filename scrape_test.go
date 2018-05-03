package main

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestScrapePage(t *testing.T) {
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

		bs, es := scrapePage(newPage("", n))

		assert.Equal(t, c.links, len(bs))
		assert.Equal(t, 0, len(es))
	}
}

func TestScrapePageError(t *testing.T) {
	n, err := html.Parse(strings.NewReader(htmlWithBody(`<a href=":" />`)))
	assert.Nil(t, err)

	bs, es := scrapePage(newPage("", n))

	assert.Equal(t, 0, len(bs))
	assert.Equal(t, 1, len(es))
}

func TestResolveURLWithAbsoluteURL(t *testing.T) {
	n, err := html.Parse(strings.NewReader(""))
	assert.Nil(t, err)

	u, err := url.Parse("http://localhost/foo/bar")
	assert.Nil(t, err)

	u, err = resolveURL(newPage("http://localhost", n), u)
	assert.Nil(t, err)
	assert.Equal(t, "http://localhost/foo/bar", u.String())
}

func TestResolveURLWithBaseTag(t *testing.T) {
	n, err := html.Parse(strings.NewReader(`<html><head><base href="/foo/" /></head></html>`))
	assert.Nil(t, err)

	u, err := url.Parse("bar")
	assert.Nil(t, err)

	u, err = resolveURL(newPage("http://localhost", n), u)
	assert.Nil(t, err)
	assert.Equal(t, "http://localhost/foo/bar", u.String())
}

func TestResolveURLError(t *testing.T) {
	n, err := html.Parse(strings.NewReader(`<html><head><base href=":" /></head></html>`))
	assert.Nil(t, err)

	u, err := url.Parse("bar")
	assert.Nil(t, err)

	_, err = resolveURL(newPage("http://localhost", n), u)
	assert.NotNil(t, err)
}
