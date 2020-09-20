package main

import (
	"errors"
	"net/url"
	"time"

	"github.com/yterajima/go-sitemap"
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

	us := map[string]struct{}{}

	sitemap.SetFetch(func(s string, _ interface{}) ([]byte, error) {
		u, err := url.Parse(s)

		if err != nil {
			return nil, err
		}

		r, err := f.client.Get(u, nil, time.Duration(0))

		if err != nil {
			return nil, err
		} else if r.StatusCode() != 200 {
			return nil, errors.New("failed to load sitemap")
		}

		return r.Body(), err
	})

	sm, err := sitemap.Get(u.String(), nil)
	if err != nil {
		return nil, err
	}

	for _, u := range sm.URL {
		us[u.Loc] = struct{}{}
	}

	return us, nil
}
