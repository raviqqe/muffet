package main

type fetcherOptions struct {
	Concurrency         int
	IgnoreFragments     bool
	MaxRedirections     int
	SkipTLSVerification bool
}

func (o *fetcherOptions) Initialize() {
	if o.Concurrency <= 0 {
		o.Concurrency = defaultConcurrency
	}

	if o.MaxRedirections <= 0 {
		o.MaxRedirections = defaultMaxRedirections
	}
}
