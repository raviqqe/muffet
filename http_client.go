package main

import (
	"net/url"
	"time"
)

type httpClient interface {
	Get(url *url.URL, headers map[string]string, timeout time.Duration) (httpResponse, error)
}

type httpResponse interface {
	StatusCode() int
	Header(string) string
	Body() []byte
}
