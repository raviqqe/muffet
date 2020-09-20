package main

import (
	"time"
)

type linkFetcherOptions struct {
	Concurrency     int
	Headers         map[string]string
	IgnoreFragments bool
	MaxRedirections int
	Timeout         time.Duration
}

func (o *linkFetcherOptions) Initialize() {
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
