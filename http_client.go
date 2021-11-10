package main

import (
	"net/url"
)

type httpClient interface {
	// Get sends a HTTP request with the GET method.
	// Non-200 status codes always result in non-nil errors.
	Get(url *url.URL) (httpResponse, error)
}
