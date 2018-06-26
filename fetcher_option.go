package main

import (
	"regexp"
	"time"
)

type fetcherOptions struct {
	Concurrency         int
	ExcludedPatterns    []*regexp.Regexp
	Headers             map[string]string
	IgnoreFragments     bool
	MaxRedirections     int
	RemoveNewlines      bool
	SkipTLSVerification bool
	Timeout             time.Duration
}

func (o *fetcherOptions) Initialize() {
	if o.Concurrency <= 0 {
		o.Concurrency = defaultConcurrency
	}

	if o.MaxRedirections <= 0 {
		o.MaxRedirections = defaultMaxRedirections
	}

	if o.Timeout <= 0 {
		o.Timeout = defaultTimeout
	}
}
