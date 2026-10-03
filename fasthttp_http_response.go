package main

import (
	"strings"

	"github.com/valyala/fasthttp"
)

type fasthttpHttpResponse struct {
	url      *fasthttp.URI
	response *fasthttp.Response
}

func newFasthttpHttpResponse(u *fasthttp.URI, r *fasthttp.Response) httpResponse {
	return fasthttpHttpResponse{u, r}
}

func (r fasthttpHttpResponse) URL() string {
	return r.url.String()
}

func (r fasthttpHttpResponse) StatusCode() int {
	return r.response.StatusCode()
}

func (r fasthttpHttpResponse) Header(key string) string {
	return string(r.response.Header.Peek(key))
}

func (r fasthttpHttpResponse) Body() ([]byte, error) {
	switch strings.ToLower(strings.TrimSpace(string(r.response.Header.Peek("Content-Encoding")))) {
	case "gzip":
		if bs, err := r.response.BodyGunzip(); err == nil {
			return bs, nil
		}
	case "deflate":
		if bs, err := r.response.BodyInflate(); err == nil {
			return bs, nil
		}
	case "br":
		if bs, err := r.response.BodyUnbrotli(); err == nil {
			return bs, nil
		}
	}

	return r.response.Body(), nil
}
