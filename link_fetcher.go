package main

import (
	"fmt"
	"mime"
	"net/url"
	"strings"
)

type linkFetcher struct {
	client              httpClient
	pageParser          *pageParser
	connectionSemaphore semaphore
	cache               cache
	options             linkFetcherOptions
}

func newLinkFetcher(c httpClient, pp *pageParser, o linkFetcherOptions) linkFetcher {
	return linkFetcher{
		c,
		pp,
		newSemaphore(o.Concurrency),
		newCache(),
		o,
	}
}

func (f linkFetcher) Fetch(u string) (fetchResult, error) {
	u, fr, err := separateFragment(u)
	if err != nil {
		return fetchResult{}, err
	}

	r, err := f.sendRequestWithCache(u)
	if err != nil {
		return fetchResult{}, err
	}

	if p, ok := r.Page(); ok && !f.options.IgnoreFragments && fr != "" {
		if _, ok := p.IDs()[fr]; !ok {
			return fetchResult{}, fmt.Errorf("id #%v not found", fr)
		}
	}

	return r, nil
}

func (f linkFetcher) sendRequestWithCache(u string) (fetchResult, error) {
	x, store := f.cache.LoadOrStore(u)

	if store == nil {
		if err, ok := x.(error); ok {
			return fetchResult{}, err
		}

		return x.(fetchResult), nil
	}

	r, err := f.sendRequest(u)

	if err == nil {
		store(r)
	} else {
		store(err)
	}

	return r, err
}

func (f linkFetcher) sendRequest(s string) (fetchResult, error) {
	f.connectionSemaphore.Request()
	defer f.connectionSemaphore.Release()

	u, err := url.Parse(s)
	if err != nil {
		return fetchResult{}, err
	}

	r, err := f.client.Get(u, f.options.Headers)

	if err != nil {
		return fetchResult{}, err
	} else if s := strings.TrimSpace(r.Header("Content-Type")); s != "" {
		t, _, err := mime.ParseMediaType(s)

		if err != nil {
			return fetchResult{}, err
		} else if t != "text/html" {
			return newFetchResult(r.StatusCode(), nil), nil
		}
	}

	p, err := f.pageParser.Parse(r.URL(), r.Body())
	if err != nil {
		return fetchResult{}, err
	}

	return newFetchResult(r.StatusCode(), p), nil
}

func separateFragment(s string) (string, string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", "", err
	}

	f := u.Fragment
	u.Fragment = ""

	return u.String(), f, nil
}
