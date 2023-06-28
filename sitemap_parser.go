package main

import (
	"bytes"
	"fmt"
	"net/url"

	sitemap "github.com/oxffaa/gopher-parse-sitemap"
)

type sitemapPageParser struct{}

func newSitemapPageParser() *sitemapPageParser {
	return &sitemapPageParser{}
}

func (f *sitemapPageParser) Parse(uu *url.URL, bs []byte) (*sitemapXmlPage, error) {
	us := map[string]error{}
	c := func(e interface{ GetLocation() string }) error {
		s := e.GetLocation()
		u, err := url.Parse(s)
		if err != nil {
			us[s] = err
			return nil
		}

		us[uu.ResolveReference(u).String()] = nil

		return nil
	}

	err := sitemap.Parse(bytes.NewReader(bs), func(e sitemap.Entry) error {
		return c(e)
	})

	if err == nil {
		return newSitemapXmlPage(uu, us), nil
	}

	err = sitemap.ParseIndex(bytes.NewReader(bs), func(e sitemap.IndexEntry) error {
		return c(e)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse sitemap: %v", err)
	}

	return newSitemapXmlPage(uu, us), nil
}
