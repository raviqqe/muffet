package main

import (
	"errors"
	"net/url"
	"time"
)

type fakeHTTPClient struct {
	data map[string]*fakeHTTPResponse
}

func newFakeHTTPClient(data map[string]*fakeHTTPResponse) httpClient {
	return &fakeHTTPClient{data}
}

func (c *fakeHTTPClient) Get(u *url.URL, headers map[string]string, timeout time.Duration) (httpResponse, error) {
	if r, ok := c.data[u.String()]; ok {
		return r, nil
	}

	return nil, errors.New("fake http request failed")
}
