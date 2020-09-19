package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-colorable"
)

func main() {
	ok, err := newCommand(colorable.NewColorableStdout()).Run(os.Args[1:])
	if err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			panic(err)
		}

		os.Exit(1)
	}

	if !ok {
		os.Exit(1)
	}
}
