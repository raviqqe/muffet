package main

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/logrusorgru/aurora/v3"
	"github.com/temoto/robotstxt"
)

type command struct {
	stdout, stderr    io.Writer
	terminal          bool
	httpClientFactory httpClientFactory
}

func newCommand(stdout, stderr io.Writer, terminal bool, f httpClientFactory) *command {
	return &command{stdout, stderr, terminal, f}
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
	} else if args.Help {
		printHelp(c.stderr)
		return true, nil
	} else if args.Version {
		_, err := fmt.Fprintln(c.stderr, version)
		return err == nil, err
	}

	client := newThrottledHTTPClient(
		c.httpClientFactory.Create(
			httpClientOptions{
				MaxConnectionsPerHost: args.MaxConnectionsPerHost,
				BufferSize:            args.BufferSize,
				MaxRedirections:       args.MaxRedirections,
				SkipTLSVerification:   args.SkipTLSVerification,
				Timeout:               time.Duration(args.Timeout) * time.Second,
			},
		),
		args.MaxConnections,
	)

	pp := newPageParser(newLinkFinder(args.ExcludedPatterns))

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
		return false, fmt.Errorf("failed to fetch root page: %v", err)
	} else if p == nil {
		return false, errors.New("root page is not HTML")
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
		args.OnePageOnly,
	)

	go checker.Check(p)

	formatter := newPageResultFormatter(args.Verbose, c.terminal)
	ok := true

	for r := range checker.Results() {
		if !r.OK() || args.Verbose {
			c.print(formatter.Format(r))
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
	s := fmt.Sprint(xs...)

	if c.terminal {
		s = aurora.Red(s).String()
	}

	if _, err := fmt.Fprintln(c.stderr, s); err != nil {
		panic(err)
	}
}
