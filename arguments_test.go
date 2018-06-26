package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetArguments(t *testing.T) {
	for _, ss := range [][]string{
		{"https://foo.com"},
		{"-c", "1", "https://foo.com"},
		{"--concurrency", "1", "https://foo.com"},
		{"-e", "regex1", "-e", "regex2", "https://foo.com"},
		{"--exclude", "regex1", "--exclude", "regex2", "https://foo.com"},
		{"-j", "MyHeader: foo", "-j", "YourHeader: bar", "https://foo.com"},
		{"--header", "MyHeader: foo", "--header", "YourHeader: bar", "https://foo.com"},
		{"-l", "4", "https://foo.com"},
		{"--limit-redirections", "4", "https://foo.com"},
		{"-n", "https://foo.com"},
		{"--remove-newlines", "https://foo.com"},
		{"-s", "https://foo.com"},
		{"--follow-sitemap-xml", "https://foo.com"},
		{"-t", "10", "https://foo.com"},
		{"--timeout", "10", "https://foo.com"},
		{"-x", "https://foo.com"},
		{"--skip-tls-verification", "https://foo.com"},
		{"-v", "https://foo.com"},
		{"--verbose", "https://foo.com"},
		{"-v", "-f", "https://foo.com"},
		{"-v", "--ignore-fragments", "https://foo.com"},
	} {
		_, err := getArguments(ss)
		assert.Nil(t, err)
	}
}

func TestGetArgumentsError(t *testing.T) {
	for _, ss := range [][]string{
		{"-c", "foo", "https://foo.com"},
		{"--concurrency", "foo", "https://foo.com"},
		{"-e", "(", "https://foo.com"},
		{"-j", "MyHeader", "https://foo.com"},
		{"--header", "MyHeader", "https://foo.com"},
		{"-l", "foo", "https://foo.com"},
		{"--limit-redirections", "foo", "https://foo.com"},
		{"-t", "foo", "https://foo.com"},
		{"--timeout", "foo", "https://foo.com"},
	} {
		_, err := getArguments(ss)
		assert.NotNil(t, err)
	}
}

func TestParseArguments(t *testing.T) {
	assert.Panics(t, func() {
		parseArguments("", nil)
	})
}

func TestParseHeaders(t *testing.T) {
	for _, c := range []struct {
		arguments []string
		answer    map[string]string
	}{
		{
			nil,
			map[string]string{},
		},
		{
			[]string{"MyHeader: foo"},
			map[string]string{"MyHeader": "foo"},
		},
		{
			[]string{"MyHeader: foo", "YourHeader: bar"},
			map[string]string{"MyHeader": "foo", "YourHeader": "bar"},
		},
	} {
		hs, err := parseHeaders(c.arguments)

		assert.Nil(t, err)
		assert.Equal(t, c.answer, hs)
	}
}

func TestParseHeadersError(t *testing.T) {
	_, err := parseHeaders([]string{"MyHeader"})
	assert.NotNil(t, err)
}
