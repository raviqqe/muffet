package main

import (
	"net/url"
	"regexp"

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

type scraper struct {
	excludedPatterns []*regexp.Regexp
}

func newScraper(rs []*regexp.Regexp) scraper {
	return scraper{rs}
}

func (sc scraper) Scrape(p page) map[string]error {
	us := map[string]error{}

	for _, n := range scrape.FindAllNested(p.Body(), func(n *html.Node) bool {
		_, ok := atomToAttributes[n.DataAtom]
		return ok
	}) {
		for _, a := range atomToAttributes[n.DataAtom] {
			s := scrape.Attr(n, a)

			if s == "" || sc.isURLExcluded(s) {
				continue
			}

			u, err := url.Parse(s)

			if err != nil {
				us[s] = err
				continue
			}

			if _, ok := validSchemes[u.Scheme]; !ok {
				continue
			}

			u, err = resolveURL(p, u)

			if err != nil {
				us[s] = err
				continue
			}

			us[u.String()] = nil
		}
	}

	return us
}

func (sc scraper) isURLExcluded(u string) bool {
	for _, r := range sc.excludedPatterns {
		if r.MatchString(u) {
			return true
		}
	}

	return false
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
