package main

import (
	"net/http"
	"net/url"
)

type retryHttpClient struct {
	client   httpClient
	maxCount uint
}

func newRetryHttpClient(c httpClient, maxCount uint) httpClient {
	return &retryHttpClient{c, maxCount}
}

func (c *retryHttpClient) Get(u *url.URL, header http.Header) (httpResponse, error) {
	return c.client.Get(u, header)
}
