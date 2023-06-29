package main

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
)

var validSchemes = map[string]struct{}{
	"":      {},
	"http":  {},
	"https": {},
}

type linkFilterer struct {
	excludedPatterns []*regexp.Regexp
	includedPatterns []*regexp.Regexp
}

func newLinkFilterer(es []*regexp.Regexp, is []*regexp.Regexp) linkFilterer {
	return linkFilterer{excludedPatterns: es, includedPatterns: is}
}

func (f linkFilterer) Filter(u *url.URL) bool {
	if _, ok := validSchemes[u.Scheme]; ok && !f.isLinkExcluded(s) && f.isLinkIncluded(s) {
		return true
	}

	return false
}

func (linkFilterer) parseLinks(n *html.Node, a string) []string {
	s := scrape.Attr(n, a)
	ss := []string{}

	switch a {
	case "srcset":
		for _, s := range strings.Split(s, ",") {
			ss = append(ss, imageDescriptorPattern.ReplaceAllString(strings.TrimSpace(s), ""))
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

func (f linkFilterer) isLinkExcluded(u string) bool {
	return f.matches(u, f.excludedPatterns)
}

func (f linkFilterer) isLinkIncluded(u string) bool {
	return len(f.includedPatterns) == 0 || f.matches(u, f.includedPatterns)
}

func (f linkFilterer) matches(u string, rs []*regexp.Regexp) bool {
	for _, r := range rs {
		if r.MatchString(u) {
			return true
		}
	}

	return false
}
