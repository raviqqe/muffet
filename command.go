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
	writer io.Writer
}

func newCommand(writer io.Writer) command {
	return command{writer}
}

func (c command) Run(rawArgs []string) (bool, error) {
	args, err := getArguments(rawArgs)
	if err != nil {
		return false, err
	}

	client := newFasthttpHTTPClient(
		&fasthttp.Client{
			MaxConnsPerHost: args.Concurrency,
			ReadBufferSize:  args.BufferSize,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: args.SkipTLSVerification,
			},
		},
		args.MaxRedirections,
		args.Timeout,
	)

	pp := newPageParser(newLinkFinder(args.ExcludedPatterns), args.FollowURLParams)

	f := newLinkFetcher(
		client,
		pp,
		linkFetcherOptions{
			args.Concurrency,
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
	if _, err := fmt.Fprintln(c.writer, xs...); err != nil {
		panic(err)
	}
}
