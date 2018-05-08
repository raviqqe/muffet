package main

import (
	"fmt"
	"os"
)

func main() {
	args, err := getArguments(os.Args[1:])

	if err != nil {
		fail(err)
	}

	c, err := newChecker(args.URL, checkerOptions{
		fetcherOptions{
			args.Concurrency,
			args.IgnoreFragments,
			args.MaxRedirections,
			args.SkipTLSVerification,
			args.Timeout,
		},
		args.FollowRobotsTxt,
		args.FollowSitemapXML,
	})

	if err != nil {
		fail(err)
	}

	go c.Check()

	b := false

	for r := range c.Results() {
		if !r.OK() || args.Verbose {
			fmt.Println(r.String(args.Verbose))
		}

		if !r.OK() {
			b = true
		}
	}

	if b {
		os.Exit(1)
	}
}

func fail(x interface{}) {
	fmt.Fprintln(os.Stderr, x)
	os.Exit(1)
}
