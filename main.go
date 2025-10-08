package main

import (
	"os"

	_ "github.com/breml/rootcerts"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

func main() {
	//for testing if works set to always use tlsHttpClientFactory, probably should set to a flag to use this or fasthttp
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
