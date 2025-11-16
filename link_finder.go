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

var atomToAttributes = map[atom.Atom][]string{
	atom.A:      {"href"},
	atom.Frame:  {"src"},
	atom.Iframe: {"src"},
	atom.Img:    {"src"},
	atom.Link:   {"href"},
	atom.Script: {"src"},
	atom.Source: {"src", "srcset"},
	atom.Track:  {"src"},
	atom.Meta:   {"content"},
}

var imageDescriptorPattern = regexp.MustCompile(`(\S)\s+\S+\s*$`)

type linkFinder struct {
	linkFilterer linkFilterer
}

func newLinkFinder(f linkFilterer) linkFinder {
	return linkFinder{f}
}

func (f linkFinder) Find(n *html.Node, base *url.URL) map[string]error {
	ls := map[string]error{}

	for _, n := range scrape.FindAllNested(n, func(n *html.Node) bool {
		_, ok := atomToAttributes[n.DataAtom]
		return ok
	}) {

		// `preconnect` and `dns-prefetch` links are not HTTP resources.
		if n.DataAtom == atom.Link {

			if rel := scrape.Attr(n, "rel"); rel == "preconnect" || rel == "dns-prefetch" {
				continue
			}
		}

		for _, a := range atomToAttributes[n.DataAtom] {
			ss := f.parseLinks(n, a)

			for _, s := range ss {
				s := f.trimUrl(s)

				if s == "" {
					continue
				}

				u, err := url.Parse(s)
				if err != nil {
					ls[s] = err
					continue
				}

				u = base.ResolveReference(u)

				if f.linkFilterer.IsValid(u) {
					ls[u.String()] = nil
				}
			}
		}
	}

	return ls
}

func (f linkFinder) parseLinks(n *html.Node, a string) []string {
	s := scrape.Attr(n, a)
	ss := []string{}

	switch a {
	case "srcset":
		for s := range strings.SplitSeq(s, ", ") {
			ss = append(ss, f.trimUrl(imageDescriptorPattern.ReplaceAllString(s, "$1")))
		}
	case "content":
		switch scrape.Attr(n, "property") {
		case "og:image", "og:audio", "og:video", "og:image:url", "og:image:secure_url", "twitter:image":
			ss = append(ss, s)
		}
	default:
		ss = append(ss, s)
	}

	return ss
}

func (linkFinder) trimUrl(s string) string {
	s = strings.TrimSpace(s)

	if !strings.HasPrefix(s, "data:") {
		return s
	}

	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}

		return r
	}, s)
}
