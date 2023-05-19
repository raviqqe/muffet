package main

import (
	"net/http"
	"net/url"
)

type fakeHttpClient struct {
	handler func(*url.URL) (*fakeHttpResponse, error)
}

func newFakeHttpClient(h func(*url.URL) (*fakeHttpResponse, error)) *fakeHttpClient {
	return &fakeHttpClient{h}
}

func (c *fakeHttpClient) Get(u *url.URL, _ http.Header) (httpResponse, error) {
	return c.handler(u)
}
