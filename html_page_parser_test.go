package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const HTML_MIME_TYPE = "text/html"

func parseURL(t *testing.T, s string) *url.URL {
	u, err := url.Parse(s)

	assert.Nil(t, err)

	return u
}

func TestHtmlPageParserParsePage(t *testing.T) {
	_, err := newHtmlPageParser(newTestLinkFinder()).Parse(parseURL(t, "http://foo.com"), HTML_MIME_TYPE, nil)
	assert.Nil(t, err)
}

func TestHtmlPageParserSetCorrectURL(t *testing.T) {
	u := parseURL(t, "http://foo.com")

	p, err := newHtmlPageParser(newTestLinkFinder()).Parse(u, HTML_MIME_TYPE, nil)
	assert.Nil(t, err)
	assert.Equal(t, u.String(), p.URL().String())
}

func TestHtmlPageParserIgnorePageURLFragment(t *testing.T) {
	p, err := newHtmlPageParser(newTestLinkFinder()).Parse(parseURL(t, "http://foo.com#id"), HTML_MIME_TYPE, nil)
	assert.Nil(t, err)
	assert.Equal(t, "http://foo.com", p.URL().String())
}

func TestHtmlPageParserKeepQuery(t *testing.T) {
	s := "http://foo.com?bar=baz"

	p, err := newHtmlPageParser(newTestLinkFinder()).Parse(parseURL(t, s), HTML_MIME_TYPE, nil)
	assert.Nil(t, err)
	assert.Equal(t, s, p.URL().String())
}

func TestHtmlPageParserResolveLinksWithBaseTag(t *testing.T) {
	p, err := newHtmlPageParser(newTestLinkFinder()).Parse(
		parseURL(t, "http://foo.com"),
		HTML_MIME_TYPE,
		[]byte(`
			<html>
			  <head>
					<base href="foo/" />
				</head>
				<body>
				  <a href="bar" />
				</body>
			</html>
		`),
	)
	assert.Nil(t, err)
	assert.Equal(t, map[string]error{"http://foo.com/foo/bar": nil}, p.Links())
}

func TestHtmlPageParserResolveLinksWithBlankBaseTag(t *testing.T) {
	p, err := newHtmlPageParser(newTestLinkFinder()).Parse(
		parseURL(t, "http://foo.com"),
		HTML_MIME_TYPE,
		[]byte(`
			<html>
			  <head>
					<base href="_blank" />
				</head>
				<body>
				  <a href="bar" />
				</body>
			</html>
		`),
	)
	assert.Nil(t, err)
	assert.Equal(t, map[string]error{"http://foo.com/bar": nil}, p.Links())
}

func TestHtmlPageParserFailToParseWithInvalidBaseTag(t *testing.T) {
	_, err := newHtmlPageParser(newTestLinkFinder()).Parse(
		parseURL(t, "http://foo.com"),
		HTML_MIME_TYPE,
		[]byte(`
			<html>
			  <head>
					<base href=":" />
				</head>
				<body>
				</body>
			</html>
		`),
	)
	assert.NotNil(t, err)
}

func TestHtmlPageParserParseID(t *testing.T) {
	p, err := newHtmlPageParser(newTestLinkFinder()).Parse(
		parseURL(t, "http://foo.com"),
		HTML_MIME_TYPE,
		[]byte(`<p id="foo" />`),
	)
	assert.Nil(t, err)
	assert.Equal(t, map[string]struct{}{"foo": {}}, p.Fragments())
}

func TestHtmlPageParserParseName(t *testing.T) {
	p, err := newHtmlPageParser(newTestLinkFinder()).Parse(
		parseURL(t, "http://foo.com"),
		HTML_MIME_TYPE,
		[]byte(`<p name="foo" />`),
	)
	assert.Nil(t, err)
	assert.Equal(t, map[string]struct{}{"foo": {}}, p.Fragments())
}

func TestHtmlPageParserParseIDAndName(t *testing.T) {
	p, err := newHtmlPageParser(newTestLinkFinder()).Parse(
		parseURL(t, "http://foo.com"),
		HTML_MIME_TYPE,
		[]byte(`<p id="foo" name="bar" />`),
	)
	assert.Nil(t, err)
	assert.Equal(t, map[string]struct{}{"foo": {}, "bar": {}}, p.Fragments())
}

func TestHtmlPageParserParseLinks(t *testing.T) {
	for _, ss := range [][2]string{
		{
			`<a href="foo">bar</a>`,
			"http://foo.com/foo",
		},
		{
			`<img src="foo.img" />`,
			"http://foo.com/foo.img",
		},
	} {
		p, err := newHtmlPageParser(newTestLinkFinder()).Parse(
			parseURL(t, "http://foo.com"),
			HTML_MIME_TYPE,
			[]byte(ss[0]),
		)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(p.Links()))
		_, ok := p.Links()[ss[1]]
		assert.True(t, ok)
	}
}
