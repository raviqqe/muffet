package main

import (
	"net/url"
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
	atom.Img:    {"src", "srcset"},
	atom.Link:   {"href"},
	atom.Script: {"src"},
	atom.Source: {"src", "srcset"},
	atom.Track:  {"src"},
	atom.Meta:   {"content"},
}

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
		for _, s := range parseSrcSet(s) {
			ss = append(ss, f.trimUrl(s))
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

// parseSrcSet returns the URLs of image candidates in a srcset attribute
// following https://html.spec.whatwg.org/multipage/images.html#parsing-a-srcset-attribute.
// A URL ends at whitespace, so it may contain commas. A comma after a URL or
// a descriptor separates candidates, with or without whitespace around it.
func parseSrcSet(s string) []string {
	ss := []string{}

	for {
		s = strings.TrimLeft(s, srcSetSpaces+",")

		if s == "" {
			return ss
		}

		i := strings.IndexAny(s, srcSetSpaces)
		if i < 0 {
			i = len(s)
		}

		u := s[:i]
		s = s[i:]

		if t := strings.TrimRight(u, ","); t != u {
			ss = append(ss, t)
			continue
		}

		ss = append(ss, u)
		s = skipSrcSetDescriptors(strings.TrimLeft(s, srcSetSpaces))
	}
}

const srcSetSpaces = " \t\n\f\r"

// skipSrcSetDescriptors skips descriptors up to and including the comma that
// ends the current image candidate. Commas in parentheses do not end it.
func skipSrcSetDescriptors(s string) string {
	inParens := false

	for i, r := range s {
		switch {
		case inParens:
			inParens = r != ')'
		case r == '(':
			inParens = true
		case r == ',':
			return s[i+1:]
		}
	}

	return ""
}
