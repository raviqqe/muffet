package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"mime"
	"net/url"
	"strings"
	"sync"

	"github.com/valyala/fasthttp"
	"golang.org/x/net/html"
)

type fetcher struct {
	client              *fasthttp.Client
	connectionSemaphore semaphore
	cache               *sync.Map
	options             fetcherOptions
	scraper
}

func newFetcher(o fetcherOptions) fetcher {
	o.Initialize()

	return fetcher{
		&fasthttp.Client{
			MaxConnsPerHost: o.Concurrency,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: o.SkipTLSVerification,
			},
		},
		newSemaphore(o.Concurrency),
		&sync.Map{},
		o,
		newScraper(o.ExcludedPatterns),
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
	g := &sync.WaitGroup{}
	g.Add(1)
	defer g.Done()

	if x, ok := f.cache.LoadOrStore(u, g); ok {
		if g, ok := x.(*sync.WaitGroup); ok {
			g.Wait()
			x, _ = f.cache.Load(u)
		}

		switch x := x.(type) {
		case fetchResult:
			return x, nil
		default:
			return fetchResult{}, x.(error)
		}
	}

	r, err := f.sendRequest(u)

	if err == nil {
		f.cache.Store(u, r)
	} else {
		f.cache.Store(u, err)
	}

	return r, err
}

func (f fetcher) sendRequest(u string) (fetchResult, error) {
	f.connectionSemaphore.Request()
	defer f.connectionSemaphore.Release()

	req, res := fasthttp.Request{}, fasthttp.Response{}
	req.SetRequestURI(u)

	for k, v := range f.options.Headers {
		req.Header.Add(k, v)
	}

	r := 0

redirects:
	for {
		err := f.client.DoTimeout(&req, &res, f.options.Timeout)

		if err != nil {
			return fetchResult{}, err
		}

		switch res.StatusCode() / 100 {
		case 2:
			break redirects
		case 3:
			r++

			if r > f.options.MaxRedirections {
				return fetchResult{}, errors.New("too many redirections")
			}

			bs := res.Header.Peek("Location")

			if len(bs) == 0 {
				return fetchResult{}, errors.New("location header not found")
			}

			req.URI().UpdateBytes(bs)
		default:
			return fetchResult{}, fmt.Errorf("%v", res.StatusCode())
		}
	}

	if s := strings.TrimSpace(string(res.Header.Peek("Content-Type"))); s != "" {
		t, _, err := mime.ParseMediaType(s)

		if err != nil {
			return fetchResult{}, err
		} else if t != "text/html" {
			return newFetchResult(res.StatusCode()), nil
		}
	}

	n, err := html.Parse(bytes.NewReader(res.Body()))

	if err != nil {
		return fetchResult{}, err
	}

	p, err := newPage(req.URI().String(), n, f.scraper)

	if err != nil {
		return fetchResult{}, err
	}

	return newFetchResultWithPage(res.StatusCode(), p), nil
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
