package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"

	"github.com/temoto/robotstxt"
	"github.com/valyala/fasthttp"
)

type command struct {
	stdout, stderr io.Writer
}

func newCommand(stdout, stderr io.Writer) *command {
	return &command{stdout, stderr}
}

func (c *command) Run(args []string) bool {
	ok, err := c.runWithError(args)

	if err != nil {
		c.printError(err)
	}

	return ok
}

func (c *command) runWithError(ss []string) (bool, error) {
	args, err := getArguments(ss)
	if err != nil {
		return false, err
	}

	client := newThrottledHTTPClient(
		newFasthttpHTTPClient(
			&fasthttp.Client{
				MaxConnsPerHost: args.Concurrency,
				ReadBufferSize:  args.BufferSize,
				TLSConfig: &tls.Config{
					InsecureSkipVerify: args.SkipTLSVerification,
				},
			},
			args.MaxRedirections,
			args.Timeout,
		),
		args.Concurrency,
	)

	pp := newPageParser(newLinkFinder(args.ExcludedPatterns), args.FollowURLParams)

	f := newLinkFetcher(
		client,
		pp,
		linkFetcherOptions{
			args.Headers,
			args.IgnoreFragments,
		},
	)

	_, p, err := f.Fetch(args.URL)
	if err != nil {
		return false, err
	} else if p == nil {
		return false, errors.New("non-HTML page")
	}

	rd := (*robotstxt.RobotsData)(nil)

	if args.FollowRobotsTxt {
		rd, err = newRobotsTxtFetcher(client).Fetch(p.URL())

		if err != nil {
			return false, err
		}
	}

	sm := (map[string]struct{})(nil)

	if args.FollowSitemapXML {
		sm, err = newSitemapFetcher(client).Fetch(p.URL())

		if err != nil {
			return false, err
		}
	}

	checker := newChecker(
		f,
		newLinkValidator(p.URL().Hostname(), rd, sm),
		args.Concurrency,
		args.OnePageOnly,
	)

	go checker.Check(p)

	ok := true

	for r := range checker.Results() {
		if !r.OK() || args.Verbose {
			c.print(r.String(args.Verbose))
		}

		if !r.OK() {
			ok = false
		}
	}

	return ok, nil
}

func (c command) print(xs ...interface{}) {
	if _, err := fmt.Fprintln(c.stdout, xs...); err != nil {
		panic(err)
	}
}

func (c command) printError(xs ...interface{}) {
	if _, err := fmt.Fprintln(c.stderr, xs...); err != nil {
		panic(err)
	}
}
