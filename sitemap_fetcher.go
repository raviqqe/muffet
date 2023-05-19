package main

import (
	"bytes"
	"fmt"
	"net/url"

	sitemap "github.com/oxffaa/gopher-parse-sitemap"
)

type sitemapFetcher struct {
	client httpClient
}

func newSitemapFetcher(c httpClient) *sitemapFetcher {
	return &sitemapFetcher{c}
}

func (f *sitemapFetcher) Fetch(uu *url.URL) (map[string]struct{}, error) {
	u := *uu
	u.Path = "sitemap.xml"

	r, err := f.client.Get(&u, nil)
	if err != nil {
		return nil, f.formatGetError(err)
	}

	us := map[string]struct{}{}

	bs, err := r.Body()
	if err != nil {
		return nil, f.formatGetError(err)
	}

	err = sitemap.Parse(bytes.NewReader(bs), func(e sitemap.Entry) error {
		us[e.GetLocation()] = struct{}{}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse sitemap.xml: %v", err)
	}

	return us, nil
}

func (*sitemapFetcher) formatGetError(err error) error {
	return fmt.Errorf("failed to GET sitemap.xml: %v", err)
}
