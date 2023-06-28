package main

import (
	"net/url"
)

type sitemapXmlPage struct {
	url   *url.URL
	links map[string]error
}

func newSitemapXmlPage(u *url.URL, links map[string]error) *sitemapXmlPage {
	return &sitemapXmlPage{u, links}
}

func (p *sitemapXmlPage) URL() *url.URL {
	return p.url
}

func (p *sitemapXmlPage) Fragments() map[string]struct{} {
	return nil
}

func (p *sitemapXmlPage) Links() map[string]error {
	return p.links
}
