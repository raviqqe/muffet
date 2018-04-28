package main

import (
	"net/url"
	"sync"

	"github.com/fatih/color"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var validSchemes = map[string]struct{}{
	"":      {},
	"http":  {},
	"https": {},
}

var atomToAttribute = map[atom.Atom]string{
	atom.A:      "href",
	atom.Frame:  "src",
	atom.Iframe: "src",
	atom.Img:    "src",
	atom.Link:   "href",
	atom.Script: "src",
	atom.Source: "src",
	atom.Track:  "src",
}

type checker struct {
	fetcher   fetcher
	daemons   daemons
	hostname  string
	results   chan pageResult
	donePages concurrentStringSet
}

func newChecker(s string, c int, i bool) (checker, error) {
	f := newFetcher(c, i)
	r, err := f.Fetch(s)

	if err != nil {
		return checker{}, err
	}

	p, ok := r.Page()

	if !ok {
		panic("unreachable")
	}

	ch := checker{
		f,
		newDaemons(c),
		p.URL().Hostname(),
		make(chan pageResult, c),
		newConcurrentStringSet(),
	}

	ch.daemons.Add(func() { ch.checkPage(p) })

	return ch, nil
}

func (c checker) Results() <-chan pageResult {
	return c.results
}

func (c checker) Check() {
	c.daemons.Run()

	close(c.results)
}

func (c checker) checkPage(p page) {
	ns := scrape.FindAllNested(p.Body(), func(n *html.Node) bool {
		_, ok := atomToAttribute[n.DataAtom]
		return ok
	})

	sc, ec := make(chan string, len(ns)), make(chan string, len(ns))
	w := sync.WaitGroup{}

	for _, n := range ns {
		w.Add(1)

		go func(n *html.Node) {
			defer w.Done()

			u, err := url.Parse(scrape.Attr(n, atomToAttribute[n.DataAtom]))

			if err != nil {
				ec <- err.Error()
				return
			}

			if _, ok := validSchemes[u.Scheme]; !ok {
				return
			}

			u, err = resolveURL(p, u)

			if err != nil {
				ec <- err.Error()
				return
			}

			r, err := c.fetcher.Fetch(u.String())

			if err == nil {
				sc <- color.GreenString("%v", r.StatusCode()) + "\t" + u.String()
			} else {
				ec <- color.RedString(err.Error()) + "\t" + u.String()
			}

			if p, ok := r.Page(); ok && n.DataAtom == atom.A && !c.donePages.Add(p.URL().String()) && p.URL().Hostname() == c.hostname {
				c.daemons.Add(func() {
					c.checkPage(p)
				})
			}
		}(n)
	}

	w.Wait()

	c.results <- newPageResult(p.URL().String(), stringChannelToSlice(sc), stringChannelToSlice(ec))
}

func resolveURL(p page, u *url.URL) (*url.URL, error) {
	if u.IsAbs() {
		return u, nil
	}

	b := p.URL()

	if n, ok := scrape.Find(p.Body(), func(n *html.Node) bool {
		return n.DataAtom == atom.Base
	}); ok {
		u, err := url.Parse(scrape.Attr(n, "href"))

		if err != nil {
			return nil, err
		}

		b = b.ResolveReference(u)
	}

	u = b.ResolveReference(u)

	return u, nil
}

func stringChannelToSlice(sc <-chan string) []string {
	ss := make([]string, 0, len(sc))

	for i := 0; i < cap(ss); i++ {
		ss = append(ss, <-sc)
	}

	return ss
}
