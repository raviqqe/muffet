package main

import (
	"crypto/tls"
	"net"

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
			Dial: func(addr string) (net.Conn, error) {
				return fasthttp.DialTimeout(addr, tcpTimeout)
			}},
		o.MaxRedirections,
		o.Timeout,
	)
}
