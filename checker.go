package main

import (
	"bytes"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var validSchemes = map[string]struct{}{
	"":      struct{}{},
	"http":  struct{}{},
	"https": struct{}{},
}

// Checker represents a web page checker.
type Checker struct {
	fetcher     Fetcher
	rootPage    Page
	rootURL     *url.URL
	results     chan Result
	doneURLs    *sync.Map
	concurrency int
}

// NewChecker creates a new checker.
func NewChecker(f Fetcher, s string, c int) (Checker, error) {
	p, err := f.Fetch(s)

	if err != nil {
		return Checker{}, err
	}

	u, err := url.Parse(s)

	if err != nil {
		return Checker{}, err
	}

	return Checker{f, p, u, make(chan Result, c), &sync.Map{}, c}, nil
}

// Results returns a reference to results of web page checks.
func (c Checker) Results() <-chan Result {
	return c.results
}

// Check start checking web pages recursively from a root page.
func (c Checker) Check() {
	ps := make(chan Page, c.concurrency)
	ps <- c.rootPage

	w := sync.WaitGroup{}

	go func() {
		for p := range ps {
			w.Add(1)

			go func(p Page) {
				c.checkPage(p, ps)
				w.Done()
			}(p)
		}
	}()

	time.Sleep(10 * time.Millisecond)
	w.Wait()

	close(c.results)
}

// Check web pages recursively from the root.
func (c Checker) checkPage(p Page, ps chan Page) {
	n, err := html.Parse(bytes.NewReader(p.Body()))

	if err != nil {
		c.results <- NewResultWithError(p.URL(), err)
		return
	}

	r, err := url.Parse(p.URL())

	if err != nil {
		c.results <- NewResultWithError(p.URL(), err)
		return
	}

	ns := scrape.FindAll(n, func(n *html.Node) bool {
		return n.DataAtom == atom.A
	})

	sc, ec := make(chan string, len(ns)), make(chan string, len(ns))
	w := sync.WaitGroup{}

	for _, n := range ns {
		w.Add(1)

		go func(n *html.Node) {
			defer w.Done()

			u, err := url.Parse(scrape.Attr(n, "href"))

			if err != nil {
				ec <- err.Error()
				return
			}

			if _, ok := validSchemes[u.Scheme]; !ok {
				return
			}

			if !u.IsAbs() {
				u = r.ResolveReference(u)
			}

			p, err := c.fetcher.Fetch(u.String())

			if err == nil {
				sc <- fmt.Sprintf("link is alive (%v)", u)

				u.Fragment = ""
				u.RawQuery = ""

				if _, exist := c.doneURLs.LoadOrStore(u.String(), nil); !exist && u.Hostname() == c.rootURL.Hostname() {
					ps <- p
				}
			} else {
				ec <- fmt.Sprintf("%v (%v)", err, u)
			}
		}(n)
	}

	w.Wait()

	c.results <- NewResult(p.URL(), stringChannelToSlice(sc), stringChannelToSlice(ec))
}

func stringChannelToSlice(sc <-chan string) []string {
	ss := make([]string, 0, len(sc))

	for i := 0; i < len(sc); i++ {
		ss = append(ss, <-sc)
	}

	return ss
}
