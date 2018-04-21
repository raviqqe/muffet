package main

import "github.com/valyala/fasthttp"

type fetcher struct {
	connectionSemaphore semaphore
}

func newFetcher(c int) fetcher {
	return fetcher{newSemaphore(c)}
}

func (f fetcher) Fetch(u string) (page, error) {
	f.connectionSemaphore.Request()
	defer f.connectionSemaphore.Release()

	_, b, err := fasthttp.Get(nil, u)

	if err != nil {
		return page{}, err
	}

	return newPage(u, b), nil
}
