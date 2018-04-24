package main

import (
	"net/url"

	"golang.org/x/net/html"
)

type page struct {
	url  *url.URL
	body *html.Node
}

func newPage(s string, n *html.Node) page {
	u, err := url.Parse(s)

	if err != nil {
		panic(err)
	}

	u.Fragment = ""
	u.RawQuery = ""

	return page{u, n}
}

func (p page) URL() *url.URL {
	return p.url
}

func (p page) Body() *html.Node {
	return p.body
}
