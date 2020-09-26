package main

import (
	"errors"
	"net/url"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func TestSitemapFetcherFetchSitemap(t *testing.T) {
	s := "http://foo.com"
	u, err := url.Parse(s)
	assert.Nil(t, err)

	sm, err := newSitemapFetcher(
		newFakeHTTPClient(
			func(u *url.URL) (*fakeHTTPResponse, error) {
				if u.String() != s+"/sitemap.xml" {
					return nil, errors.New("")
				}

				return newFakeHTTPResponse(
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
				), nil
			})).Fetch(u)

	assert.Nil(t, err)

	_, ok := sm["http://foo.com/bar"]
	assert.True(t, ok)
}

func TestSitemapFetcherFailToFetchSitemap(t *testing.T) {
	u, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	_, err = newSitemapFetcher(
		newFakeHTTPClient(
			func(u *url.URL) (*fakeHTTPResponse, error) {
				return nil, errors.New("foo")
			},
		),
	).Fetch(u)

	cupaloy.SnapshotT(t, err.Error())
}

func TestSitemapFetcherFailToParseSitemap(t *testing.T) {
	u, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	_, err = newSitemapFetcher(
		newFakeHTTPClient(
			func(u *url.URL) (*fakeHTTPResponse, error) {
				return newFakeHTTPResponse(200, "", "text/xml", []byte(`<`)), nil
			},
		),
	).Fetch(u)

	cupaloy.SnapshotT(t, err.Error())
}
