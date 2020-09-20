package main

import (
	"net/url"

	"github.com/temoto/robotstxt"
)

type urlValidator struct {
	hostname    string
	sitemapURLs map[string]struct{}
	robotsTxt   *robotstxt.RobotsData
}

func newURLValidator(hostname string, robotsTxt *robotstxt.RobotsData, sitemap map[string]struct{}) urlValidator {
	return urlValidator{hostname, sitemap, robotsTxt}
}

func (i urlValidator) Validate(u *url.URL) bool {
	if len(i.sitemapURLs) != 0 {
		if _, ok := i.sitemapURLs[u.String()]; !ok {
			return false
		}
	}

	if i.robotsTxt != nil && !i.robotsTxt.TestAgent(u.Path, "muffet") {
		return false
	}

	return u.Hostname() == i.hostname
}
