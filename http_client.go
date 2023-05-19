package main

import (
	"net/http"
	"net/url"
)

type httpClient interface {
	// Get sends an HTTP request with a GET method.
	// It depends on implementation of each client what is considered as errors.
	Get(url *url.URL, header http.Header) (httpResponse, error)
}
