package main

import (
	"bytes"
	"errors"
	"fmt"
	"mime"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type fetcher struct {
	client              httpClient
	connectionSemaphore semaphore
	cache               cache
	options             fetcherOptions
	scraper
}

func newFetcher(c httpClient, o fetcherOptions) fetcher {
	o.Initialize()

	return fetcher{
		c,
		newSemaphore(o.Concurrency),
		newCache(),
		o,
		newScraper(o.ExcludedPatterns, o.FollowURLParams),
	}
}

func (f fetcher) Fetch(u string) (fetchResult, error) {
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

func (f fetcher) sendRequestWithCache(u string) (fetchResult, error) {
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

func (f fetcher) sendRequest(rawURL string) (fetchResult, error) {
	f.connectionSemaphore.Request()
	defer f.connectionSemaphore.Release()

	i := 0

	for {
		u, err := url.Parse(rawURL)
		if err != nil {
			return fetchResult{}, err
		}

		r, err := f.client.Get(u, f.options.Headers, f.options.Timeout)
		if err != nil {
			return fetchResult{}, err
		}

		switch r.StatusCode() / 100 {
		case 2:
			return f.createSuccessfulResult(r, rawURL)
		case 3:
			i++

			if i > f.options.MaxRedirections {
				return fetchResult{}, errors.New("too many redirections")
			}

			rawURL = r.Header("Location")

			if len(rawURL) == 0 {
				return fetchResult{}, errors.New("location header not found")
			}
		default:
			return fetchResult{}, fmt.Errorf("%v", r.StatusCode())
		}
	}
}

func (f fetcher) createSuccessfulResult(r httpResponse, u string) (fetchResult, error) {
	if s := strings.TrimSpace(r.Header("Content-Type")); s != "" {
		t, _, err := mime.ParseMediaType(s)

		if err != nil {
			return fetchResult{}, err
		} else if t != "text/html" {
			return newFetchResult(r.StatusCode(), nil), nil
		}
	}

	n, err := html.Parse(bytes.NewReader(r.Body()))
	if err != nil {
		return fetchResult{}, err
	}

	p, err := newPage(u, n, f.scraper)
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
