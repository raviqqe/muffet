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

func (r *fakeHTTPResponse) StatusCode() int {
	return r.statusCode
}

func (r *fakeHTTPResponse) Header(name string) string {
	switch strings.ToLower(name) {
	case "location":
		return r.location
	case "content-type":
		return r.location
	}

	return ""
}

func (r *fakeHTTPResponse) Body() []byte {
	return r.body
}
