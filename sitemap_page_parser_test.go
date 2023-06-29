package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

const SITEMAP_MIME_TYPE = "application/xml"

func TestSitemapPageParserParsePage(t *testing.T) {
	p, err := newSitemapPageParser(newTestLinkFilterer()).Parse(parseURL(t, "https://foo.com/sitemap.xml"), SITEMAP_MIME_TYPE, []byte(`
		<?xml version="1.0" encoding="UTF-8"?>
		<urlset
			xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
			xmlns:news="http://www.google.com/schemas/sitemap-news/0.9"
			xmlns:xhtml="http://www.w3.org/1999/xhtml"
			xmlns:image="http://www.google.com/schemas/sitemap-image/1.1"
			xmlns:video="http://www.google.com/schemas/sitemap-video/1.1"
		>
			<url>
				<loc>https://foo.com/</loc>
			</url>
		</urlset>
	`))

	assert.Nil(t, err)
	assert.Equal(t, "https://foo.com/sitemap.xml", p.URL().String())
	assert.Equal(t, map[string]error{"https://foo.com/": nil}, p.Links())
	assert.Equal(t, map[string]struct{}(nil), p.Fragments())
}

func TestSitemapPageParserParseIndexPage(t *testing.T) {
	p, err := newSitemapPageParser(newTestLinkFilterer()).Parse(parseURL(t, "https://foo.com/sitemap.xml"), SITEMAP_MIME_TYPE, []byte(`
		<?xml version="1.0" encoding="UTF-8"?>
		<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
			<sitemap>
				<loc>https://foo.com/sitemap-0.xml</loc>
			</sitemap>
		</sitemapindex>
	`))

	assert.Nil(t, err)
	assert.Equal(t, "https://foo.com/sitemap.xml", p.URL().String())
	assert.Equal(t, map[string]error{"https://foo.com/sitemap-0.xml": nil}, p.Links())
	assert.Equal(t, map[string]struct{}(nil), p.Fragments())
}

func TestSitemapPageParserExcludeLink(t *testing.T) {
	p, err := newSitemapPageParser(
		newLinkFilterer([]*regexp.Regexp{regexp.MustCompile("private")}, nil),
	).Parse(parseURL(t, "https://foo.com/sitemap.xml"), SITEMAP_MIME_TYPE, []byte(`
		<?xml version="1.0" encoding="UTF-8"?>
		<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
			<sitemap>
				<loc>https://foo.com/private-sitemap.xml</loc>
			</sitemap>
			<sitemap>
				<loc>https://foo.com/public-sitemap.xml</loc>
			</sitemap>
		</sitemapindex>
	`))

	assert.Nil(t, err)
	assert.Equal(t, "https://foo.com/sitemap.xml", p.URL().String())
	assert.Equal(t, map[string]error{"https://foo.com/public-sitemap.xml": nil}, p.Links())
	assert.Equal(t, map[string]struct{}(nil), p.Fragments())
}

func TestSitemapPageParserFailToParseXML(t *testing.T) {
	p, err := newSitemapPageParser(newTestLinkFilterer()).Parse(parseURL(t, "https://foo.com/sitemap.xml"), SITEMAP_MIME_TYPE, []byte(""))

	assert.Nil(t, err)
	assert.Nil(t, p)
}

func TestSitemapPageParserFailToParseHTML(t *testing.T) {
	p, err := newSitemapPageParser(newTestLinkFilterer()).Parse(parseURL(t, "https://foo.com/sitemap.xml"), "text/html", []byte(""))

	assert.Nil(t, err)
	assert.Nil(t, p)
}
