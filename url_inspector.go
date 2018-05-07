package main

import (
	"net/url"

	"github.com/yterajima/go-sitemap"
)

type urlInspector struct {
	hostname     string
	includedURLs map[string]struct{}
}

func newURLInspector(s string, sm bool) (urlInspector, error) {
	u, err := url.Parse(s)

	if err != nil {
		return urlInspector{}, err
	}

	us := map[string]struct{}{}

	if sm {
		u.Path = "sitemap.xml"
		m, err := sitemap.Get(u.String(), nil)

		if err != nil {
			return urlInspector{}, err
		}

		for _, u := range m.URL {
			us[u.Loc] = struct{}{}
		}
	}

	return urlInspector{u.Hostname(), us}, nil
}

func (i urlInspector) Inspect(u *url.URL) bool {
	if len(i.includedURLs) != 0 {
		if _, ok := i.includedURLs[u.String()]; !ok {
			return false
		}
	}

	return u.Hostname() == i.hostname
}
