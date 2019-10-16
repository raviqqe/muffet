package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	s, err := command(os.Args[1:], os.Stdout)

	if err != nil {
		fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(s)
}

func command(ss []string, w io.Writer) (int, error) {
	args, err := getArguments(ss)

	if err != nil {
		return 0, err
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
		args.FollowRobotsTxt,
		args.FollowSitemapXML,
		args.FollowURLParams,
		args.SkipTLSVerification,
	})

	if err != nil {
		return 0, err
	}

	go c.Check()

	s := 0

	for r := range c.Results() {
		if !r.OK() || args.Verbose {
			fprintln(w, r.String(args.Verbose))
		}

		if !r.OK() {
			s = 1
		}
	}

	return s, nil
}

func fprintln(w io.Writer, xs ...interface{}) {
	if _, err := fmt.Fprintln(w, xs...); err != nil {
		panic(err)
	}
}
