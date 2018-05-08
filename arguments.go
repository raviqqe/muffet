package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/docopt/docopt-go"
)

var usage = fmt.Sprintf(`Muffet, the web repairgirl

Usage:
	muffet [-c <concurrency>] [-f] [-l <times>] [-r] [-s] [-t <seconds>] [-v] [-x] <url>

Options:
	-c, --concurrency <concurrency>   Roughly maximum number of concurrent HTTP connections. [default: %v]
	-f, --ignore-fragments            Ignore URL fragments.
	-h, --help                        Show this help.
	-l, --limit-redirections <times>  Limit a number of redirections. [default: %v]
	-r, --follow-robots-txt           Follow robots.txt when scraping.
	-s, --follow-sitemap-xml          Scrape only pages listed in sitemap.xml.
	-t, --timeout <seconds>           Set timeout for HTTP requests in seconds. [default: %v]
	-v, --verbose                     Show successful results too.
	-x, --skip-tls-verification       Skip TLS certificates verification.`,
	defaultConcurrency, defaultMaxRedirections, defaultTimeout.Seconds())

type arguments struct {
	Concurrency int
	FollowRobotsTxt,
	FollowSitemapXML,
	IgnoreFragments bool
	MaxRedirections int
	Timeout         time.Duration
	URL             string
	Verbose,
	SkipTLSVerification bool
}

func getArguments(ss []string) (arguments, error) {
	if ss == nil {
		ss = os.Args[1:]
	}

	args, err := docopt.ParseArgs(usage, ss, "0.3.0")

	if err != nil {
		return arguments{}, err
	}

	c, err := parseInt(args["--concurrency"].(string))

	if err != nil {
		return arguments{}, err
	}

	r, err := parseInt(args["--limit-redirections"].(string))

	if err != nil {
		return arguments{}, err
	}

	t, err := parseInt(args["--timeout"].(string))

	if err != nil {
		return arguments{}, err
	}

	return arguments{
		c,
		args["--follow-robots-txt"].(bool),
		args["--follow-sitemap-xml"].(bool),
		args["--ignore-fragments"].(bool),
		r,
		time.Duration(t) * time.Second,
		args["<url>"].(string),
		args["--verbose"].(bool),
		args["--skip-tls-verification"].(bool),
	}, nil
}

func parseInt(s string) (int, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	return int(i), err
}
