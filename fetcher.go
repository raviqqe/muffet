package main

import "github.com/valyala/fasthttp"

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

	_, b, err := fasthttp.Get(nil, u)

	if err != nil {
		return Page{}, err
	}

	return newPage(u, b), nil
}
