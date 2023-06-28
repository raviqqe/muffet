package main

import (
	"bytes"
	"net/url"

	sitemap "github.com/oxffaa/gopher-parse-sitemap"
)

type sitemapPageParser struct{}

func newSitemapPageParser() *sitemapPageParser {
	return &sitemapPageParser{}
}

func (f *sitemapPageParser) Parse(rawURL string, typ string, bs []byte) (page, error) {
	if typ != "application/xml" {
		return nil, nil
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	ls := map[string]error{}
	c := func(e interface{ GetLocation() string }) error {
		ls[e.GetLocation()] = nil
		return nil
	}

	err = sitemap.Parse(bytes.NewReader(bs), func(e sitemap.Entry) error {
		return c(e)
	})

	// TODO Detect XML files as sitemaps.
	if err == nil {
		return newSitemapXmlPage(u, ls), nil
	}

	err = sitemap.ParseIndex(bytes.NewReader(bs), func(e sitemap.IndexEntry) error {
		return c(e)
	})

	// TODO Detect XML files as sitemaps.
	if err != nil {
		return nil, nil
	}

	return newSitemapXmlPage(u, ls), nil
}
