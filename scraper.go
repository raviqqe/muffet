package main

import (
	"net/url"
	"regexp"
	"strings"

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

func (sc scraper) Scrape(n *html.Node, base *url.URL) map[string]error {
	us := map[string]error{}

	for _, n := range scrape.FindAllNested(n, func(n *html.Node) bool {
		_, ok := atomToAttributes[n.DataAtom]
		return ok
	}) {
		for _, a := range atomToAttributes[n.DataAtom] {
			s := scrape.Attr(n, a)

			if s == "" || sc.isURLExcluded(s) {
				continue
			}

			s, err := url.PathUnescape(s)
			s = strings.Replace(s, "\t", "", -1)
			s = strings.Replace(s, "\r", "", -1)
			s = strings.Replace(s, "\n", "", -1)
			u, err := urlParse(s)

			if err != nil {
				us[s] = err
				continue
			}

			if _, ok := validSchemes[u.Scheme]; !ok {
				continue
			}

			us[base.ResolveReference(u).String()] = nil
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
