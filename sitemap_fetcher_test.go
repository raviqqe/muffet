package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSitemapFetcherFetchSitemap(t *testing.T) {
	s := "http://foo.com"
	u, err := url.Parse(s)
	assert.Nil(t, err)

	sm, err := newSitemapFetcher(
		newFakeHTTPClient(
			map[string]*fakeHTTPResponse{
				s + "/sitemap.xml": newFakeHTTPResponse(
					200,
					s,
					"text/xml",
					[]byte(`
						<?xml version="1.0" encoding="UTF-8"?>
						<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
							<url>
								<loc>http://foo.com/bar</loc>
							</url>
						</urlset>
					`),
				),
			})).Fetch(u)

	assert.Nil(t, err)

	_, ok := sm["http://foo.com/bar"]
	assert.True(t, ok)
}
