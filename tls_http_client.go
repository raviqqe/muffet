package main

import (
	"net/http"
	"net/url"

	fh "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

type tlsHttpClient struct {
	client tls_client.HttpClient
	header http.Header
}

func newTlsHttpClient(c tls_client.HttpClient, header http.Header) httpClient {
	return &tlsHttpClient{client: c, header: header}
}

// Get performs an HTTP GET using the underlying tls-client HttpClient and
// adapts the response to the package's httpResponse interface.
func (t *tlsHttpClient) Get(u *url.URL, header http.Header) (httpResponse, error) {
	// Build request using fhttp types (tls-client expects fhttp requests)
	req, err := fh.NewRequest(fh.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	for k, vs := range t.header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	for k, vs := range header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	if !includeHeader(t.header, "Accept") {
		req.Header.Add("Accept", "*/*")
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return newTlsHttpResponse(u, resp), nil
}
