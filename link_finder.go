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

var imageDescriptorPattern = regexp.MustCompile(" [^ ]*$")

type linkFinder struct {
	excludedPatterns []*regexp.Regexp
	includedPatterns []*regexp.Regexp
}

func newLinkFinder(es []*regexp.Regexp, is []*regexp.Regexp) linkFinder {
	return linkFinder{excludedPatterns: es, includedPatterns: is}
}

func (f linkFinder) Find(n *html.Node, base *url.URL) map[string]error {
	ls := map[string]error{}

	for _, n := range scrape.FindAllNested(n, func(n *html.Node) bool {
		_, ok := atomToAttributes[n.DataAtom]
		return ok
	}) {
		for _, a := range atomToAttributes[n.DataAtom] {
			ss := f.parseLinks(n, a)

			for _, s := range ss {
				s := strings.TrimSpace(s)

				if s == "" {
					continue
				}

				u, err := url.Parse(s)
				if err != nil {
					ls[s] = err
					continue
				}

				s = base.ResolveReference(u).String()

				if _, ok := validSchemes[u.Scheme]; ok && !f.isLinkExcluded(s) && f.isLinkIncluded(s) {
					ls[s] = nil
				}
			}
		}
	}

	return ls
}

func (linkFinder) parseLinks(n *html.Node, a string) []string {
	s := scrape.Attr(n, a)
	ss := []string{}

	if a == "srcset" {
		for _, s := range strings.Split(s, ",") {
			ss = append(ss, imageDescriptorPattern.ReplaceAllString(strings.TrimSpace(s), ""))
		}
	} else {
		ss = append(ss, s)
	}

	return ss
}

func (f linkFinder) isLinkExcluded(u string) bool {
	return f.matches(u, f.excludedPatterns)
}

func (f linkFinder) isLinkIncluded(u string) bool {
	return len(f.includedPatterns) == 0 || f.matches(u, f.includedPatterns)
}

func (f linkFinder) matches(u string, rs []*regexp.Regexp) bool {
	for _, r := range rs {
		if r.MatchString(u) {
			return true
		}
	}

	return false
}
