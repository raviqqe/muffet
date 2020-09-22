package main

import (
	"net/url"
)

type httpClient interface {
	Get(url *url.URL, headers map[string]string) (httpResponse, error)
}
