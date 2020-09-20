package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRobotsTxtFetcherFetchRobotsTxt(t *testing.T) {
	s := "http://foo.com"
	u, err := url.Parse(s)
	assert.Nil(t, err)

	r, err := newRobotsTxtFetcher(
		newFakeHTTPClient(
			map[string]*fakeHTTPResponse{
				s + "/robots.txt": newFakeHTTPResponse(
					200,
					s,
					"text/plain",
					[]byte(`
						User-Agent: *
						Disallow: /bar
					`),
				),
			})).Fetch(u)

	assert.Nil(t, err)
	assert.False(t, r.TestAgent("/bar", "foo"))
}
