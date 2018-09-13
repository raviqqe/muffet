package main

type checkerOptions struct {
	fetcherOptions
	FollowRobotsTxt,
	FollowSitemapXML,
	SkipTLSVerification bool
}
