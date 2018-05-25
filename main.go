package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	s, err := command(os.Args[1:], os.Stdout)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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
			args.Headers,
			args.IgnoreFragments,
			args.MaxRedirections,
			args.SkipTLSVerification,
			args.Timeout,
		},
		args.FollowRobotsTxt,
		args.FollowSitemapXML,
	})

	if err != nil {
		return 0, err
	}

	go c.Check()

	s := 0

	for r := range c.Results() {
		if !r.OK() || args.Verbose {
			fmt.Fprintln(w, r.String(args.Verbose))
		}

		if !r.OK() {
			s = 1
		}
	}

	return s, nil
}
