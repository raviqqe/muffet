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

func (v urlValidator) Validate(u *url.URL) bool {
	if v.sitemapURLs != nil {
		if _, ok := v.sitemapURLs[u.String()]; !ok {
			return false
		}
	}

	if v.robotsData != nil && !v.robotsData.TestAgent(u.Path, "muffet") {
		return false
	}

	return u.Hostname() == v.hostname
}
