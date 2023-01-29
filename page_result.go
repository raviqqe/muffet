package main

import "time"

type pageResult struct {
	URL                string
	SuccessLinkResults []*successLinkResult
	ErrorLinkResults   []*errorLinkResult
	Elapsed            time.Duration
}

type successLinkResult struct {
	URL        string
	StatusCode int
	Elapsed    time.Duration
}

type errorLinkResult struct {
	URL     string
	Error   error
	Elapsed time.Duration
}

func (r *pageResult) OK() bool {
	return len(r.ErrorLinkResults) == 0
}
