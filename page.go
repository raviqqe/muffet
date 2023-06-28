package main

import (
	"net/url"
)

type page interface {
	URL() *url.URL
	Fragments() map[string]struct{}
	Links() map[string]error
}

type htmlPage struct {
	url       *url.URL
	fragments map[string]struct{}
	links     map[string]error
}

func newHtmlPage(u *url.URL, fragments map[string]struct{}, links map[string]error) *htmlPage {
	return &htmlPage{u, fragments, links}
}

func (p *htmlPage) URL() *url.URL {
	return p.url
}

func (p *htmlPage) Fragments() map[string]struct{} {
	return p.fragments
}

func (p *htmlPage) Links() map[string]error {
	return p.links
}
