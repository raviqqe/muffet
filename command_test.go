package main

import (
	"errors"
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestCommand(h func(*url.URL) (*fakeHTTPResponse, error)) *command {
	return newCommand(
		ioutil.Discard,
		ioutil.Discard,
		newFakeHTTPClientFactory(h),
	)
}

func TestCommand(t *testing.T) {
	ok := createTestCommand(
		func(u *url.URL) (*fakeHTTPResponse, error) {
			s := "http://foo.com"

			if u.String() != s {
				return nil, errors.New("")
			}

			return newFakeHTTPResponse(200, s, "text/html", nil), nil
		},
	).Run([]string{"http://foo.com"})

	assert.True(t, ok)
}

func TestCommandFailToRun(t *testing.T) {
	ok := createTestCommand(
		func(u *url.URL) (*fakeHTTPResponse, error) {
			return nil, errors.New("")
		},
	).Run([]string{"http://foo.com"})

	assert.False(t, ok)
}
