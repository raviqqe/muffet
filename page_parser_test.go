package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPageParserParsePage(t *testing.T) {
	_, err := newPageParser(newLinkFinder(nil)).Parse("http://foo.com", nil)
	assert.Nil(t, err)
}

func TestPageParserFailWithInvalidURL(t *testing.T) {
	_, err := newPageParser(newLinkFinder(nil)).Parse(":", nil)
	assert.NotNil(t, err)
}

func TestPageParserSetCorrectURL(t *testing.T) {
	s := "http://foo.com"

	p, err := newPageParser(newLinkFinder(nil)).Parse(s, nil)
	assert.Nil(t, err)
	assert.Equal(t, s, p.URL().String())
}

func TestPageParserIgnorePageURLFragment(t *testing.T) {
	s := "http://foo.com#id"

	p, err := newPageParser(newLinkFinder(nil)).Parse(s, nil)
	assert.Nil(t, err)
	assert.Equal(t, "http://foo.com", p.URL().String())
}

func TestPageParserKeepQuery(t *testing.T) {
	s := "http://foo.com?bar=baz"

	p, err := newPageParser(newLinkFinder(nil)).Parse(s, nil)
	assert.Nil(t, err)
	assert.Equal(t, s, p.URL().String())
}

func TestPageParserResolveLinksWithBaseTag(t *testing.T) {
	p, err := newPageParser(newLinkFinder(nil)).Parse(
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

func TestPageParserResolveLinksWithBlankBaseTag(t *testing.T) {
	p, err := newPageParser(newLinkFinder(nil)).Parse(
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

func TestPageParserFailToParseWithInvalidBaseTag(t *testing.T) {
	_, err := newPageParser(newLinkFinder(nil)).Parse(
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

func TestPageParserParseID(t *testing.T) {
	p, err := newPageParser(newLinkFinder(nil)).Parse(
		"http://foo.com",
		[]byte(`<p id="foo" />`),
	)
	assert.Nil(t, err)
	assert.Equal(t, map[string]struct{}{"foo": {}}, p.Fragments())
}

func TestPageParserParseName(t *testing.T) {
	p, err := newPageParser(newLinkFinder(nil)).Parse(
		"http://foo.com",
		[]byte(`<p name="foo" />`),
	)
	assert.Nil(t, err)
	assert.Equal(t, map[string]struct{}{"foo": {}}, p.Fragments())
}

func TestPageParserParseIDAndName(t *testing.T) {
	p, err := newPageParser(newLinkFinder(nil)).Parse(
		"http://foo.com",
		[]byte(`<p id="foo" name="bar" />`),
	)
	assert.Nil(t, err)
	assert.Equal(t, map[string]struct{}{"foo": {}, "bar": {}}, p.Fragments())
}

func TestPageParserParseLinks(t *testing.T) {
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
		p, err := newPageParser(newLinkFinder(nil)).Parse(
			"http://foo.com",
			[]byte(ss[0]),
		)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(p.Links()))
		_, ok := p.Links()[ss[1]]
		assert.True(t, ok)
	}
}
