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
	atom.Img:    {"src"},
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
		ss = append(ss, parseSrcSet(s)...)
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

// parseSrcSet extracts the URLs of an image candidate list following the
// syntax of the `srcset` attribute defined by the HTML specification. URLs are
// separated by commas but may themselves contain commas, so each candidate is
// a run of non-whitespace characters whose trailing comma, if any, acts as the
// separator and whose optional descriptor is discarded.
func parseSrcSet(s string) []string {
	ss := []string{}

	for i := 0; i < len(s); {
		for i < len(s) && (isAsciiWhitespace(s[i]) || s[i] == ',') {
			i++
		}

		start := i

		for i < len(s) && !isAsciiWhitespace(s[i]) {
			i++
		}

		if start == i {
			continue
		}

		u := s[start:i]

		if strings.HasSuffix(u, ",") {
			ss = append(ss, strings.TrimRight(u, ","))
			continue
		}

		ss = append(ss, u)

		for depth := 0; i < len(s) && !(depth == 0 && s[i] == ','); i++ {
			switch s[i] {
			case '(':
				depth++
			case ')':
				if depth > 0 {
					depth--
				}
			}
		}
	}

	return ss
}

func isAsciiWhitespace(b byte) bool {
	switch b {
	case '\t', '\n', '\f', '\r', ' ':
		return true
	}

	return false
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
