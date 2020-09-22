package main

import (
	"net/url"
)

type page struct {
	url       *url.URL
	fragments map[string]struct{}
	links     map[string]error
}

func newPage(u *url.URL, fragments map[string]struct{}, links map[string]error) *page {
	return &page{u, fragments, links}
}

func (p *page) URL() *url.URL {
	return p.url
}

func (p *page) Fragments() map[string]struct{} {
	return p.fragments
}

func (p *page) Links() map[string]error {
	return p.links
}
