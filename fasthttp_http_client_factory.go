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
	} else if o.DnsResolver != "" {
		nd := &net.Dialer{}
		td := &fasthttp.TCPDialer{
			Concurrency: concurrency,
			Resolver: &net.Resolver{
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					return nd.DialContext(ctx, "udp", o.DnsResolver)
				},
			},
		}

		d = func(address string) (net.Conn, error) {
			return td.DialTimeout(address, tcpTimeout)
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
