package main

import (
	"fmt"
	"sync"

	"github.com/valyala/fasthttp"
)

type fetcher struct {
	client              *fasthttp.Client
	connectionSemaphore semaphore
	cache               *sync.Map
}

func newFetcher(c int) fetcher {
	return fetcher{
		&fasthttp.Client{MaxConnsPerHost: c},
		newSemaphore(c),
		&sync.Map{},
	}
}

func (f fetcher) Fetch(u string) (*page, error) {
	if err, ok := f.cache.Load(u); ok && err == nil {
		return nil, nil
	} else if ok {
		return nil, err.(error)
	}

	f.connectionSemaphore.Request()
	defer f.connectionSemaphore.Release()

	s, b, err := f.client.Get(nil, u)

	if err == nil && s/100 != 2 {
		err = fmt.Errorf("invalid status code %v", s)
	}

	f.cache.Store(u, err)

	if err != nil {
		return nil, err
	}

	p := newPage(u, b)
	return &p, nil
}
