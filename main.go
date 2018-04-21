package main

import (
	"fmt"
	"os"
)

func main() {
	args, err := getArguments()

	if err != nil {
		printToStderr(err)
		os.Exit(1)
	}

	c, err := newChecker(args.url, args.concurrency, args.connectionsPerHost)

	if err != nil {
		printToStderr(err)
		os.Exit(1)
	}

	go c.Check()

	b := false

	for r := range c.Results() {
		if !r.OK() || args.verbose {
			printToStderr(r.String(args.verbose))
		}

		if !r.OK() {
			b = true
		}
	}

	if b {
		os.Exit(1)
	}
}

func printToStderr(x interface{}) {
	fmt.Fprintln(os.Stderr, x)
}
