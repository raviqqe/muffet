package main

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestNewPage(t *testing.T) {
	newPage("https://foo.com", dummyHTML(t), newScraper(nil))
}

func TestNewPageError(t *testing.T) {
	assert.Panics(t, func() {
		newPage(":", dummyHTML(t), newScraper(nil))
	})
}

func TestPageURL(t *testing.T) {
	s := "https://foo.com"
	u, err := url.Parse(s)
	assert.Nil(t, err)

	p, err := newPage(s, dummyHTML(t), newScraper(nil))
	assert.Nil(t, err)

	assert.Equal(t, u, p.URL())
}

func TestPageIDs(t *testing.T) {
	n, err := html.Parse(strings.NewReader(`<p id="foo">Hello!</p>`))
	assert.Nil(t, err)

	p, err := newPage("https://foo.com", n, newScraper(nil))
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

		p, err := newPage("https://foo.com", n, newScraper(nil))
		assert.Nil(t, err)

		assert.Equal(t, 1, len(p.Links()))

		_, ok := p.Links()[ss[1]]

		assert.True(t, ok)
	}
}
