package main

import (
	"net/url"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var validSchemes = map[string]struct{}{
	"":      {},
	"http":  {},
	"https": {},
}

var atomToAttributes = map[atom.Atom][]string{
	atom.A:      {"href"},
	atom.Frame:  {"src"},
	atom.Iframe: {"src"},
	atom.Img:    {"src"},
	atom.Link:   {"href"},
	atom.Script: {"src"},
	atom.Source: {"src", "srcset"},
	atom.Track:  {"src"},
}

func scrapePage(p page) (map[string]bool, map[string]error) {
	bs := map[string]bool{}
	es := map[string]error{}

	for _, n := range scrape.FindAllNested(p.Body(), func(n *html.Node) bool {
		_, ok := atomToAttributes[n.DataAtom]
		return ok
	}) {
		for _, a := range atomToAttributes[n.DataAtom] {
			s := scrape.Attr(n, a)

			if s == "" {
				continue
			}

			u, err := url.Parse(s)

			if err != nil {
				es[s] = err
				continue
			}

			if _, ok := validSchemes[u.Scheme]; !ok {
				continue
			}

			u, err = resolveURL(p, u)

			if err != nil {
				es[s] = err
				continue
			}

			bs[u.String()] = n.DataAtom == atom.A
		}
	}

	return bs, es
}

func resolveURL(p page, u *url.URL) (*url.URL, error) {
	if u.IsAbs() {
		return u, nil
	}

	b := p.URL()

	if n, ok := scrape.Find(p.Body(), func(n *html.Node) bool {
		return n.DataAtom == atom.Base
	}); ok {
		u, err := url.Parse(scrape.Attr(n, "href"))

		if err != nil {
			return nil, err
		}

		b = b.ResolveReference(u)
	}

	u = b.ResolveReference(u)

	return u, nil
}
