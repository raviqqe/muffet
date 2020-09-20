package main

import (
	"net/url"

	"github.com/temoto/robotstxt"
)

type linkValidator struct {
	hostname    string
	sitemapURLs map[string]struct{}
	robotsData  *robotstxt.RobotsData
}

func newLinkValidator(hostname string, robotsData *robotstxt.RobotsData, sitemap map[string]struct{}) linkValidator {
	return linkValidator{hostname, sitemap, robotsData}
}

func (v linkValidator) Validate(u *url.URL) bool {
	if v.sitemapURLs != nil {
		if _, ok := v.sitemapURLs[u.String()]; !ok {
			return false
		}
	}

	if v.robotsData != nil && !v.robotsData.TestAgent(u.Path, agentName) {
		return false
	}

	return u.Hostname() == v.hostname
}
