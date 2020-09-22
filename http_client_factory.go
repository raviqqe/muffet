package main

import "time"

type httpClientOptions struct {
	Concurrency,
	BufferSize,
	MaxRedirections int
	SkipTLSVerification bool
	Timeout             time.Duration
}

type httpClientFactory interface {
	Create(httpClientOptions) httpClient
}
