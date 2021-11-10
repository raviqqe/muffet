package main

import (
	"net/url"

	"go.uber.org/ratelimit"
)

// TODO Throttle requests for each host.
type throttledHttpClient struct {
	client                httpClient
	connections           semaphore
	hosts                 map[string]*Host
	requestPerSecond      int
	maxConnectionsPerHost int
}

type Host struct {
	limiter     ratelimit.Limiter
	connections semaphore
}

func newThrottledHttpClient(c httpClient, requestPerSecond int, maxConnections, maxConnectionsPerHost int) httpClient {
	return &throttledHttpClient{
		c,
		newSemaphore(maxConnections),
		map[string]*Host{},
		requestPerSecond,
		maxConnectionsPerHost,
	}
}

func (c *throttledHttpClient) Get(u *url.URL) (httpResponse, error) {
	c.connections.Request()
	defer c.connections.Release()

	h := c.getHost(u.Hostname())
	h.limiter.Take()
	h.connections.Request()
	defer h.connections.Release()

	return c.client.Get(u)
}

func (c *throttledHttpClient) getHost(name string) *Host {
	h, ok := c.hosts[name]

	if !ok {
		h = &Host{c.createLimiter(), newSemaphore(c.maxConnectionsPerHost)}
		c.hosts[name] = h
	}

	return h
}

func (c *throttledHttpClient) createLimiter() ratelimit.Limiter {
	l := ratelimit.NewUnlimited()

	if c.requestPerSecond > 0 {
		l = ratelimit.New(c.requestPerSecond)
	}

	return l
}
