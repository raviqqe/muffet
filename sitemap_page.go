package main

import (
	"net/url"
)

type sitemapPage struct {
	url   *url.URL
	links map[string]error
}

func newSitemapPage(u *url.URL, links map[string]error) *sitemapPage {
	return &sitemapPage{u, links}
}

func (p *sitemapPage) URL() *url.URL {
	return p.url
}

func (p *sitemapPage) Fragments() map[string]struct{} {
	return nil
}

func (p *sitemapPage) Links() map[string]error {
	return p.links
}
