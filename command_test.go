package main

import (
	"bytes"
	"errors"
	"io"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func newTestCommand(h func(*url.URL) (*fakeHttpResponse, error)) *command {
	return newTestCommandWithStdout(io.Discard, h)
}

func newTestCommandWithStdout(stdout io.Writer, h func(*url.URL) (*fakeHttpResponse, error)) *command {
	return newCommand(
		stdout,
		io.Discard,
		false,
		newFakeHttpClientFactory(h),
	)
}

func newTestCommandWithStderr(stderr io.Writer, h func(*url.URL) (*fakeHttpResponse, error)) *command {
	return newCommand(
		io.Discard,
		stderr,
		false,
		newFakeHttpClientFactory(h),
	)
}

func TestCommandRun(t *testing.T) {
	ok := newTestCommand(
		func(u *url.URL) (*fakeHttpResponse, error) {
			s := "http://foo.com"

			if u.String() != s {
				return nil, errors.New("")
			}

			return newFakeHtmlResponse(s, ""), nil
		},
	).Run([]string{"http://foo.com"})

	assert.True(t, ok)
}

func TestCommandRunWithLinks(t *testing.T) {
	visited := false

	ok := newTestCommand(
		func(u *url.URL) (*fakeHttpResponse, error) {
			switch u.String() {
			case "http://foo.com":
				return newFakeHtmlResponse(
					"http://foo.com",
					`<html><body><a href="/foo" /></body></html>`,
				), nil
			case "http://foo.com/foo":
				visited = true
				return newFakeHtmlResponse("http://foo.com", ""), nil
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
		func(u *url.URL) (*fakeHttpResponse, error) {
			return newFakeHtmlResponse("http://foo.com", ""), nil
		},
	).Run([]string{"-v", "http://foo.com"})

	assert.True(t, ok)
	assert.Greater(t, b.Len(), 0)
}

func TestCommandFailToRun(t *testing.T) {
	ok := newTestCommand(
		func(u *url.URL) (*fakeHttpResponse, error) {
			return nil, errors.New("")
		},
	).Run([]string{"http://foo.com"})

	assert.False(t, ok)
}

func TestCommandFailToRunWithInvalidLink(t *testing.T) {
	b := &bytes.Buffer{}

	ok := newTestCommandWithStdout(
		b,
		func(u *url.URL) (*fakeHttpResponse, error) {
			if u.String() == "http://foo.com" {
				return newFakeHtmlResponse(
					"http://foo.com",
					`<html><body><a href="/foo" /></body></html>`,
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
		func(u *url.URL) (*fakeHttpResponse, error) {
			return newFakeHtmlResponse("", ""), nil
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
		func(u *url.URL) (*fakeHttpResponse, error) {
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
		func(u *url.URL) (*fakeHttpResponse, error) {
			return newFakeHttpResponse(
				200,
				"",
				nil,
				map[string]string{"content-type": "image/png"},
			), nil
		},
	).Run([]string{"http://foo.com"})

	assert.False(t, ok)
	cupaloy.SnapshotT(t, b.String())
}

func TestCommandColorErrorMessage(t *testing.T) {
	b := &bytes.Buffer{}

	c := newCommand(
		io.Discard,
		b,
		true,
		newFakeHttpClientFactory(func(u *url.URL) (*fakeHttpResponse, error) {
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
		func(u *url.URL) (*fakeHttpResponse, error) {
			return newFakeHtmlResponse("", ""), nil
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
		func(u *url.URL) (*fakeHttpResponse, error) {
			return newFakeHtmlResponse("", ""), nil
		},
	)

	ok := c.Run([]string{"--version"})
	assert.True(t, ok)

	r, err := regexp.Compile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
	assert.Nil(t, err)
	assert.True(t, r.MatchString(strings.TrimSpace(b.String())))
}

func TestCommandFailToRunWithJSONOutput(t *testing.T) {
	b := &bytes.Buffer{}

	ok := newTestCommandWithStdout(
		b,
		func(u *url.URL) (*fakeHttpResponse, error) {
			if u.String() == "http://foo.com" {
				return newFakeHtmlResponse(
					"http://foo.com",
					`<html><body><a href="/foo" /></body></html>`,
				), nil
			}

			return nil, errors.New("foo")
		},
	).Run([]string{"--json", "http://foo.com"})

	assert.False(t, ok)
	assert.Greater(t, b.Len(), 0)
}

func TestCommandDoNotIncludeSuccessfulPageInJSONOutput(t *testing.T) {
	b := &bytes.Buffer{}

	ok := newTestCommandWithStdout(
		b,
		func(u *url.URL) (*fakeHttpResponse, error) {
			return newFakeHtmlResponse("", ""), nil
		},
	).Run([]string{"--json", "http://foo.com"})

	assert.True(t, ok)
	assert.Equal(t, strings.TrimSpace(b.String()), "[]")
}

func TestCommandIncludeSuccessfulPageInJSONOutputWhenRequested(t *testing.T) {
	b := &bytes.Buffer{}

	ok := newTestCommandWithStdout(
		b,
		func(u *url.URL) (*fakeHttpResponse, error) {
			return newFakeHtmlResponse("", ""), nil
		},
	).Run([]string{"--json", "--include-success-in-json", "http://foo.com"})

	assert.True(t, ok)
	assert.Equal(t, strings.TrimSpace(b.String()), "[{\"url\":\"\",\"links\":[]}]")
}

func TestCommandFailWithVerboseAndJSONOptions(t *testing.T) {
	b := &bytes.Buffer{}

	ok := newTestCommandWithStderr(b, nil).Run(
		[]string{"--json", "--verbose", "http://foo.com"},
	)

	assert.False(t, ok)
	cupaloy.SnapshotT(t, b.String())
}
