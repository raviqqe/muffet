package main

import (
	"bytes"
	"encoding/gob"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestNewPage(t *testing.T) {
	_, err := newPage("https://foo.com", dummyHTML(t), newScraper(nil))
	assert.Nil(t, err)
}

func TestNewPageError(t *testing.T) {
	_, err := newPage(":", dummyHTML(t), newScraper(nil))
	assert.NotNil(t, err)
}

func TestPageURL(t *testing.T) {
	s := "https://foo.com"
	u, err := url.Parse(s)
	assert.Nil(t, err)

	p, err := newPage(s, dummyHTML(t), newScraper(nil))
	assert.Nil(t, err)

	assert.Equal(t, u, p.URL())
}

func TestPageURLWithBaseTag(t *testing.T) {
	n, err := html.Parse(strings.NewReader(`<base href="_blank" />`))
	assert.Nil(t, err)

	p, err := newPage("https://foo.com", n, newScraper(nil))
	assert.Nil(t, err)

	assert.Equal(t, "https://foo.com", p.URL().String())
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

func TestPageEncodeDecode(t *testing.T) {
	for _, s := range []string{
		"",
		`
			<a href=":" />
			<a href="mailto:me@right.here" />
			<a href="/bar" />
			<a href="#foo" />
		`,
	} {
		n, err := html.Parse(strings.NewReader(htmlWithBody(s)))
		assert.Nil(t, err)

		p, err := newPage("https://foo.com", n, newScraper(nil))
		assert.Nil(t, err)

		b := bytes.NewBuffer(nil)
		assert.Nil(t, gob.NewEncoder(b).Encode(p))

		q := &page{}
		assert.Nil(t, gob.NewDecoder(b).Decode(q))

		assertPagesEqual(t, p, q)
	}
}

func assertPagesEqual(t *testing.T, p, q *page) {
	if p == nil {
		assert.Equal(t, p, q)
		return
	}

	assert.Equal(t, p.URL(), q.URL())
	assert.Equal(t, p.IDs(), q.IDs())

	for k, v := range p.Links() {
		if v == nil {
			assert.Equal(t, nil, q.Links()[k])
			continue
		}

		assert.Equal(t, v.Error(), q.Links()[k].Error())
	}
}

func TestPageUnmarshalError(t *testing.T) {
	p := page{}
	assert.NotNil(t, p.UnmarshalBinary(nil))
}
