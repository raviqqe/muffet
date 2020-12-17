package main

import (
	"github.com/valyala/fasthttp"
)

type fasthttpHTTPResponse struct {
	url      *fasthttp.URI
	response *fasthttp.Response
}

func newFasthttpHTTPResponse(u *fasthttp.URI, r *fasthttp.Response) httpResponse {
	return fasthttpHTTPResponse{u, r}
}

func (r fasthttpHTTPResponse) URL() string {
	return r.url.String()
}

func (r fasthttpHTTPResponse) StatusCode() int {
	return r.response.StatusCode()
}

func (r fasthttpHTTPResponse) Header(key string) string {
	return string(r.response.Header.Peek(key))
}

func (r fasthttpHTTPResponse) Body() ([]byte, error) {
	switch string(r.response.Header.Peek("Content-Encoding")) {
	case "gzip":
		return r.response.BodyGunzip()
	}

	return r.response.Body(), nil
}
