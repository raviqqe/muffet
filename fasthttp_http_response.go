package main

import (
	"bytes"
	"compress/flate"
	"io"
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
		return r.response.BodyGunzip()
	case "deflate":
		bs, err := r.response.BodyInflate()
		if err == nil {
			return bs, nil
		}

		zr := flate.NewReader(bytes.NewReader(r.response.Body()))
		defer zr.Close()

		raw, flateErr := io.ReadAll(zr)
		if flateErr == nil {
			return raw, nil
		}

		return nil, err
	case "br":
		return r.response.BodyUnbrotli()
	}

	return r.response.Body(), nil
}
