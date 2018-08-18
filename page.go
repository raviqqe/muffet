package main

import (
	"net/url"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type page struct {
	url   *url.URL
	ids   map[string]struct{}
	links map[string]error
}

func newPage(s string, n *html.Node, sc scraper) (*page, error) {
	u, err := urlParse(s, sc.removeNewlines)

	if err != nil {
		return nil, err
	}

	u.Fragment = ""
	u.RawQuery = ""

	ids := map[string]struct{}{}

	scrape.FindAllNested(n, func(n *html.Node) bool {
		if s := scrape.Attr(n, "id"); s != "" {
			ids[s] = struct{}{}
		}

		return false
	})

	b := u

	if n, ok := scrape.Find(n, func(n *html.Node) bool {
		return n.DataAtom == atom.Base
	}); ok {
		u, err := urlParse(scrape.Attr(n, "href"), sc.removeNewlines)

		if err != nil {
			return nil, err
		}

		b = b.ResolveReference(u)
	}

	return &page{u, ids, sc.Scrape(n, b)}, nil
}

func (p page) URL() *url.URL {
	return p.url
}

func (p page) IDs() map[string]struct{} {
	return p.ids
}

func (p page) Links() map[string]error {
	return p.links
}
