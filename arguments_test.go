package main

import (
	"net/http"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func TestGetArguments(t *testing.T) {
	for _, ss := range [][]string{
		{"https://foo.com"},
		{"--accepted-status-codes", "200..300,403", "https://foo.com"},
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
		{"--format=json", "https://foo.com"},
		{"-h"},
		{"--help"},
		{"--version"},
	} {
		_, err := getArguments(ss)
		assert.Nil(t, err)
	}
}

func TestGetArgumentsError(t *testing.T) {
	for _, ss := range [][]string{
		{},
		{"--accepted-status-codes", "foo", "https://foo.com"},
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
		_, err := getArguments(ss)
		assert.NotNil(t, err)
	}
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
