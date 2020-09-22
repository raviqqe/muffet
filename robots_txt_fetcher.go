package main

import (
	"errors"
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
		return nil, err
	} else if r.StatusCode() != 200 {
		return nil, errors.New("failed to load robots.txt")
	}

	return robotstxt.FromBytes(r.Body())
}
