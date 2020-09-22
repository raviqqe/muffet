package main

import (
	"net/url"

	"github.com/yterajima/go-sitemap"
)

type sitemapFetcher struct {
	client httpClient
}

func newSitemapFetcher(c httpClient) *sitemapFetcher {
	sitemap.SetFetch(func(s string, _ interface{}) ([]byte, error) {
		u, err := url.Parse(s)
		if err != nil {
			return nil, err
		}

		r, err := c.Get(u, nil)

		if err != nil {
			return nil, err
		}

		return r.Body(), err
	})

	return &sitemapFetcher{c}
}

func (f *sitemapFetcher) Fetch(uu *url.URL) (map[string]struct{}, error) {
	u := *uu
	u.Path = "sitemap.xml"

	us := map[string]struct{}{}

	sm, err := sitemap.Get(u.String(), nil)
	if err != nil {
		return nil, err
	}

	for _, u := range sm.URL {
		us[u.Loc] = struct{}{}
	}

	return us, nil
}
