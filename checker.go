package main

import (
	"crypto/tls"
	"errors"
	"sync"

	"github.com/fatih/color"
	"github.com/valyala/fasthttp"
)

type checker struct {
	fetcher
	daemons      daemons
	urlInspector urlInspector
	results      chan pageResult
	donePages    concurrentStringSet
}

func newChecker(s string, o checkerOptions) (checker, error) {
	o.Initialize()

	c := &fasthttp.Client{
		MaxConnsPerHost: o.Concurrency,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: o.SkipTLSVerification,
		},
	}
	f := newFetcher(c, o.fetcherOptions)
	r, err := f.Fetch(s)

	if err != nil {
		return checker{}, err
	}

	p, ok := r.Page()

	if !ok {
		return checker{}, errors.New("non-HTML page")
	}

	ui, err := newURLInspector(c, p.URL().String(), o.FollowRobotsTxt, o.FollowSitemapXML)

	if err != nil {
		return checker{}, err
	}

	ch := checker{
		f,
		newDaemons(o.Concurrency),
		ui,
		make(chan pageResult, o.Concurrency),
		newConcurrentStringSet(),
	}

	ch.addPage(p)

	return ch, nil
}

func (c checker) Results() <-chan pageResult {
	return c.results
}

func (c checker) Check() {
	c.daemons.Run()

	close(c.results)
}

func (c checker) checkPage(p *page) {
	us := p.Links()

	sc := make(chan string, len(us))
	ec := make(chan string, len(us))
	w := sync.WaitGroup{}

	for u, err := range us {
		if err != nil {
			ec <- formatLinkError(u, err)
			continue
		}

		w.Add(1)

		go func(u string) {
			defer w.Done()

			r, err := c.fetcher.Fetch(u)

			if err == nil {
				sc <- formatLinkSuccess(u, r.StatusCode())
			} else {
				ec <- formatLinkError(u, err)
			}

			if p, ok := r.Page(); ok && c.urlInspector.Inspect(p.URL()) {
				c.addPage(p)
			}
		}(u)
	}

	w.Wait()

	c.results <- newPageResult(p.URL().String(), stringChannelToSlice(sc), stringChannelToSlice(ec))
}

func (c checker) addPage(p *page) {
	if !c.donePages.Add(p.URL().String()) {
		c.daemons.Add(func() { c.checkPage(p) })
	}
}

func stringChannelToSlice(sc <-chan string) []string {
	ss := make([]string, 0, len(sc))

	for i := 0; i < cap(ss); i++ {
		ss = append(ss, <-sc)
	}

	return ss
}

func formatLinkSuccess(u string, s int) string {
	return color.GreenString("%v", s) + "\t" + u
}

func formatLinkError(u string, err error) string {
	return color.RedString(err.Error()) + "\t" + u
}
