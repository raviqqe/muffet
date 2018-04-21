package main

import "net/url"

type page struct {
	url  *url.URL
	body []byte
}

func newPage(s string, b []byte) page {
	u, err := url.Parse(s)

	if err != nil {
		panic(err)
	}

	u.Fragment = ""
	u.RawQuery = ""

	return page{u, b}
}

func (p page) URL() *url.URL {
	return p.url
}

func (p page) Body() []byte {
	return p.body
}
