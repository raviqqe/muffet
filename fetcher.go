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

func (f fetcher) Fetch(s string) (*page, error) {
	s, fr, err := separateFragment(s)

	if err != nil {
		return nil, err
	}

	if err, ok := f.cache.Load(s); ok && err == nil {
		return nil, nil
	} else if ok {
		return nil, err.(error)
	}

	p, err := f.fetchPageWithFragment(s, fr)

	f.cache.Store(s, err)

	return p, err
}

func (f fetcher) fetchPageWithFragment(u, fr string) (*page, error) {
	p, err := f.fetchPage(u)

	if err != nil {
		return nil, err
	}

	if !f.ignoreFragments && fr != "" {
		if _, ok := scrape.Find(p.Body(), func(n *html.Node) bool {
			return scrape.Attr(n, "id") == fr
		}); !ok {
			return nil, fmt.Errorf("id #%v not found", fr)
		}
	}

	return p, nil
}

func (f fetcher) fetchPage(u string) (*page, error) {
	f.connectionSemaphore.Request()
	defer f.connectionSemaphore.Release()

	req, res := fasthttp.Request{}, fasthttp.Response{}
	req.SetRequestURI(u)

redirects:
	for {
		err := f.client.Do(&req, &res)

		if err != nil {
			return nil, err
		}

		switch res.StatusCode() / 100 {
		case 2:
			break redirects
		case 3:
			bs := res.Header.Peek("Location")

			if len(bs) == 0 {
				return nil, errors.New("location header not found")
			}

			req.SetRequestURIBytes(bs)
		default:
			return nil, fmt.Errorf("invalid status code %v", res.StatusCode())
		}
	}

	n, err := html.Parse(bytes.NewReader(res.Body()))

	if err != nil {
		return nil, err
	}

	p := newPage(req.URI().String(), n)
	return &p, nil
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
