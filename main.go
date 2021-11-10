package main

import (
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

func main() {
	ok := newCommand(
		colorable.NewColorableStdout(),
		os.Stderr,
		isatty.IsTerminal(os.Stdout.Fd()),
		newFasthttpHttpClientFactory(),
	).Run(os.Args[1:])

	if !ok {
		os.Exit(1)
	}
}
