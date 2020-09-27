package main

import "net/url"

type throttledHTTPClient struct {
	client    httpClient
	semaphore semaphore
}

func newThrottledHTTPClient(c httpClient, maxConnections int) httpClient {
	return &throttledHTTPClient{c, newSemaphore(maxConnections)}
}

func (c *throttledHTTPClient) Get(u *url.URL, hs map[string]string) (httpResponse, error) {
	c.semaphore.Request()
	defer c.semaphore.Release()

	return c.client.Get(u, hs)
}
