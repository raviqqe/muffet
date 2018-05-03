package main

import (
	"sync"

	"github.com/fatih/color"
)

type checker struct {
	fetcher   fetcher
	daemons   daemons
	hostname  string
	results   chan pageResult
	donePages concurrentStringSet
}

func newChecker(s string, c int, i bool) (checker, error) {
	f := newFetcher(c, i)
	p, err := f.FetchPage(s)

	if err != nil {
		return checker{}, err
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
	bs, es := scrapePage(p)

	ec := make(chan string, len(bs)+len(es))

	for u, err := range es {
		ec <- formatLinkError(u, err)
	}

	sc := make(chan string, len(bs))
	w := sync.WaitGroup{}

	for u, b := range bs {
		w.Add(1)

		go func(u string, isHTML bool) {
			defer w.Done()

			r, err := c.fetcher.FetchLink(u)

			if err == nil {
				sc <- formatLinkSuccess(u, r.StatusCode())
			} else {
				ec <- formatLinkError(u, err)
			}

			if p, ok := r.Page(); ok && isHTML && !c.donePages.Add(p.URL().String()) && p.URL().Hostname() == c.hostname {
				c.daemons.Add(func() { c.checkPage(p) })
			}
		}(u, b)
	}

	w.Wait()

	c.results <- newPageResult(p.URL().String(), stringChannelToSlice(sc), stringChannelToSlice(ec))
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
