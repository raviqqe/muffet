package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"

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

	client := newFasthttpHTTPClient(&fasthttp.Client{
		MaxConnsPerHost: args.Concurrency,
		ReadBufferSize:  args.BufferSize,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: args.SkipTLSVerification,
		},
	})

	f := newFetcher(client, fetcherOptions{
		args.Concurrency,
		args.ExcludedPatterns,
		args.Headers,
		args.IgnoreFragments,
		args.FollowURLParams,
		args.MaxRedirections,
		args.Timeout,
		args.OnePageOnly,
	})

	r, err := f.Fetch(args.URL)
	if err != nil {
		return false, err
	}

	p, ok := r.Page()

	if !ok {
		return false, errors.New("non-HTML page")
	}

	ui, err := newURLValidator(client, p.URL().String(), args.FollowRobotsTxt, args.FollowSitemapXML)
	if err != nil {
		return false, err
	}

	checker := newChecker(f, ui, args.Concurrency)

	go checker.Check(p)

	ok = true

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
