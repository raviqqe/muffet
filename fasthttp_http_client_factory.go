package main

import (
	"crypto/tls"

	"github.com/valyala/fasthttp"
)

type fasthttpHTTPClientFactory struct {
}

func newFasthttpHTTPClientFactory() *fasthttpHTTPClientFactory {
	return &fasthttpHTTPClientFactory{}
}

func (*fasthttpHTTPClientFactory) Create(o httpClientOptions) httpClient {
	return newFasthttpHTTPClient(
		&fasthttp.Client{
			MaxConnsPerHost: o.Concurrency,
			ReadBufferSize:  o.BufferSize,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: o.SkipTLSVerification,
			},
		},
		o.MaxRedirections,
		o.Timeout,
	)
}
