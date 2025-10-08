package main

import (
	"time"

	tls_client "github.com/bogdanfinn/tls-client"
)

type tls_http_client_factory struct {
}

func newTlsHttpClientFactory() *tls_http_client_factory {
	return &tls_http_client_factory{}
}

func (*tls_http_client_factory) Create(o httpClientOptions) httpClient {
	opts := []tls_client.HttpClientOption{
		tls_client.WithTransportOptions(&tls_client.TransportOptions{
			MaxConnsPerHost: o.MaxConnectionsPerHost,
			ReadBufferSize:  o.BufferSize,
		}),
		tls_client.WithTimeoutSeconds(int(o.Timeout / time.Second)),
	}

	if o.SkipTLSVerification {
		opts = append(opts, tls_client.WithInsecureSkipVerify())
	}

	client, err := tls_client.NewHttpClient(tls_client.NewLogger(), opts...)
	if err != nil {
		//todo
	}

	return newTlsHttpClient(
		client,
		o.Header,
	)
}
