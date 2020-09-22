package main

import (
	"os"

	"github.com/mattn/go-colorable"
)

func main() {
	ok := newCommand(
		colorable.NewColorableStdout(),
		os.Stderr,
		newFasthttpHTTPClientFactory(),
	).Run(os.Args[1:])

	if !ok {
		os.Exit(1)
	}
}
