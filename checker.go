package main

import (
	"errors"
	"sync"

	"github.com/fatih/color"
)

type checker struct {
	fetcher
	scraper
	daemons      daemons
	urlInspector urlInspector
	results      chan pageResult
	donePages    concurrentStringSet
}

func newChecker(s string, o checkerOptions) (checker, error) {
	o.Initialize()

	f := newFetcher(o.fetcherOptions)
	r, err := f.Fetch(s)

	if err != nil {
		return checker{}, err
	}

	p, ok := r.Page()

	if !ok {
		return checker{}, errors.New("non-HTML page")
	}

	ui, err := newURLInspector(p.URL().String(), o.FollowRobotsTxt, o.FollowSitemapXML)

	if err != nil {
		return checker{}, err
	}

	ch := checker{
		f,
		newScraper(o.ExcludedPatterns),
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

func (c checker) checkPage(p page) {
	ss, es := c.scraper.Scrape(p)

	ec := make(chan string, len(ss)+len(es))

	for u, err := range es {
		ec <- formatLinkError(u, err)
	}

	sc := make(chan string, len(ss))
	w := sync.WaitGroup{}

	for u := range ss {
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

func (c checker) addPage(p page) {
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
