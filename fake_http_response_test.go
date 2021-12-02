package main

import "strings"

type fakeHttpResponse struct {
	statusCode int
	location   string
	body       []byte
	headers    map[string]string
}

func newFakeHttpResponse(statusCode int, location string, body []byte, headers map[string]string) *fakeHttpResponse {
	hs := make(map[string]string, len(headers))

	for k, v := range headers {
		hs[strings.ToLower(k)] = v
	}

	return &fakeHttpResponse{statusCode, location, body, hs}
}

func newFakeHtmlResponse(location string, body string) *fakeHttpResponse {
	return newFakeHttpResponse(
		200,
		location,
		[]byte(body),
		map[string]string{"content-type": "text/html"},
	)
}

func (r *fakeHttpResponse) URL() string {
	return r.location
}

func (r *fakeHttpResponse) StatusCode() int {
	return r.statusCode
}

func (r *fakeHttpResponse) Header(name string) string {
	if v, ok := r.headers[strings.ToLower(name)]; ok {
		return v
	}

	return ""
}

func (r *fakeHttpResponse) Body() ([]byte, error) {
	return r.body, nil
}
