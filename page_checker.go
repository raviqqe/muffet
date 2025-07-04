package main

import (
	"context"
	"errors"
	"net"
	"net/url"
	"sync"
)

type pageChecker struct {
	fetcher        *linkFetcher
	linkValidator  *linkValidator
	daemonManager  *daemonManager
	results        chan *pageResult
	donePages      concurrentStringSet
	onePageOnly    bool
	ignoreTimeouts ignoreTimeoutsGroup
}

func newPageChecker(f *linkFetcher, v *linkValidator, onePageOnly bool, ignoreTimeouts ignoreTimeoutsGroup) *pageChecker {
	return &pageChecker{
		f,
		v,
		newDaemonManager(concurrency),
		make(chan *pageResult, concurrency),
		newConcurrentStringSet(),
		onePageOnly,
		ignoreTimeouts,
	}
}

func (c *pageChecker) Results() <-chan *pageResult {
	return c.results
}

func (c *pageChecker) Check(page page) {
	c.addPage(page)
	c.daemonManager.Run()

	close(c.results)
}

func (c *pageChecker) checkPage(p page) {
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
			} else if !c.shouldIgnoreNetworkError(err, u) {
				ec <- &errorLinkResult{u, err}
			}

			if !c.onePageOnly && p != nil && c.linkValidator.Validate(p.URL()) {
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

func (c *pageChecker) addPage(p page) {
	if !c.donePages.Add(p.URL().String()) {
		c.daemonManager.Add(func() { c.checkPage(p) })
	}
}

func (c *pageChecker) shouldIgnoreNetworkError(err error, rawURL string) bool {
	if c.ignoreTimeouts == ignoreTimeoutsGroupNone || !isTimeoutError(err) {
		return false
	}

	if c.ignoreTimeouts == ignoreTimeoutsGroupAll {
		return true
	}

	u, err := url.Parse(rawURL)
	return err == nil &&
		(c.ignoreTimeouts == ignoreTimeoutsGroupAll ||
			u.Hostname() != c.linkValidator.hostname)
}

func isTimeoutError(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout() || errors.Is(err, context.DeadlineExceeded)
}
