package main

import (
	"net/url"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

type fasthttpHttpClient struct {
	client  *fasthttp.Client
	timeout time.Duration
	headers map[string]string
}

func newFasthttpHttpClient(c *fasthttp.Client, timeout time.Duration, headers map[string]string) httpClient {
	return &fasthttpHttpClient{c, timeout, headers}
}

func (c *fasthttpHttpClient) Get(u *url.URL) (httpResponse, error) {
	req, res := fasthttp.Request{}, fasthttp.Response{}
	req.SetRequestURI(u.String())

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}

	// Some HTTP servers require "Accept" headers to be set explicitly.
	if !includeHeader(c.headers, "Accept") {
		req.Header.Add("Accept", "*/*")
	}

	err := c.client.DoTimeout(&req, &res, c.timeout)
	if err != nil {
		return nil, err
	}

	return newFasthttpHttpResponse(req.URI(), &res), nil
}

func includeHeader(hs map[string]string, h string) bool {
	for k := range hs {
		if strings.EqualFold(k, h) {
			return true
		}
	}

	return false
}
