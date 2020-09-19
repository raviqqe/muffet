package main

import (
	"fmt"
	"io"
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

	checker, err := newChecker(args.URL, checkerOptions{
		fetcherOptions{
			args.Concurrency,
			args.ExcludedPatterns,
			args.Headers,
			args.IgnoreFragments,
			args.FollowURLParams,
			args.MaxRedirections,
			args.Timeout,
			args.OnePageOnly,
		},
		args.BufferSize,
		args.FollowRobotsTxt,
		args.FollowSitemapXML,
		args.FollowURLParams,
		args.SkipTLSVerification,
	})

	if err != nil {
		return false, err
	}

	go checker.Check()

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
