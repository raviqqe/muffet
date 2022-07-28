package main

import "time"

type httpClientOptions struct {
	MaxConnectionsPerHost,
	MaxResponseBodySize,
	BufferSize int
	Proxy               string
	SkipTLSVerification bool
	Timeout             time.Duration
	Headers             map[string]string
}

type httpClientFactory interface {
	Create(httpClientOptions) httpClient
}
