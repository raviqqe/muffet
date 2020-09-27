package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func newTestCommand(h func(*url.URL) (*fakeHTTPResponse, error)) *command {
	return newTestCommandWithStderr(ioutil.Discard, h)
}

func newTestCommandWithStderr(stderr io.Writer, h func(*url.URL) (*fakeHTTPResponse, error)) *command {
	return newCommand(
		ioutil.Discard,
		stderr,
		false,
		newFakeHTTPClientFactory(h),
	)
}

func TestCommandRun(t *testing.T) {
	ok := newTestCommand(
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

func TestCommandRunWithLinks(t *testing.T) {
	visited := false

	ok := newTestCommand(
		func(u *url.URL) (*fakeHTTPResponse, error) {
			switch u.String() {
			case "http://foo.com":
				return newFakeHTTPResponse(
					200,
					"http://foo.com",
					"text/html",
					[]byte(`<html><body><a href="/foo" /></body></html>`),
				), nil
			case "http://foo.com/foo":
				visited = true
				return newFakeHTTPResponse(200, "http://foo.com", "text/html", nil), nil
			}

			return nil, errors.New("")
		},
	).Run([]string{"http://foo.com"})

	assert.True(t, ok)
	assert.True(t, visited)
}

func TestCommandFailToRun(t *testing.T) {
	ok := newTestCommand(
		func(u *url.URL) (*fakeHTTPResponse, error) {
			return nil, errors.New("")
		},
	).Run([]string{"http://foo.com"})

	assert.False(t, ok)
}

func TestCommandFailToFetchRootPage(t *testing.T) {
	b := &bytes.Buffer{}

	ok := newTestCommandWithStderr(
		b,
		func(u *url.URL) (*fakeHTTPResponse, error) {
			return nil, errors.New("foo")
		},
	).Run([]string{"http://foo.com"})

	assert.False(t, ok)
	cupaloy.SnapshotT(t, b.Bytes())
}

func TestCommandFailToGetHTMLRootPage(t *testing.T) {
	b := &bytes.Buffer{}

	ok := newTestCommandWithStderr(
		b,
		func(u *url.URL) (*fakeHTTPResponse, error) {
			return newFakeHTTPResponse(200, "", "image/png", nil), nil
		},
	).Run([]string{"http://foo.com"})

	assert.False(t, ok)
	cupaloy.SnapshotT(t, b.Bytes())
}

func TestCommandColorErrorMessage(t *testing.T) {
	b := &bytes.Buffer{}

	c := newCommand(
		ioutil.Discard,
		b,
		true,
		newFakeHTTPClientFactory(func(u *url.URL) (*fakeHTTPResponse, error) {
			return nil, errors.New("foo")
		}),
	)

	ok := c.Run([]string{"http://foo.com"})

	assert.False(t, ok)
	cupaloy.SnapshotT(t, b.Bytes())
}
