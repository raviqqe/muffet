package main

import (
	"errors"
	"net/url"
	"time"
)

type fakeHTTPClient struct{}

func newFakeHTTPClient() httpClient {
	return &fakeHTTPClient{}
}

func (c *fakeHTTPClient) Get(u *url.URL, headers map[string]string, timeout time.Duration) (httpResponse, error) {
	return nil, errors.New("fake http request failed")
}
