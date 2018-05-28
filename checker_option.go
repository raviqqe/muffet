package main

import "regexp"

type checkerOptions struct {
	ExcludedPatterns []*regexp.Regexp
	fetcherOptions
	FollowRobotsTxt,
	FollowSitemapXML bool
}
