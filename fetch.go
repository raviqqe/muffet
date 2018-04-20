package main

import (
	"io"
	"net/http"
)

var connectionSemaphore = newSemaphore(512)

// Page represents a web page fetched already.
type Page struct {
	url  string
	body io.Reader
}

// URL returns a URL of a fetched page.
func (p Page) URL() string {
	return p.url
}

// Body returns a body reader of a fetched page.
func (p Page) Body() io.Reader {
	return p.body
}

func fetch(u string) (Page, error) {
	connectionSemaphore.Request()
	defer connectionSemaphore.Release()

	r, err := http.Get(u)

	if err != nil {
		return Page{}, err
	}

	return Page{u, r.Body}, nil
}
