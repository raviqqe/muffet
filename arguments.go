package main

import (
	"strconv"

	"github.com/docopt/docopt-go"
)

type arguments struct {
	concurrency int
	url         string
	verbose     bool
}

func getArguments() (arguments, error) {
	usage := `Muffet, the web repairgirl

Usage:
	muffet [-c <concurrency>] [-v] <url>

Options:
	-c, --concurrency <concurrency>  Roughly maximum number of concurrent HTTP connections. [default: 1000]
	-h, --help  Show this help.
	-v, --verbose  Show successful results too.`

	args, err := docopt.ParseArgs(usage, nil, "0.1.0")

	if err != nil {
		return arguments{}, err
	}

	c, err := strconv.ParseInt(args["--concurrency"].(string), 10, 32)

	if err != nil {
		return arguments{}, err
	}

	return arguments{int(c), args["<url>"].(string), args["--verbose"].(bool)}, nil
}
