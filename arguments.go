package main

import (
	"os"
	"strconv"

	"github.com/docopt/docopt-go"
)

const usage = `Muffet, the web repairgirl

Usage:
	muffet [-c <concurrency>] [-f] [-v] <url>

Options:
	-c, --concurrency <concurrency>  Roughly maximum number of concurrent HTTP connections. [default: 512]
	-f, --ignore-fragments           Ignore URL fragments.
	-h, --help                       Show this help.
	-v, --verbose                    Show successful results too.`

type arguments struct {
	concurrency     int
	url             string
	verbose         bool
	ignoreFragments bool
}

func getArguments(ss []string) (arguments, error) {
	if ss == nil {
		ss = os.Args[1:]
	}

	args, err := docopt.ParseArgs(usage, ss, "0.1.0")

	if err != nil {
		return arguments{}, err
	}

	c, err := strconv.ParseInt(args["--concurrency"].(string), 10, 32)

	if err != nil {
		return arguments{}, err
	}

	return arguments{int(c), args["<url>"].(string), args["--verbose"].(bool), args["--ignore-fragments"].(bool)}, nil
}
