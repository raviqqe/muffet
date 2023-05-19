package main

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

type fasthttpHttpClient struct {
	client  *fasthttp.Client
	timeout time.Duration
	header  http.Header
}

func newFasthttpHttpClient(c *fasthttp.Client, timeout time.Duration, header http.Header) httpClient {
	return &fasthttpHttpClient{c, timeout, header}
}

func (c *fasthttpHttpClient) Get(u *url.URL, headers http.Header) (httpResponse, error) {
	req, res := fasthttp.Request{}, fasthttp.Response{}
	req.SetRequestURI(u.String())
	req.SetConnectionClose()

	for k, vs := range c.header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	for k, vs := range headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	// S me HTTP servers require "Accept" headers set explicitly.
	if !includeHeader(c.header, "Accept") {
		req.Header.Add("Accept", "*/*")
	}

	err := c.client.DoTimeout(&req, &res, c.timeout)
	if err != nil {
		return nil, err
	}

	return newFasthttpHttpResponse(req.URI(), &res), nil
}

func includeHeader(h http.Header, k string) bool {
	for kk := range h {
		if strings.EqualFold(kk, k) {
			return true
		}
	}

	return false
}
