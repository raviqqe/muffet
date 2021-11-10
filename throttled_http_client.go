package main

import (
	"net/url"

	"go.uber.org/ratelimit"
)

// TODO Throttle requests for each host.
type throttledHTTPClient struct {
	client    httpClient
	limiter   ratelimit.Limiter
	semaphore semaphore
}

func newThrottledHTTPClient(c httpClient, rps int, maxConnections int) httpClient {
	l := ratelimit.NewUnlimited()

	if rps > 0 {
		l = ratelimit.New(rps)
	}

	return &throttledHTTPClient{c, l, newSemaphore(maxConnections)}
}

func (c *throttledHTTPClient) Get(u *url.URL) (httpResponse, error) {
	c.semaphore.Request()
	defer c.semaphore.Release()

	c.limiter.Take()

	return c.client.Get(u)
}
