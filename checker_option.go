package main

type checkerOptions struct {
	Concurrency int
	FollowRobotsTxt,
	FollowSitemapXML,
	IgnoreFragments bool
	MaxRedirections     int
	SkipTLSVerification bool
}

func (o *checkerOptions) Initialize() {
	if o.Concurrency <= 0 {
		o.Concurrency = defaultConcurrency
	}

	if o.MaxRedirections <= 0 {
		o.MaxRedirections = defaultMaxRedirections
	}
}
