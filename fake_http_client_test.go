package main

import (
	"net/url"
)

type fakeHTTPClient struct {
	handler func(*url.URL) (*fakeHTTPResponse, error)
}

func newFakeHTTPClient(h func(*url.URL) (*fakeHTTPResponse, error)) *fakeHTTPClient {
	return &fakeHTTPClient{h}
}

func (c *fakeHTTPClient) Get(u *url.URL, headers map[string]string) (httpResponse, error) {
	return c.handler(u)
}
