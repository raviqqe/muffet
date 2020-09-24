package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/docopt/docopt-go"
)

var usage = fmt.Sprintf(`Muffet, the web repairgirl

Usage:
	muffet [-b <size>] [-c <concurrency>] [-e <pattern>...] [-f] [-j <header>...] [-l <times>] [-p] [-r] [-s] [-t <seconds>] [-v] [-x] <url>

Options:
	-b, --buffer-size <size>          HTTP response buffer size in bytes. [default: %v]
	-c, --max-connections <count>     Maximum number of concurrent HTTP connections. [default: %v]
	-e, --exclude <pattern>...        Exclude URLs matched with given regular expressions.
	-f, --ignore-fragments            Ignore URL fragments.
	-h, --help                        Show this help.
	-j, --header <header>...          Custom headers.
	-l, --max-redirections <count>    Maximum number of redirections. [default: %v]
	-p, --one-page-only               Only check links found in the given URL, do not follow links.
	-r, --follow-robots-txt           Follow robots.txt when scraping pages.
	-s, --follow-sitemap-xml          Scrape only pages listed in sitemap.xml.
	-t, --timeout <seconds>           Set timeout for HTTP requests in seconds. [default: %v]
	-v, --verbose                     Show successful results too.
	-x, --skip-tls-verification       Skip TLS certificates verification.`,
	defaultBufferSize, defaultMaxConnections, defaultMaxRedirections, defaultHTTPTimeout.Seconds())

type arguments struct {
	BufferSize       int
	MaxConnections   int
	ExcludedPatterns []*regexp.Regexp
	FollowRobotsTxt,
	FollowSitemapXML bool
	Headers         map[string]string
	IgnoreFragments bool
	MaxRedirections int
	Timeout         time.Duration
	URL             string
	Verbose,
	SkipTLSVerification bool
	OnePageOnly bool
}

func getArguments(regexps []string) (arguments, error) {
	args := parseArguments(usage, regexps)

	b, err := parseInt(args["--buffer-size"].(string))
	if err != nil {
		return arguments{}, err
	}

	c, err := parseInt(args["--max-connections"].(string))
	if err != nil {
		return arguments{}, err
	}

	ss, _ := args["--exclude"].([]string)
	rs, err := compileRegexps(ss)
	if err != nil {
		return arguments{}, err
	}

	ss, _ = args["--header"].([]string)
	hs, err := parseHeaders(ss)
	if err != nil {
		return arguments{}, err
	}

	r, err := parseInt(args["--max-redirections"].(string))
	if err != nil {
		return arguments{}, err
	}

	t, err := parseInt(args["--timeout"].(string))
	if err != nil {
		return arguments{}, err
	}

	return arguments{
		b,
		c,
		rs,
		args["--follow-robots-txt"].(bool),
		args["--follow-sitemap-xml"].(bool),
		hs,
		args["--ignore-fragments"].(bool),
		r,
		time.Duration(t) * time.Second,
		args["<url>"].(string),
		args["--verbose"].(bool),
		args["--skip-tls-verification"].(bool),
		args["--one-page-only"].(bool),
	}, nil
}

func parseArguments(u string, ss []string) map[string]interface{} {
	args, err := docopt.ParseArgs(u, ss, version)
	if err != nil {
		panic(err)
	}

	return args
}

func parseInt(s string) (int, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	return int(i), err
}

func compileRegexps(regexps []string) ([]*regexp.Regexp, error) {
	rs := make([]*regexp.Regexp, 0, len(regexps))

	for _, s := range regexps {
		r, err := regexp.Compile(s)
		if err != nil {
			return nil, err
		}

		rs = append(rs, r)
	}

	return rs, nil
}

func parseHeaders(headers []string) (map[string]string, error) {
	m := make(map[string]string, len(headers))

	for _, s := range headers {
		i := strings.IndexRune(s, ':')

		if i < 0 {
			return nil, errors.New("invalid header format")
		}

		m[s[:i]] = strings.TrimSpace(s[i+1:])
	}

	return m, nil
}
