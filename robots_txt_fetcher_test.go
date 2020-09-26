package main

import (
	"errors"
	"net/url"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func TestRobotsTxtFetcherFetchRobotsTxt(t *testing.T) {
	s := "http://foo.com"
	u, err := url.Parse(s)
	assert.Nil(t, err)

	r, err := newRobotsTxtFetcher(
		newFakeHTTPClient(
			func(u *url.URL) (*fakeHTTPResponse, error) {
				if u.String() != s+"/robots.txt" {
					return nil, errors.New("")
				}

				return newFakeHTTPResponse(
					200,
					s,
					"text/plain",
					[]byte(`
						User-Agent: *
						Disallow: /bar
					`),
				), nil
			})).Fetch(u)

	assert.Nil(t, err)
	assert.False(t, r.TestAgent("/bar", "foo"))
}

func TestRobotsTxtFetcherFailToFetchRobotsTxt(t *testing.T) {
	u, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	_, err = newRobotsTxtFetcher(
		newFakeHTTPClient(
			func(u *url.URL) (*fakeHTTPResponse, error) {
				return nil, errors.New("foo")
			})).Fetch(u)

	cupaloy.SnapshotT(t, err.Error())
}
