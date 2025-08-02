package main

import (
	"net/http"
	"net/url"
)

type retryHttpClient struct {
	client  httpClient
	retries uint
}

func newRetryHttpClient(c httpClient, retries uint) httpClient {
	return &retryHttpClient{c, retries}
}

func (c *retryHttpClient) Get(u *url.URL, header http.Header) (httpResponse, error) {
	return c.client.Get(u, header)
}
