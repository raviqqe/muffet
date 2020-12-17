package main

import "strings"

type fakeHTTPResponse struct {
	statusCode  int
	location    string
	contentType string
	body        []byte
}

func newFakeHTTPResponse(statusCode int, location string, contentType string, body []byte) *fakeHTTPResponse {
	return &fakeHTTPResponse{statusCode, location, contentType, body}
}

func (r *fakeHTTPResponse) URL() string {
	return r.location
}

func (r *fakeHTTPResponse) StatusCode() int {
	return r.statusCode
}

func (r *fakeHTTPResponse) Header(name string) string {
	if strings.ToLower(name) == "content-type" {
		return r.contentType
	}

	return ""
}

func (r *fakeHTTPResponse) Body() ([]byte, error) {
	return r.body, nil
}
