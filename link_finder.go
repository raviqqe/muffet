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
			s := scrape.Attr(n, a)
			ss := []string{}

			if a == "srcset" {
				for _, s := range strings.Split(s, ",") {
					ss = append(ss, imageDescriptorPattern.ReplaceAllString(strings.TrimSpace(s), ""))
				}
			} else {
				ss = append(ss, s)
			}

			for _, s := range ss {
				s := strings.TrimSpace(s)

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
