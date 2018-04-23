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

	c, err := newChecker(args.url, args.concurrency)

	if err != nil {
		fail(err)
	}

	go c.Check()

	b := false

	for r := range c.Results() {
		if !r.OK() || args.verbose {
			fmt.Println(r.String(args.verbose))
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
