package main

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type fetcher struct {
	connectionSemaphore semaphore
}

func newFetcher(c int) fetcher {
	return fetcher{newSemaphore(c)}
}

func (f fetcher) Fetch(u string) (page, error) {
	f.connectionSemaphore.Request()
	defer f.connectionSemaphore.Release()

	s, b, err := fasthttp.Get(nil, u)

	if err != nil {
		return page{}, err
	}

	if s/100 != 2 {
		return page{}, fmt.Errorf("invalid status code: %v", s)
	}

	return newPage(u, b), nil
}
