package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckedHttpClientFailWithValidStatusCode(t *testing.T) {
	u, err := url.Parse(testUrl)

	assert.Nil(t, err)

	r, err := newCheckedHttpClient(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				return newFakeHttpResponse(200, testUrl, nil, nil), nil
			},
		),
		statusCodeSet{{200, 201}: {}},
	).Get(u, nil)

	assert.Nil(t, err)
	assert.NotNil(t, r)
}

func TestCheckedHttpClientFailWithInvalidStatusCode(t *testing.T) {
	u, err := url.Parse(testUrl)

	assert.Nil(t, err)

	r, err := newCheckedHttpClient(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				return newFakeHttpResponse(404, testUrl, nil, nil), nil
			},
		),
		statusCodeSet{{200, 201}: {}},
	).Get(u, nil)

	assert.Nil(t, r)
	assert.Equal(t, err.Error(), "404")
}
