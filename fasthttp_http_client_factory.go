package main

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

type fasthttpHttpClientFactory struct {
}

func newFasthttpHttpClientFactory() *fasthttpHttpClientFactory {
	return &fasthttpHttpClientFactory{}
}

func (*fasthttpHttpClientFactory) Create(o httpClientOptions) httpClient {
	d := func(addr string) (net.Conn, error) {
		return fasthttp.DialTimeout(addr, tcpTimeout)
	}

	if o.Proxy != "" {
		d = fasthttpproxy.FasthttpHTTPDialerTimeout(o.Proxy, tcpTimeout)
	}

	if o.CustomDnsAddr != "" {
		d = func(addr string) (net.Conn, error) {
			dialer := fasthttp.TCPDialer{
				Concurrency: 1000,
				Resolver: &net.Resolver{
					PreferGo:     true,
					StrictErrors: false,
					Dial: func(ctx context.Context, network string, address string) (net.Conn, error) {
						internalDialer := net.Dialer{}
						return internalDialer.DialContext(ctx, "udp", o.CustomDnsAddr)
					},
				},
			}
			return dialer.DialTimeout(addr, tcpTimeout)
		}
	}

	return newFasthttpHttpClient(
		&fasthttp.Client{
			MaxConnsPerHost: o.MaxConnectionsPerHost,
			ReadBufferSize:  o.BufferSize,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: o.SkipTLSVerification,
			},
			Dial:                     d,
			DisablePathNormalizing:   true,
			NoDefaultUserAgentHeader: true,
			MaxResponseBodySize:      o.MaxResponseBodySize,
		},
		o.Timeout,
		o.Header,
	)
}
