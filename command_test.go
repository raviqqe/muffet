package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func newTestCommand(h func(*url.URL) (*fakeHTTPResponse, error)) *command {
	return newTestCommandWithStdout(ioutil.Discard, h)
}

func newTestCommandWithStdout(stdout io.Writer, h func(*url.URL) (*fakeHTTPResponse, error)) *command {
	return newCommand(
		stdout,
		ioutil.Discard,
		false,
		newFakeHTTPClientFactory(h),
	)
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

func TestCommandRunWithVerboseOption(t *testing.T) {
	b := &bytes.Buffer{}

	ok := newTestCommandWithStdout(
		b,
		func(u *url.URL) (*fakeHTTPResponse, error) {
			return newFakeHTTPResponse(200, "http://foo.com", "text/html", nil), nil
		},
	).Run([]string{"-v", "http://foo.com"})

	assert.True(t, ok)
	assert.Greater(t, b.Len(), 0)
}

func TestCommandFailToRun(t *testing.T) {
	ok := newTestCommand(
		func(u *url.URL) (*fakeHTTPResponse, error) {
			return nil, errors.New("")
		},
	).Run([]string{"http://foo.com"})

	assert.False(t, ok)
}

func TestCommandFailToRunWithInvalidLink(t *testing.T) {
	b := &bytes.Buffer{}

	ok := newTestCommandWithStdout(
		b,
		func(u *url.URL) (*fakeHTTPResponse, error) {
			if u.String() == "http://foo.com" {
				return newFakeHTTPResponse(
					200,
					"http://foo.com",
					"text/html",
					[]byte(`<html><body><a href="/foo" /></body></html>`),
				), nil
			}

			return nil, errors.New("")
		},
	).Run([]string{"http://foo.com"})

	assert.False(t, ok)
	assert.Regexp(t, `http://foo\.com/foo`, b.String())
}

func TestCommandFailToParseArguments(t *testing.T) {
	b := &bytes.Buffer{}

	c := newTestCommandWithStderr(
		b,
		func(u *url.URL) (*fakeHTTPResponse, error) {
			return newFakeHTTPResponse(200, "", "text/html", nil), nil
		},
	)

	ok := c.Run(nil)

	assert.False(t, ok)
	assert.Greater(t, b.Len(), 0)
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
	cupaloy.SnapshotT(t, b.String())
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
	cupaloy.SnapshotT(t, b.String())
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
	cupaloy.SnapshotT(t, b.String())
}

func TestCommandShowHelp(t *testing.T) {
	b := &bytes.Buffer{}

	c := newTestCommandWithStdout(
		b,
		func(u *url.URL) (*fakeHTTPResponse, error) {
			return newFakeHTTPResponse(200, "", "text/html", nil), nil
		},
	)

	ok := c.Run([]string{"--help"})

	assert.True(t, ok)
	assert.Greater(t, b.Len(), 0)
}

func TestCommandShowVersion(t *testing.T) {
	b := &bytes.Buffer{}

	c := newTestCommandWithStdout(
		b,
		func(u *url.URL) (*fakeHTTPResponse, error) {
			return newFakeHTTPResponse(200, "", "text/html", nil), nil
		},
	)

	ok := c.Run([]string{"--version"})
	assert.True(t, ok)

	r, err := regexp.Compile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
	assert.Nil(t, err)
	assert.True(t, r.MatchString(strings.TrimSpace(b.String())))
}
