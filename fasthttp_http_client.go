package main

import (
	"net/url"
	"time"

	"github.com/valyala/fasthttp"
)

type fasthttpHTTPClient struct {
	client *fasthttp.Client
}

func newFasthttpHTTPClient(c *fasthttp.Client) httpClient {
	return fasthttpHTTPClient{c}
}

func (c fasthttpHTTPClient) Get(u *url.URL, headers map[string]string, timeout time.Duration) (httpResponse, error) {
	req, res := fasthttp.Request{}, fasthttp.Response{}
	req.SetRequestURI(u.String())
	req.SetConnectionClose()

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	err := c.client.DoTimeout(&req, &res, timeout)
	if err != nil {
		return nil, err
	}

	return newFasthttpHTTPResponse(&res), nil
}
