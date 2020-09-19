package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-colorable"
)

func main() {
	ok, err := command(os.Args[1:], colorable.NewColorableStdout())

	if err != nil {
		fprintln(os.Stderr, err)
	}

	if !ok {
		os.Exit(1)
	}
}

func command(ss []string, w io.Writer) (bool, error) {
	args, err := getArguments(ss)

	if err != nil {
		return false, err
	}

	c, err := newChecker(args.URL, checkerOptions{
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

	go c.Check()

	ok := true

	for r := range c.Results() {
		if !r.OK() || args.Verbose {
			fprintln(w, r.String(args.Verbose))
		}

		if !r.OK() {
			ok = false
		}
	}

	return ok, nil
}

func fprintln(w io.Writer, xs ...interface{}) {
	if _, err := fmt.Fprintln(w, xs...); err != nil {
		panic(err)
	}
}
