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

func (f *sitemapPageParser) Parse(urlString string, bs []byte) (*sitemapXmlPage, error) {
	uu, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	us := map[string]error{}
	c := func(e interface{ GetLocation() string }) error {
		us[e.GetLocation()] = nil
		return nil
	}

	err = sitemap.Parse(bytes.NewReader(bs), func(e sitemap.Entry) error {
		return c(e)
	})

	// TODO Detect XML files as sitemaps.
	if err == nil {
		return newSitemapXmlPage(uu, us), nil
	}

	err = sitemap.ParseIndex(bytes.NewReader(bs), func(e sitemap.IndexEntry) error {
		return c(e)
	})

	// TODO Detect XML files as sitemaps.
	if err != nil {
		return nil, nil
	}

	return newSitemapXmlPage(uu, us), nil
}
