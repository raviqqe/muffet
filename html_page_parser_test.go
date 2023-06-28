package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHtmlPageParserParsePage(t *testing.T) {
	_, err := newHtmlPageParser(newLinkFinder(nil, nil)).Parse("http://foo.com", nil)
	assert.Nil(t, err)
}

func TestHtmlPageParserFailWithInvalidURL(t *testing.T) {
	_, err := newHtmlPageParser(newLinkFinder(nil, nil)).Parse(":", nil)
	assert.NotNil(t, err)
}

func TestHtmlPageParserSetCorrectURL(t *testing.T) {
	s := "http://foo.com"

	p, err := newHtmlPageParser(newLinkFinder(nil, nil)).Parse(s, nil)
	assert.Nil(t, err)
	assert.Equal(t, s, p.URL().String())
}

func TestHtmlPageParserIgnorePageURLFragment(t *testing.T) {
	s := "http://foo.com#id"

	p, err := newHtmlPageParser(newLinkFinder(nil, nil)).Parse(s, nil)
	assert.Nil(t, err)
	assert.Equal(t, "http://foo.com", p.URL().String())
}

func TestHtmlPageParserKeepQuery(t *testing.T) {
	s := "http://foo.com?bar=baz"

	p, err := newHtmlPageParser(newLinkFinder(nil, nil)).Parse(s, nil)
	assert.Nil(t, err)
	assert.Equal(t, s, p.URL().String())
}

func TestHtmlPageParserResolveLinksWithBaseTag(t *testing.T) {
	p, err := newHtmlPageParser(newLinkFinder(nil, nil)).Parse(
		"http://foo.com",
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
	p, err := newHtmlPageParser(newLinkFinder(nil, nil)).Parse(
		"http://foo.com",
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
	_, err := newHtmlPageParser(newLinkFinder(nil, nil)).Parse(
		"http://foo.com",
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
	p, err := newHtmlPageParser(newLinkFinder(nil, nil)).Parse(
		"http://foo.com",
		[]byte(`<p id="foo" />`),
	)
	assert.Nil(t, err)
	assert.Equal(t, map[string]struct{}{"foo": {}}, p.Fragments())
}

func TestHtmlPageParserParseName(t *testing.T) {
	p, err := newHtmlPageParser(newLinkFinder(nil, nil)).Parse(
		"http://foo.com",
		[]byte(`<p name="foo" />`),
	)
	assert.Nil(t, err)
	assert.Equal(t, map[string]struct{}{"foo": {}}, p.Fragments())
}

func TestHtmlPageParserParseIDAndName(t *testing.T) {
	p, err := newHtmlPageParser(newLinkFinder(nil, nil)).Parse(
		"http://foo.com",
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
		p, err := newHtmlPageParser(newLinkFinder(nil, nil)).Parse(
			"http://foo.com",
			[]byte(ss[0]),
		)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(p.Links()))
		_, ok := p.Links()[ss[1]]
		assert.True(t, ok)
	}
}
