package main

import "net/url"

type pageParser interface {
	// Returned pages can be nil to indicate unrecognized file formats.
	Parse(*url.URL, string, []byte) (page, error)
}
