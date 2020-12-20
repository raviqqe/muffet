package main

import "time"

type httpClientOptions struct {
	MaxConnectionsPerHost,
	BufferSize,
	MaxRedirections int
	Proxy                    string
	SkipTLSVerification      bool
	Timeout                  time.Duration
	NoDefaultUserAgentHeader bool
}

type httpClientFactory interface {
	Create(httpClientOptions) httpClient
}
