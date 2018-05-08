package main

import (
	"bytes"
	"crypto/tls"
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

func newFetcher(c int, f, t bool) fetcher {
	return fetcher{
		&fasthttp.Client{
			MaxConnsPerHost: c,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: t,
			},
		},
		newSemaphore(c),
		&sync.Map{},
		f,
	}
}

func (f fetcher) FetchPage(s string) (page, error) {
	_, p, err := f.sendRequest(s)

	return p, err
}

func (f fetcher) FetchLink(s string) (linkResult, error) {
	s, fr, err := separateFragment(s)

	if err != nil {
		return linkResult{}, err
	}

	if x, ok := f.cache.Load(s); ok {
		switch x := x.(type) {
		case int:
			return newLinkResult(x), nil
		case error:
			return linkResult{}, x
		}
	}

	c, p, err := f.sendRequestWithFragment(s, fr)

	if err == nil {
		f.cache.Store(s, c)
	} else {
		f.cache.Store(s, err)
	}

	return newLinkResultWithPage(c, p), err
}

func (f fetcher) sendRequestWithFragment(u, fr string) (int, page, error) {
	c, p, err := f.sendRequest(u)

	if err != nil {
		return 0, page{}, err
	}

	if !f.ignoreFragments && fr != "" {
		if _, ok := scrape.Find(p.Body(), func(n *html.Node) bool {
			return scrape.Attr(n, "id") == fr
		}); !ok {
			return 0, page{}, fmt.Errorf("id #%v not found", fr)
		}
	}

	return c, p, nil
}

func (f fetcher) sendRequest(u string) (int, page, error) {
	f.connectionSemaphore.Request()
	defer f.connectionSemaphore.Release()

	req, res := fasthttp.Request{}, fasthttp.Response{}
	req.SetRequestURI(u)

redirects:
	for {
		err := f.client.Do(&req, &res)

		if err != nil {
			return 0, page{}, err
		}

		switch res.StatusCode() / 100 {
		case 2:
			break redirects
		case 3:
			bs := res.Header.Peek("Location")

			if len(bs) == 0 {
				return 0, page{}, errors.New("location header not found")
			}

			req.SetRequestURIBytes(bs)
		default:
			return 0, page{}, fmt.Errorf("%v", res.StatusCode())
		}
	}

	n, err := html.Parse(bytes.NewReader(res.Body()))

	if err != nil {
		return 0, page{}, err
	}

	return res.StatusCode(), newPage(req.URI().String(), n), nil
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
