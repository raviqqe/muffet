package main

import (
	"io"
	"net/url"

	fh "github.com/bogdanfinn/fhttp"
)

type tlsHttpResponse struct {
	url      *url.URL
	response *fh.Response
}

func newTlsHttpResponse(url *url.URL, response *fh.Response) httpResponse {
	return &tlsHttpResponse{url, response}
}

func (r *tlsHttpResponse) Body() ([]byte, error) {
	defer r.response.Body.Close()
	return io.ReadAll(r.response.Body)
}

func (r *tlsHttpResponse) Header(str string) string {
	return r.response.Header.Get(str)
}

func (r *tlsHttpResponse) StatusCode() int {
	return r.response.StatusCode
}

func (r *tlsHttpResponse) URL() string {
	return r.url.String()
}
