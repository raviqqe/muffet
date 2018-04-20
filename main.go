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

	c, err := NewChecker(newFetcher(args.concurrency), args.url, args.concurrency)

	if err != nil {
		printToStderr(err)
		os.Exit(1)
	}

	go c.Check()

	b := false

	for r := range c.Results() {
		printToStderr(r)

		b = b && r.IsError()
	}

	if b {
		os.Exit(1)
	}
}

func printToStderr(x interface{}) {
	fmt.Fprintln(os.Stderr, x)
}
