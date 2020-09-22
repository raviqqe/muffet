package main

import (
	"net/url"
	"regexp"
	"strings"
	"unicode"

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

type linkFinder struct {
	excludedPatterns []*regexp.Regexp
}

func newLinkFinder(rs []*regexp.Regexp) linkFinder {
	return linkFinder{rs}
}

func (f linkFinder) Find(n *html.Node, base *url.URL) map[string]error {
	ls := map[string]error{}

	for _, n := range scrape.FindAllNested(n, func(n *html.Node) bool {
		_, ok := atomToAttributes[n.DataAtom]
		return ok
	}) {
		for _, a := range atomToAttributes[n.DataAtom] {
			s := normalizeURL(scrape.Attr(n, a))

			if s == "" || f.isLinkExcluded(s) {
				continue
			}

			u, err := url.Parse(s)
			if err != nil {
				ls[s] = err
				continue
			} else if _, ok := validSchemes[u.Scheme]; ok {
				ls[base.ResolveReference(u).String()] = nil
			}
		}
	}

	return ls
}

func (f linkFinder) isLinkExcluded(u string) bool {
	for _, r := range f.excludedPatterns {
		if r.MatchString(u) {
			return true
		}
	}

	return false
}

func normalizeURL(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}

		return r
	}, s)
}
