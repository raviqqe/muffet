package main

import (
	"fmt"
	"sync"

	"github.com/valyala/fasthttp"
)

type fetcher struct {
	connectionSemaphore semaphore
	cache               *sync.Map
}

func newFetcher(c int) fetcher {
	return fetcher{newSemaphore(c), &sync.Map{}}
}

func (f fetcher) Fetch(u string, h func(page)) error {
	if err, ok := f.cache.Load(u); ok && err == nil {
		return nil
	} else if ok {
		return err.(error)
	}

	f.connectionSemaphore.Request()
	defer f.connectionSemaphore.Release()

	s, b, err := fasthttp.Get(nil, u)
	f.cache.Store(u, err)

	if err != nil {
		return err
	}

	if s/100 != 2 {
		return fmt.Errorf("invalid status code %v", s)
	}

	if h != nil {
		h(newPage(u, b))
	}

	return nil
}
