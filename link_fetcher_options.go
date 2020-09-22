package main

type linkFetcherOptions struct {
	// TODO Move to httpClient?
	Headers         map[string]string
	IgnoreFragments bool
}
