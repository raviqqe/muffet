package main

import (
	"net/http"
	"net/url"
)

type throttledHttpClient struct {
	client            httpClient
	connections       semaphore
	hostThrottlerPool *hostThrottlerPool
}

func newThrottledHttpClient(c httpClient, requestPerSecond int, maxConnections, maxConnectionsPerHost int) httpClient {
	return &throttledHttpClient{
		c,
		newSemaphore(maxConnections),
		newHostThrottlerPool(requestPerSecond, maxConnectionsPerHost),
	}
}

func (c *throttledHttpClient) Get(u *url.URL, header http.Header) (httpResponse, error) {
	c.connections.Request()
	defer c.connections.Release()

	t := c.hostThrottlerPool.Get(u.Hostname())
	t.Request()
	defer t.Release()

	return c.client.Get(u, header)
}
