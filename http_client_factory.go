package main

import (
	"net/http"
	"time"
)

type httpClientOptions struct {
	MaxConnectionsPerHost,
	MaxResponseBodySize,
	BufferSize int
	Proxy               string
	SkipTLSVerification bool
	Timeout             time.Duration
	Header              http.Header
	DnsResolver         string
}

type httpClientFactory interface {
	Create(httpClientOptions) httpClient
}
