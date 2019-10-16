package main

type checkerOptions struct {
	fetcherOptions
	FollowRobotsTxt,
	FollowSitemapXML,
	FollowURLParams,
	SkipTLSVerification bool
}
