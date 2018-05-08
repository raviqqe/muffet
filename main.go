package main

import (
	"fmt"
	"os"
)

func main() {
	args, err := getArguments(nil)

	if err != nil {
		fail(err)
	}

	c, err := newChecker(args.URL, checkerOptions{
		args.Concurrency,
		args.FollowRobotsTxt,
		args.FollowSitemapXML,
		args.IgnoreFragments,
		args.MaxRedirections,
		args.SkipTLSVerification,
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
