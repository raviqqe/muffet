package main

import (
	"os"

	_ "github.com/breml/rootcerts"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

func main() {
	//change to nicer implementation later, for testing
	useTls := false
	for _, a := range os.Args[1:] {
		if a == "--use-tls-client" {
			useTls = true
		}
	}
	var factory httpClientFactory
	factory = newFasthttpHttpClientFactory()
	if useTls {
		factory = newTlsHttpClientFactory()
	}

	ok := newCommand(
		colorable.NewColorableStdout(),
		os.Stderr,
		isatty.IsTerminal(os.Stdout.Fd()),
		factory,
	).Run(os.Args[1:])

	if !ok {
		os.Exit(1)
	}
}
