package main

import (
	"io"
)

// Page represents a web page fetched already.
type Page struct {
	url  string
	body io.Reader
}

// newPage creates a new web page.
func newPage(u string, b io.Reader) Page {
	return Page{u, b}
}

// URL returns a URL of a fetched page.
func (p Page) URL() string {
	return p.url
}

// Body returns a body reader of a fetched page.
func (p Page) Body() io.Reader {
	return p.body
}
