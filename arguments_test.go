package main

import (
	"testing"
)

func TestGetArguments(t *testing.T) {
	for _, ss := range [][]string{
		{"https://foo.com"},
		{"-c", "1", "https://foo.com"},
		{"--concurrency", "1", "https://foo.com"},
		{"-c", "foo", "https://foo.com"},
		{"-s", "https://foo.com"},
		{"--follow-sitemap-xml", "https://foo.com"},
		{"-t", "https://foo.com"},
		{"--skip-tls-verification", "https://foo.com"},
		{"-v", "https://foo.com"},
		{"--verbose", "https://foo.com"},
		{"-v", "-f", "https://foo.com"},
		{"-v", "--ignore-fragments", "https://foo.com"},
	} {
		getArguments(ss)
	}
}
