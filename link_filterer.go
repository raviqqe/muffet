package main

import (
	"net/url"
	"regexp"
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

func (f linkFilterer) IsValid(u *url.URL) bool {
	s := u.String()

	if _, ok := validSchemes[u.Scheme]; !ok {
		return false
	}

	return !f.isLinkExcluded(s) && f.isLinkIncluded(s)
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
