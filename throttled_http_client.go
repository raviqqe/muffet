package main

import (
	"net/url"
)

// TODO Throttle requests for each host.
type throttledHttpClient struct {
	client                httpClient
	connections           semaphore
	hostThrottlers        map[string]*hostThrottler
	requestPerSecond      int
	maxConnectionsPerHost int
}

func newThrottledHttpClient(c httpClient, requestPerSecond int, maxConnections, maxConnectionsPerHost int) httpClient {
	return &throttledHttpClient{
		c,
		newSemaphore(maxConnections),
		map[string]*hostThrottler{},
		requestPerSecond,
		maxConnectionsPerHost,
	}
}

func (c *throttledHttpClient) Get(u *url.URL) (httpResponse, error) {
	c.connections.Request()
	defer c.connections.Release()

	t := c.getHostThrottler(u.Hostname())
	t.Request()
	defer t.Release()

	return c.client.Get(u)
}

func (c *throttledHttpClient) getHostThrottler(name string) *hostThrottler {
	t, ok := c.hostThrottlers[name]

	if !ok {
		t = newHostThrottler(c.requestPerSecond, c.maxConnectionsPerHost)
		c.hostThrottlers[name] = t
	}

	return t
}
