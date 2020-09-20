package main

import (
	"sync"

	"github.com/fatih/color"
)

type checker struct {
	fetcher       fetcher
	urlValidator  urlValidator
	daemonManager daemonManager
	results       chan pageResult
	donePages     concurrentStringSet
}

func newChecker(f fetcher, ui urlValidator, concurrency int) checker {
	return checker{
		f,
		ui,
		newDaemonManager(concurrency),
		make(chan pageResult, concurrency),
		newConcurrentStringSet(),
	}
}

func (c checker) Results() <-chan pageResult {
	return c.results
}

func (c checker) Check(page *page) {
	c.addPage(page)
	c.daemonManager.Run()

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

			if !c.fetcher.options.OnePageOnly {
				if p, ok := r.Page(); ok && c.urlValidator.Validate(p.URL()) {
					c.addPage(p)
				}
			}
		}(u)
	}

	w.Wait()

	close(sc)
	close(ec)

	c.results <- newPageResult(p.URL().String(), stringChannelToSlice(sc), stringChannelToSlice(ec))
}

func (c checker) addPage(p *page) {
	if !c.donePages.Add(p.URL().String()) {
		c.daemonManager.Add(func() { c.checkPage(p) })
	}
}

func stringChannelToSlice(sc <-chan string) []string {
	ss := make([]string, 0, len(sc))

	for s := range sc {
		ss = append(ss, s)
	}

	return ss
}

func formatLinkSuccess(u string, s int) string {
	return color.GreenString("%v", s) + "\t" + u
}

func formatLinkError(u string, err error) string {
	return color.RedString(err.Error()) + "\t" + u
}
