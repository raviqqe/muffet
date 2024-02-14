package main

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func TestGetArguments(t *testing.T) {
	for _, ss := range [][]string{
		{"https://foo.com"},
		{"-b", "42", "https://foo.com"},
		{"--buffer-size", "42", "https://foo.com"},
		{"-c", "1", "https://foo.com"},
		{"--max-connections", "1", "https://foo.com"},
		{"--max-connections-per-host", "1", "https://foo.com"},
		{"--max-response-body-size", "1", "https://foo.com"},
		{"-e", "regex1", "-e", "regex2", "https://foo.com"},
		{"--exclude", "regex1", "--exclude", "regex2", "https://foo.com"},
		{"--header", "MyHeader: foo", "--header", "YourHeader: bar", "https://foo.com"},
		{"--header", "User-Agent: custom-agent", "https://foo.com"},
		{"-r", "4", "https://foo.com"},
		{"--max-redirections", "4", "https://foo.com"},
		{"--follow-robots-txt", "https://foo.com"},
		{"--follow-sitemap-xml", "https://foo.com"},
		{"-t", "10", "https://foo.com"},
		{"--timeout", "10", "https://foo.com"},
		{"--rate-limit", "1", "https://foo.com"},
		{"--proxy", "localhost:8080", "https://foo.com"},
		{"--skip-tls-verification", "https://foo.com"},
		{"-v", "https://foo.com"},
		{"--verbose", "https://foo.com"},
		{"-v", "-f", "https://foo.com"},
		{"-v", "--ignore-fragments", "https://foo.com"},
		{"--one-page-only", "https://foo.com"},
		{"--json", "https://foo.com"},
		{"-h"},
		{"--help"},
		{"--version"},
	} {
		_, err := getArguments(ss, nil)
		assert.Nil(t, err)
	}
}

func TestGetArgumentsError(t *testing.T) {
	for _, ss := range [][]string{
		{},
		{"-b", "foo", "https://foo.com"},
		{"--buffer-size", "foo", "https://foo.com"},
		{"-c", "foo", "https://foo.com"},
		{"--max-connections", "foo", "https://foo.com"},
		{"-e", "(", "https://foo.com"},
		{"-j", "MyHeader", "https://foo.com"},
		{"--header", "MyHeader", "https://foo.com"},
		{"-l", "foo", "https://foo.com"},
		{"--max-redirections", "foo", "https://foo.com"},
		{"-t", "foo", "https://foo.com"},
		{"--timeout", "foo", "https://foo.com"},
	} {
		_, err := getArguments(ss, nil)
		assert.NotNil(t, err)
	}
}

func TestGetArgumentsWithIniFile(t *testing.T) {
	ini := `
bufferSize = 8192
exclude = foo.com 
exclude = bar.com
maxConnectionsPerHost = 122
`
	args := []string{"--header", "a:fo", "--max-connections-per-host", "123", "https://baz.com"}
	outArgs, err := getArguments(args, bytes.NewBufferString(ini))
	assert.Nil(t, err)

	// Just from the ini file (the global default is overriden)
	assert.Equal(t, 8192, outArgs.BufferSize)
	// Not set anywhere (the global default is taken)
	assert.Equal(t, 512, outArgs.MaxConnections)
	// Command line takes priority over ini file
	assert.Equal(t, 123, outArgs.MaxConnectionsPerHost)
	// Just on command line
	assert.Equal(t, []string{"a:fo"}, outArgs.RawHeaders)
	// Repeated entry in ini file lead to multiple items
	assert.Equal(t, []string{"foo.com", "bar.com"}, outArgs.RawExcludedPatterns)
}

func TestHelp(t *testing.T) {
	cupaloy.SnapshotT(t, help())
}

func TestParseHeader(t *testing.T) {
	for _, c := range []struct {
		arguments []string
		header    http.Header
	}{
		{
			nil,
			http.Header{},
		},
		{
			[]string{"MyHeader: foo"},
			(func() http.Header {
				h := http.Header{}
				h.Add("MyHeader", "foo")
				return h
			})(),
		},
		{
			[]string{"MyHeader: foo", "YourHeader: bar"},
			(func() http.Header {
				h := http.Header{}
				h.Add("MyHeader", "foo")
				h.Add("YourHeader", "bar")
				return h
			})(),
		},
	} {
		hs, err := parseHeaders(c.arguments)

		assert.Nil(t, err)
		assert.Equal(t, c.header, hs)
	}
}

func TestParseHeadersError(t *testing.T) {
	_, err := parseHeaders([]string{"MyHeader"})
	assert.NotNil(t, err)
}
