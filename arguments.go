package main

import (
	"fmt"
	"strconv"

	"github.com/docopt/docopt-go"
	"github.com/valyala/fasthttp"
)

type arguments struct {
	concurrency        int
	connectionsPerHost int
	url                string
	verbose            bool
}

func getArguments() (arguments, error) {
	usage := fmt.Sprintf(`Muffet, the web repairgirl

Usage:
	muffet [-c <concurrency>] [-n <connections>] [-v] <url>

Options:
	-c, --concurrency <concurrency>  Roughly maximum number of concurrent HTTP connections. [default: 1000]
	-h, --help  Show this help.
	-n, --connections-per-host <connections>  Maximum number of concurrent connections per host. [default: %v]
	-v, --verbose  Show successful results too.`, fasthttp.DefaultMaxConnsPerHost)

	args, err := docopt.ParseArgs(usage, nil, "0.1.0")

	if err != nil {
		return arguments{}, err
	}

	c, err := strconv.ParseInt(args["--concurrency"].(string), 10, 32)

	if err != nil {
		return arguments{}, err
	}

	n, err := strconv.ParseInt(args["--connections-per-host"].(string), 10, 32)

	if err != nil {
		return arguments{}, err
	}

	return arguments{
		int(c),
		int(n),
		args["<url>"].(string),
		args["--verbose"].(bool),
	}, nil
}
