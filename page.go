package main

import (
	"net/url"
)

type page struct {
	url   *url.URL
	ids   map[string]struct{}
	links map[string]error
}

func newPage(u *url.URL, ids map[string]struct{}, links map[string]error) *page {
	return &page{u, ids, links}
}

func (p *page) URL() *url.URL {
	return p.url
}

func (p *page) IDs() map[string]struct{} {
	return p.ids
}

func (p *page) Links() map[string]error {
	return p.links
}
