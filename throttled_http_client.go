package main

import (
	"net/url"

	"go.uber.org/ratelimit"
)

// TODO Throttle requests for each host.
type throttledHttpClient struct {
	client    httpClient
	limiter   ratelimit.Limiter
	semaphore semaphore
}

func newThrottledHttpClient(c httpClient, rps int, maxConnections int) httpClient {
	l := ratelimit.NewUnlimited()

	if rps > 0 {
		l = ratelimit.New(rps)
	}

	return &throttledHttpClient{c, l, newSemaphore(maxConnections)}
}

func (c *throttledHttpClient) Get(u *url.URL) (httpResponse, error) {
	c.semaphore.Request()
	defer c.semaphore.Release()

	c.limiter.Take()

	return c.client.Get(u)
}
