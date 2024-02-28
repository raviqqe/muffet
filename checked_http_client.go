package main

import (
	"fmt"
	"net/http"
	"net/url"
)

type checkedHttpClient struct {
	client              httpClient
	acceptedStatusCodes statusCodeSet
}

func newCheckedHttpClient(c httpClient, acceptedStatusCodes statusCodeSet) httpClient {
	return &checkedHttpClient{c, acceptedStatusCodes}
}

func (c *checkedHttpClient) Get(u *url.URL, header http.Header) (httpResponse, error) {
	r, err := c.client.Get(u, header)
	if err != nil {
		return nil, err
	}

	code := r.StatusCode()

	if c.acceptedStatusCodes.Contains(code) {
		return r, nil
	} else {
		return nil, fmt.Errorf("%v", code)
	}
}
