package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/docopt/docopt-go"
)

var usage = fmt.Sprintf(`Muffet, the web repairgirl

Usage:
	muffet [-c <concurrency>] [-f] [-l <times>] [-r] [-s] [-t] [-v] <url>

Options:
	-c, --concurrency <concurrency>   Roughly maximum number of concurrent HTTP connections. [default: %v]
	-f, --ignore-fragments            Ignore URL fragments.
	-h, --help                        Show this help.
	-l, --limit-redirections <times>  Limit a number of redirections. [default: %v]
	-r, --follow-robots-txt           Follow robots.txt when scraping.
	-s, --follow-sitemap-xml          Scrape only pages listed in sitemap.xml.
	-t, --skip-tls-verification       Skip TLS certificates verification.
	-v, --verbose                     Show successful results too.`,
	defaultConcurrency, defaultMaxRedirections)

type arguments struct {
	Concurrency int
	FollowRobotsTxt,
	FollowSitemapXML,
	IgnoreFragments bool
	MaxRedirections     int
	SkipTLSVerification bool
	URL                 string
	Verbose             bool
}

func getArguments(ss []string) (arguments, error) {
	if ss == nil {
		ss = os.Args[1:]
	}

	args, err := docopt.ParseArgs(usage, ss, "0.3.0")

	if err != nil {
		return arguments{}, err
	}

	c, err := strconv.ParseInt(args["--concurrency"].(string), 10, 32)

	if err != nil {
		return arguments{}, err
	}

	r, err := strconv.ParseInt(args["--limit-redirections"].(string), 10, 32)

	if err != nil {
		return arguments{}, err
	}

	return arguments{
		int(c),
		args["--follow-robots-txt"].(bool),
		args["--follow-sitemap-xml"].(bool),
		args["--ignore-fragments"].(bool),
		int(r),
		args["--skip-tls-verification"].(bool),
		args["<url>"].(string),
		args["--verbose"].(bool),
	}, nil
}
