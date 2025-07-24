package main

import (
	"fmt"
	"mime"
	"net/url"
	"strings"
	"time"
)

type linkFetcher struct {
	client      httpClient
	pageParsers []pageParser
	cache       cache
	options     linkFetcherOptions
}

type fetchResult struct {
	StatusCode int
	Page       page
}

func newLinkFetcher(c httpClient, ps []pageParser, o linkFetcherOptions) *linkFetcher {
	return &linkFetcher{c, ps, newCache(), o}
}

// Fetch fetches a link and returns a successful status code and optionally HTML page, or an error.
func (f *linkFetcher) Fetch(u string) (int, page, error) {
	u, fr, err := separateFragment(u)
	if err != nil {
		return 0, nil, err
	}

	s, p, err := f.sendRequestWithCache(u)
	if err != nil {
		return 0, nil, err
	} else if p == nil || f.options.IgnoreFragments || fr == "" || strings.HasPrefix(fr, ":~:") {
		// TODO Support text fragments.
		return s, p, nil
	} else if _, ok := p.Fragments()[fr]; !ok {
		return 0, nil, fmt.Errorf("id #%v not found", fr)
	}

	return s, p, nil
}

func (f *linkFetcher) sendRequestWithCache(u string) (int, page, error) {
	x, store := f.cache.LoadOrStore(u)

	if store == nil {
		if err, ok := x.(error); ok {
			return 0, nil, err
		}

		r := x.(fetchResult)

		return r.StatusCode, r.Page, nil
	}

	attempts := 0
	retryDelay := 500 * time.Millisecond

	for {
		s, p, err := f.sendRequest(u)

		if err == nil {
			store(fetchResult{s, p})
		} else {
			if f.options.Retries > 0 && attempts < f.options.Retries-1 {
				attempts++
				time.Sleep(retryDelay)
				retryDelay *= 2 // Exponential backoff
				continue
			}
			store(err)
		}
		return s, p, err
	}
}

func (f *linkFetcher) sendRequest(s string) (int, page, error) {
	u, err := url.Parse(s)
	if err != nil {
		return 0, nil, err
	}

	r, err := f.client.Get(u, nil)

	if err != nil {
		return 0, nil, err
	}

	t := ""

	if s := strings.TrimSpace(r.Header("Content-Type")); s != "" {
		t, _, err = mime.ParseMediaType(s)

		if err != nil {
			return 0, nil, err
		}
	}

	bs, err := r.Body()
	if err != nil {
		return 0, nil, err
	}

	for _, pp := range f.pageParsers {
		u, err := url.Parse(r.URL())
		if err != nil {
			return 0, nil, err
		}

		p, err := pp.Parse(u, t, bs)
		if err != nil {
			return 0, nil, err
		} else if p != nil {
			return r.StatusCode(), p, nil
		}
	}

	return r.StatusCode(), nil, nil
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
