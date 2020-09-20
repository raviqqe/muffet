package main

import (
	"net/url"

	"github.com/temoto/robotstxt"
)

type urlValidator struct {
	hostname    string
	sitemapURLs map[string]struct{}
	robotsData  *robotstxt.RobotsData
}

func newURLValidator(hostname string, robotsData *robotstxt.RobotsData, sitemap map[string]struct{}) urlValidator {
	return urlValidator{hostname, sitemap, robotsData}
}

func (i urlValidator) Validate(u *url.URL) bool {
	if i.sitemapURLs != nil {
		if _, ok := i.sitemapURLs[u.String()]; !ok {
			return false
		}
	}

	if i.robotsData != nil && !i.robotsData.TestAgent(u.Path, "muffet") {
		return false
	}

	return u.Hostname() == i.hostname
}
