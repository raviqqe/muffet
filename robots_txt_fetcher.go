package main

import (
	"fmt"
	"net/url"

	"github.com/temoto/robotstxt"
)

type robotsTxtFetcher struct {
	client httpClient
}

func newRobotsTxtFetcher(c httpClient) *robotsTxtFetcher {
	return &robotsTxtFetcher{c}
}

func (f *robotsTxtFetcher) Fetch(uu *url.URL) (*robotstxt.RobotsData, error) {
	u := *uu
	u.Path = "robots.txt"
	r, err := f.client.Get(&u, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch robots.txt: %v", err)
	}

	return robotstxt.FromBytes(r.Body())
}
