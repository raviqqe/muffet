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
		{"-l", "4", "https://foo.com"},
		{"--limit-redirections", "4", "https://foo.com"},
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
