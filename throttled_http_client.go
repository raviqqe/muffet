package main

import "net/url"

type throttledHTTPClient struct {
	client    httpClient
	semaphore semaphore
}

func newThrottledHTTPClient(c httpClient, concurrency int) httpClient {
	if concurrency < 1 {
		concurrency = defaultConcurrency
	}

	return &throttledHTTPClient{c, newSemaphore(concurrency)}
}

func (c *throttledHTTPClient) Get(u *url.URL, hs map[string]string) (httpResponse, error) {
	c.semaphore.Request()
	defer c.semaphore.Release()

	return c.client.Get(u, hs)
}
