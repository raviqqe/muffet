package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
)

func main() {
	args, err := getArgs()

	if err != nil {
		printToStderr(err)
		os.Exit(1)
	}

	c, err := NewChecker(args["<url>"].(string))

	if err != nil {
		printToStderr(err)
		os.Exit(1)
	}

	go c.Check()

	for r := range c.Results() {
		printToStderr(r)
	}
}

func getArgs() (map[string]interface{}, error) {
	usage := `Muffet, the web repairman

Usage:
	muffet <url>

Options:
	-h, --help  Show this help.`

	args, err := docopt.ParseArgs(usage, os.Args[1:], "0.1.0")

	if err != nil {
		return nil, err
	}

	return args, nil
}

func printToStderr(x interface{}) {
	fmt.Fprintln(os.Stderr, x)
}
