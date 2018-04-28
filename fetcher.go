package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/valyala/fasthttp"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
)

type fetcher struct {
	client              *fasthttp.Client
	connectionSemaphore semaphore
	cache               *sync.Map
	ignoreFragments     bool
}

func newFetcher(c int, f bool) fetcher {
	return fetcher{
		&fasthttp.Client{MaxConnsPerHost: c},
		newSemaphore(c),
		&sync.Map{},
		f,
	}
}

func (f fetcher) Fetch(s string) (fetchResult, error) {
	s, fr, err := separateFragment(s)

	if err != nil {
		return fetchResult{}, err
	}

	if x, ok := f.cache.Load(s); ok {
		switch x := x.(type) {
		case fetchResult:
			return x, nil
		case error:
			return fetchResult{}, x
		}
	}

	r, err := f.sendRequestWithFragment(s, fr)

	if err == nil {
		f.cache.Store(s, newFetchResult(r.StatusCode()))
	} else {
		f.cache.Store(s, err)
	}

	return r, err
}

func (f fetcher) sendRequestWithFragment(u, fr string) (fetchResult, error) {
	r, err := f.sendRequest(u)

	if err != nil {
		return fetchResult{}, err
	}

	if p, ok := r.Page(); ok && !f.ignoreFragments && fr != "" {
		if _, ok := scrape.Find(p.Body(), func(n *html.Node) bool {
			return scrape.Attr(n, "id") == fr
		}); !ok {
			return fetchResult{}, fmt.Errorf("id #%v not found", fr)
		}
	}

	return r, nil
}

func (f fetcher) sendRequest(u string) (fetchResult, error) {
	f.connectionSemaphore.Request()
	defer f.connectionSemaphore.Release()

	req, res := fasthttp.Request{}, fasthttp.Response{}
	req.SetRequestURI(u)

redirects:
	for {
		err := f.client.Do(&req, &res)

		if err != nil {
			return fetchResult{}, err
		}

		switch res.StatusCode() / 100 {
		case 2:
			break redirects
		case 3:
			bs := res.Header.Peek("Location")

			if len(bs) == 0 {
				return fetchResult{}, errors.New("location header not found")
			}

			req.SetRequestURIBytes(bs)
		default:
			return fetchResult{}, fmt.Errorf("%v", res.StatusCode())
		}
	}

	n, err := html.Parse(bytes.NewReader(res.Body()))

	if err != nil {
		return fetchResult{}, err
	}

	return newFetchResultWithPage(res.StatusCode(), newPage(req.URI().String(), n)), nil
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
