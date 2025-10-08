package main

import (
	"os"

	_ "github.com/breml/rootcerts"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

func main() {
	ok := newCommand(
		colorable.NewColorableStdout(),
		os.Stderr,
		isatty.IsTerminal(os.Stdout.Fd()),
		newTlsHttpClientFactory(),
	).Run(os.Args[1:])

	if !ok {
		os.Exit(1)
	}
}
