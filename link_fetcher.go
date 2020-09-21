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

type fetchResult struct {
	StatusCode int
	Page       *page
}

func newLinkFetcher(c httpClient, pp *pageParser, o linkFetcherOptions) linkFetcher {
	o.Initialize()

	return linkFetcher{
		c,
		pp,
		newSemaphore(o.Concurrency),
		newCache(),
		o,
	}
}

// Fetch fetches a link and returns a successful status code and optionally HTML page, or an error.
func (f linkFetcher) Fetch(u string) (int, *page, error) {
	u, fr, err := separateFragment(u)
	if err != nil {
		return 0, nil, err
	}

	s, p, err := f.sendRequestWithCache(u)
	if err != nil {
		return 0, nil, err
	}

	if p != nil && !f.options.IgnoreFragments && fr != "" {
		if _, ok := p.IDs()[fr]; !ok {
			return 0, nil, fmt.Errorf("id #%v not found", fr)
		}
	}

	return s, p, nil
}

func (f linkFetcher) sendRequestWithCache(u string) (int, *page, error) {
	x, store := f.cache.LoadOrStore(u)

	if store == nil {
		if err, ok := x.(error); ok {
			return 0, nil, err
		}

		r := x.(fetchResult)

		return r.StatusCode, r.Page, nil
	}

	s, p, err := f.sendRequest(u)

	if err == nil {
		store(fetchResult{s, p})
	} else {
		store(err)
	}

	return s, p, err
}

func (f linkFetcher) sendRequest(s string) (int, *page, error) {
	f.connectionSemaphore.Request()
	defer f.connectionSemaphore.Release()

	u, err := url.Parse(s)
	if err != nil {
		return 0, nil, err
	}

	r, err := f.client.Get(u, f.options.Headers)

	if err != nil {
		return 0, nil, err
	} else if s := strings.TrimSpace(r.Header("Content-Type")); s != "" {
		t, _, err := mime.ParseMediaType(s)

		if err != nil {
			return 0, nil, err
		} else if t != "text/html" {
			return r.StatusCode(), nil, nil
		}
	}

	p, err := f.pageParser.Parse(r.URL(), r.Body())
	if err != nil {
		return 0, nil, err
	}

	return r.StatusCode(), p, nil
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