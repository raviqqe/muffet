package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

type retryHttpClient struct {
	client       httpClient
	maxCount     uint
	initialDelay time.Duration
}

func newRetryHttpClient(c httpClient, maxCount uint, initialDelay time.Duration) httpClient {
	return &retryHttpClient{c, maxCount, initialDelay}
}

func (c *retryHttpClient) Get(u *url.URL, header http.Header) (httpResponse, error) {
	d := c.initialDelay

	for range c.maxCount + 1 {
		if r, err := c.client.Get(u, header); err == nil {
			return r, nil
		} else if e, ok := err.(net.Error); !ok || !e.Timeout() {
			return nil, err
		}

		time.Sleep(d)
		d = min(retryBackoff*d, maxRetryDelay)
	}

	return nil, fmt.Errorf("max retry count %d exceeded", c.maxCount)
}
