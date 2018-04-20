package main

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Checker represents a web page checker.
type Checker struct {
	rootPage Page
	rootURL  *url.URL
	results  chan Result
	doneURLs *sync.Map
}

// NewChecker creates a new checker.
func NewChecker(s string) (Checker, error) {
	p, err := fetch(s)

	if err != nil {
		return Checker{}, err
	}

	u, err := url.Parse(s)

	if err != nil {
		return Checker{}, err
	}

	return Checker{p, u, make(chan Result, 256), &sync.Map{}}, nil
}

// Results returns a reference to results of web page checks.
func (c Checker) Results() <-chan Result {
	return c.results
}

// Check start checking web pages recursively from a root page.
func (c Checker) Check() {
	ps := make(chan Page, 256)
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
	n, err := html.Parse(p.Body())

	if err != nil {
		c.results <- NewResultWithError(p.URL(), err)
		return
	}

	ss, es := []string{}, []string{}

	for _, n := range scrape.FindAll(n, func(n *html.Node) bool {
		return n.DataAtom == atom.A
	}) {
		u, err := url.Parse(scrape.Attr(n, "href"))

		if err != nil {
			es = append(es, err.Error())
			continue
		}

		if !u.IsAbs() {
			u = c.rootURL.ResolveReference(u)
		}

		p, err := fetch(u.String())

		if err == nil {
			ss = append(ss, fmt.Sprintf("%s is alive", u.String()))

			if _, exist := c.doneURLs.LoadOrStore(u.String(), nil); !exist && u.Hostname() == c.rootURL.Hostname() {
				ps <- p
			}
		} else {
			es = append(es, err.Error())
		}
	}

	c.results <- NewResult(p.URL(), ss, es)
}
