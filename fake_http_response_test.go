package main

import "strings"

type fakeHttpResponse struct {
	statusCode  int
	location    string
	contentType string
	body        []byte
}

func newFakeHttpResponse(statusCode int, location string, contentType string, body []byte) *fakeHttpResponse {
	return &fakeHttpResponse{statusCode, location, contentType, body}
}

func (r *fakeHttpResponse) URL() string {
	return r.location
}

func (r *fakeHttpResponse) StatusCode() int {
	return r.statusCode
}

func (r *fakeHttpResponse) Header(name string) string {
	if strings.ToLower(name) == "content-type" {
		return r.contentType
	}

	return ""
}

func (r *fakeHttpResponse) Body() ([]byte, error) {
	return r.body, nil
}
