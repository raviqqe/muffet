package main

type linkFetcherOptions struct {
	// TODO Move to throttledHTTPClient.
	Concurrency int
	// TODO Move to httpClient?
	Headers         map[string]string
	IgnoreFragments bool
}

func (o *linkFetcherOptions) Initialize() {
	if o.Concurrency <= 0 {
		o.Concurrency = defaultConcurrency
	}
}
