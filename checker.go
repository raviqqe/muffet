package main

import (
	"sync"
)

type checker struct {
	fetcher       *linkFetcher
	linkValidator *linkValidator
	daemonManager *daemonManager
	results       chan *pageResult
	donePages     concurrentStringSet
	onePageOnly   bool
}

func newChecker(f *linkFetcher, v *linkValidator, concurrency int, onePageOnly bool) *checker {
	return &checker{
		f,
		v,
		newDaemonManager(concurrency),
		make(chan *pageResult, concurrency),
		newConcurrentStringSet(),
		onePageOnly,
	}
}

func (c *checker) Results() <-chan *pageResult {
	return c.results
}

func (c *checker) Check(page *page) {
	c.addPage(page)
	c.daemonManager.Run()

	close(c.results)
}

func (c *checker) checkPage(p *page) {
	us := p.Links()

	sc := make(chan *successLinkResult, len(us))
	ec := make(chan *errorLinkResult, len(us))
	w := sync.WaitGroup{}

	for u, err := range us {
		if err != nil {
			ec <- &errorLinkResult{u, err}
			continue
		}

		w.Add(1)

		go func(u string) {
			defer w.Done()

			status, p, err := c.fetcher.Fetch(u)

			if err == nil {
				sc <- &successLinkResult{u, status}
			} else {
				ec <- &errorLinkResult{u, err}
			}

			if !c.onePageOnly &&
				p != nil &&
				c.linkValidator.Validate(p.URL()) &&
				!c.donePages.Add(p.URL().String()) {
				c.addPage(p)
			}
		}(u)
	}

	w.Wait()

	close(sc)
	close(ec)

	ss := make([]*successLinkResult, 0, len(sc))

	for s := range sc {
		ss = append(ss, s)
	}

	es := make([]*errorLinkResult, 0, len(ec))

	for e := range ec {
		es = append(es, e)
	}

	c.results <- &pageResult{p.URL().String(), ss, es}
}

func (c *checker) addPage(p *page) {
	c.daemonManager.Add(func() { c.checkPage(p) })
}
