package main

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func newTestLinkFinder() linkFinder {
	return newLinkFinder(newTestLinkFilterer())
}

func TestLinkFinderFindLinks(t *testing.T) {
	b, err := url.Parse("https://localhost")
	assert.Nil(t, err)

	for _, c := range []struct {
		html      string
		linkCount int
	}{
		{``, 0},
		{`<a href="" />`, 0},
		{`<a href="/" />`, 1},
		{`<a href="/foo" />`, 1},
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

		for _, err := range newTestLinkFinder().Find(n, b) {
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

func TestLinkFinderFindLinkWithoutEncodedSpaces(t *testing.T) {
	b, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	n, err := html.Parse(strings.NewReader(
		htmlWithBody(`<a href="http://foo.com/a%20b" />`)),
	)
	assert.Nil(t, err)

	ls := newTestLinkFinder().Find(n, b)

	err, ok := ls["http://foo.com/a%20b"]
	assert.True(t, ok)
	assert.Nil(t, err)
}

func TestLinkFinderFindLinkWithoutSpacesNotEncoded(t *testing.T) {
	b, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	n, err := html.Parse(strings.NewReader(
		htmlWithBody(`<a href="http://foo.com/a b" />`)),
	)
	assert.Nil(t, err)

	ls := newTestLinkFinder().Find(n, b)

	err, ok := ls["http://foo.com/a%20b"]
	assert.True(t, ok)
	assert.Nil(t, err)
}

func TestLinkFinderFindLinkWithLeadingAndTrailingSpaces(t *testing.T) {
	b, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	n, err := html.Parse(strings.NewReader(
		htmlWithBody(`<a href=" http://foo.com " />`)),
	)
	assert.Nil(t, err)

	ls := newTestLinkFinder().Find(n, b)

	err, ok := ls["http://foo.com"]
	assert.True(t, ok)
	assert.Nil(t, err)
}

func TestLinkFinderFailWithInvalidURL(t *testing.T) {
	b, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	n, err := html.Parse(strings.NewReader(htmlWithBody(`<a href=":" />`)))
	assert.Nil(t, err)

	ls := newTestLinkFinder().Find(n, b)

	assert.Equal(t, 1, len(ls))
	assert.NotNil(t, ls[":"])
}

func TestLinkFinderFindLinkInSrcSet(t *testing.T) {
	b, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	n, err := html.Parse(strings.NewReader(
		htmlWithBody(`<source srcset="foo.png" />`)),
	)
	assert.Nil(t, err)

	ls := newTestLinkFinder().Find(n, b)

	err, ok := ls["http://foo.com/foo.png"]
	assert.True(t, ok)
	assert.Nil(t, err)
}

func TestLinkFinderFindMultipleLinksInSrcSet(t *testing.T) {
	b, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	n, err := html.Parse(strings.NewReader(
		htmlWithBody(`<source srcset="foo.png, bar.png" />`)),
	)
	assert.Nil(t, err)

	ls := newTestLinkFinder().Find(n, b)

	err, ok := ls["http://foo.com/foo.png"]
	assert.True(t, ok)
	assert.Nil(t, err)

	err, ok = ls["http://foo.com/bar.png"]
	assert.True(t, ok)
	assert.Nil(t, err)
}

func TestLinkFinderFindMultipleLinksInSrcSetWithDescriptors(t *testing.T) {
	b, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	n, err := html.Parse(strings.NewReader(
		htmlWithBody(`<source srcset="foo.png 100w, bar.png 200w" />`)),
	)
	assert.Nil(t, err)

	ls := newTestLinkFinder().Find(n, b)

	err, ok := ls["http://foo.com/foo.png"]
	assert.True(t, ok)
	assert.Nil(t, err)

	err, ok = ls["http://foo.com/bar.png"]
	assert.True(t, ok)
	assert.Nil(t, err)
}

func TestLinkFinderFindMetaTags(t *testing.T) {
	b, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	n, err := html.Parse(strings.NewReader(
		htmlWithHead(`<meta property="og:image" content="foo.png" />`)),
	)
	assert.Nil(t, err)

	ls := newTestLinkFinder().Find(n, b)

	err, ok := ls["http://foo.com/foo.png"]
	assert.True(t, ok)
	assert.Nil(t, err)
}

func TestLinkFinderFindLinkWithSpaces(t *testing.T) {
	b, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	n, err := html.Parse(strings.NewReader(
		htmlWithBody(`
    <a href="https:
    //foo.com
    /
    foo.html" />
    `)),
	)
	assert.Nil(t, err)

	ls := newTestLinkFinder().Find(n, b)

	err, ok := ls["http://foo.com/foo.html"]
	assert.True(t, ok)
	assert.Nil(t, err)
}

func TestLinkFinderIgnoreMetaTags(t *testing.T) {
	b, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	n, err := html.Parse(strings.NewReader(
		htmlWithHead(`<meta property="og:title" content="title" />`)),
	)
	assert.Nil(t, err)

	ls := newTestLinkFinder().Find(n, b)

	assert.Len(t, ls, 0)
}

func htmlWithBody(b string) string {
	return fmt.Sprintf(`<html><body>%v</body></html>`, b)
}

func htmlWithHead(b string) string {
	return fmt.Sprintf(`<html><head>%v</head><body><p>hi</p></body></html>`, b)
}
