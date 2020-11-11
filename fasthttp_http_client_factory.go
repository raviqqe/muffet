package main

import (
	"crypto/tls"
	"net"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

type fasthttpHTTPClientFactory struct {
}

func newFasthttpHTTPClientFactory() *fasthttpHTTPClientFactory {
	return &fasthttpHTTPClientFactory{}
}

func (*fasthttpHTTPClientFactory) Create(o httpClientOptions) httpClient {
	d := func(addr string) (net.Conn, error) {
		return fasthttp.DialTimeout(addr, tcpTimeout)
	}

	if o.Proxy != "" {
		d = fasthttpproxy.FasthttpHTTPDialerTimeout(o.Proxy, tcpTimeout)
	}

	return newFasthttpHTTPClient(
		&fasthttp.Client{
			MaxConnsPerHost: o.MaxConnectionsPerHost,
			ReadBufferSize:  o.BufferSize,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: o.SkipTLSVerification,
			},
			Dial: d,
		},
		o.MaxRedirections,
		o.Timeout,
	)
}
