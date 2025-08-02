package main

import (
	"fmt"
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
	for range c.maxCount + 1 {
		if r, err := c.client.Get(u, header); err == nil {
			return r, nil
		}

	}

	return nil, fmt.Errorf("max retry count %d exceeded", c.maxCount)
}
