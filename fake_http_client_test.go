package main

import (
	"net/url"
)

type fakeHttpClient struct {
	handler func(*url.URL) (*fakeHttpResponse, error)
}

func newFakeHttpClient(h func(*url.URL) (*fakeHttpResponse, error)) *fakeHttpClient {
	return &fakeHttpClient{h}
}

func (c *fakeHttpClient) Get(u *url.URL) (httpResponse, error) {
	return c.handler(u)
}
