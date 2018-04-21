package main

import (
	"os"
	"strconv"

	"github.com/docopt/docopt-go"
)

type arguments struct {
	concurrency int
	url         string
}

func getArguments() (arguments, error) {
	usage := `Muffet, the web repairgirl

Usage:
	muffet [-c <concurrency>] <url>

Options:
	-c, --concurrency <concurrency>  Roughly maximum number of concurrent open files. [default: 1000]
	-h, --help  Show this help.`

	args, err := docopt.ParseArgs(usage, os.Args[1:], "0.1.0")

	if err != nil {
		return arguments{}, err
	}

	c, err := strconv.ParseInt(args["--concurrency"].(string), 10, 32)

	if err != nil {
		return arguments{}, err
	}

	return arguments{int(c), args["<url>"].(string)}, nil
}
