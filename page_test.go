// +build !v2

package main

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestNewPage(t *testing.T) {
	_, err := newPage("https://foo.com", dummyHTML(t), newScraper(nil, false))
	assert.Nil(t, err)
}

func TestNewPageError(t *testing.T) {
	_, err := newPage(":", dummyHTML(t), newScraper(nil, false))
	assert.NotNil(t, err)
}

func TestPageURL(t *testing.T) {
	s := "https://foo.com"
	u, err := url.Parse(s)
	assert.Nil(t, err)

	p, err := newPage(s, dummyHTML(t), newScraper(nil, false))
	assert.Nil(t, err)

	assert.Equal(t, u, p.URL())
}

func TestPageURLWithBaseTag(t *testing.T) {
	n, err := html.Parse(strings.NewReader(`<base href="_blank" />`))
	assert.Nil(t, err)

	p, err := newPage("https://foo.com", n, newScraper(nil, false))
	assert.Nil(t, err)

	assert.Equal(t, "https://foo.com", p.URL().String())
}

func TestPageIDs(t *testing.T) {
	n, err := html.Parse(strings.NewReader(`<p id="foo">Hello!</p>`))
	assert.Nil(t, err)

	p, err := newPage("https://foo.com", n, newScraper(nil, false))
	assert.Nil(t, err)

	assert.Equal(t, 1, len(p.IDs()))
}

func TestPageLinks(t *testing.T) {
	for _, ss := range [][2]string{
		{
			`<a href="foo">bar</a>`,
			"https://foo.com/foo",
		},
		{
			`<img src="foo.img" />`,
			"https://foo.com/foo.img",
		},
		{
			`
				<html>
					<head>
						<base href="foo/" />
					</head>
					<body>
						<a href="foo">bar</a>
					</body>
				</html>
			`,
			"https://foo.com/foo/foo",
		},
	} {
		n, err := html.Parse(strings.NewReader(ss[0]))
		assert.Nil(t, err)

		p, err := newPage("https://foo.com", n, newScraper(nil, false))
		assert.Nil(t, err)

		assert.Equal(t, 1, len(p.Links()))

		_, ok := p.Links()[ss[1]]

		assert.True(t, ok)
	}
}

func TestPageWithParams(t *testing.T) {
	n, err := html.Parse(strings.NewReader(`<p id="foo">Hello!</p>`))
	assert.Nil(t, err)

	url := "https://foo.com/list?page=2"

	p, err := newPage(url, n, newScraper(nil, true))
	assert.Nil(t, err)

	assert.Equal(t, p.url.String(), url)
}

func TestPageWithoutParams(t *testing.T) {
	n, err := html.Parse(strings.NewReader(`<p id="foo">Hello!</p>`))
	assert.Nil(t, err)

	url := "https://foo.com/list?page=2"

	p, err := newPage(url, n, newScraper(nil, false))
	assert.Nil(t, err)

	assert.Equal(t, p.url.String(), url[:strings.IndexByte(url, '?')])
}
