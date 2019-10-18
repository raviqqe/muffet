package main

type checkerOptions struct {
	fetcherOptions
	BufferSize int
	FollowRobotsTxt,
	FollowSitemapXML,
	FollowURLParams,
	SkipTLSVerification bool
}
