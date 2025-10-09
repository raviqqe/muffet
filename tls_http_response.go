package main

import (
	"errors"
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
	var err error
	defer func() {
		if cerr := r.response.Body.Close(); cerr != nil {
			if err != nil {
				err = errors.Join(err, cerr)
			} else {
				err = cerr
			}
		}
	}()
	b, err := io.ReadAll(r.response.Body)
	return b, err
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
