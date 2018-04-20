package main

import (
	"net/http"
)

// Fetcher represents a web page fetcher.
type Fetcher struct {
	connectionSemaphore semaphore
}

// newFetcher creates a new web page fetcher.
func newFetcher() Fetcher {
	return Fetcher{newSemaphore(512)}
}

// Fetch fetches a web page.
func (f Fetcher) Fetch(u string) (Page, error) {
	f.connectionSemaphore.Request()
	defer f.connectionSemaphore.Release()

	r, err := http.Get(u)

	if err != nil {
		return Page{}, err
	}

	return newPage(u, r.Body), nil
}
