package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

const maxRetryDelay = 10 * time.Second
const retryBackoff = 2

type retryHttpClient struct {
	client       httpClient
	maxCount     uint
	initialDelay time.Duration
}

func newRetryHttpClient(c httpClient, maxCount uint, initialDelay time.Duration) httpClient {
	return &retryHttpClient{c, maxCount, initialDelay}
}

func (c *retryHttpClient) Get(u *url.URL, header http.Header) (httpResponse, error) {
	delay := c.initialDelay

	for range c.maxCount + 1 {
		r, err := c.client.Get(u, header)
		if err == nil {
			return r, nil
		} else if e, ok := err.(net.Error); !ok || !e.Timeout() {
			return nil, err
		}

		time.Sleep(delay)
		delay = min(retryBackoff*delay, maxRetryDelay)
	}

	return nil, fmt.Errorf("max retry count %d exceeded", c.maxCount)
}
