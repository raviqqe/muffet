package main

import (
	"bytes"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/jessevdk/go-flags"
)

type arguments struct {
	BufferSize            int      `short:"b" long:"buffer-size" value-name:"<size>" default:"4096" description:"HTTP response buffer size in bytes"`
	MaxConnections        int      `short:"c" long:"max-connections" value-name:"<count>" default:"512" description:"Maximum number of HTTP connections"`
	MaxConnectionsPerHost int      `long:"max-connections-per-host" value-name:"<count>" default:"512" description:"Maximum number of HTTP connections per host"`
	MaxResponseBodySize   int      `long:"max-response-body-size" value-name:"<size>" default:"10000000" description:"Maximum response body size to read"`
	RawExcludedPatterns   []string `short:"e" long:"exclude" value-name:"<pattern>..." description:"Exclude URLs matched with given regular expressions"`
	RawIncludedPatterns   []string `short:"i" long:"include" value-name:"<pattern>..." description:"Include URLs matched with given regular expressions"`
	FollowRobotsTxt       bool     `long:"follow-robots-txt" description:"Follow robots.txt when scraping pages"`
	FollowSitemapXML      bool     `long:"follow-sitemap-xml" description:"Scrape only pages listed in sitemap.xml"`
	RawHeaders            []string `long:"header" value-name:"<header>..." description:"Custom headers"`
	// TODO Remove a short option.
	IgnoreFragments bool   `short:"f" long:"ignore-fragments" description:"Ignore URL fragments"`
	Format          string `long:"format" description:"Output format" default:"text" choice:"text" choice:"json" choice:"junit"`
	// TODO Remove this option.
	JSONOutput bool `long:"json" description:"Output results in JSON (deprecated)"`
	// TODO Remove this option.
	VerboseJSON bool `long:"experimental-verbose-json" description:"Include successful results in JSON (deprecated)"`
	// TODO Remove this option.
	JUnitOutput         bool   `long:"junit" description:"Output results as JUnit XML file (deprecated)"`
	MaxRedirections     int    `short:"r" long:"max-redirections" value-name:"<count>" default:"64" description:"Maximum number of redirections"`
	RateLimit           int    `long:"rate-limit" value-name:"<rate>" description:"Max requests per second"`
	Timeout             int    `short:"t" long:"timeout" value-name:"<seconds>" default:"10" description:"Timeout for HTTP requests in seconds"`
	Verbose             bool   `short:"v" long:"verbose" description:"Show successful results too"`
	Proxy               string `long:"proxy" value-name:"<host>" description:"HTTP proxy host"`
	SkipTLSVerification bool   `long:"skip-tls-verification" description:"Skip TLS certificate verification"`
	OnePageOnly         bool   `long:"one-page-only" description:"Only check links found in the given URL"`
	Color               color  `long:"color" description:"Color output" choice:"auto" choice:"always" choice:"never" default:"auto"`
	Help                bool   `short:"h" long:"help" description:"Show this help"`
	Version             bool   `long:"version" description:"Show version"`
	URL                 string
	ExcludedPatterns    []*regexp.Regexp
	IncludePatterns     []*regexp.Regexp
	Header              http.Header
}

func getArguments(ss []string) (*arguments, error) {
	args := arguments{}
	ss, err := flags.NewParser(&args, flags.PassDoubleDash).ParseArgs(ss)

	if err != nil {
		return nil, err
	} else if args.Version || args.Help {
		return &args, nil
	} else if len(ss) != 1 {
		return nil, errors.New("invalid number of arguments")
	}

	reconcileDeprecatedArguments(&args)

	args.URL = ss[0]

	args.ExcludedPatterns, err = compileRegexps(args.RawExcludedPatterns)
	if err != nil {
		return nil, err
	}

	args.IncludePatterns, err = compileRegexps(args.RawIncludedPatterns)
	if err != nil {
		return nil, err
	}

	args.Header, err = parseHeaders(args.RawHeaders)
	if err != nil {
		return nil, err
	}

	if args.Format == "junit" && args.Verbose {
		return nil, errors.New("verbose option not supported for JUnit output")
	}

	return &args, nil
}

func help() string {
	p := flags.NewParser(&arguments{}, flags.PassDoubleDash)
	p.Usage = "[options] <url>"

	// Parse() is run here to show default values in help.
	// This seems to be a bug in go-flags.
	p.Parse() // nolint:errcheck

	b := &bytes.Buffer{}
	p.WriteHelp(b)
	return b.String()
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

func parseHeaders(headers []string) (http.Header, error) {
	m := make(http.Header, len(headers))

	for _, s := range headers {
		i := strings.IndexRune(s, ':')

		if i < 0 {
			return nil, errors.New("invalid header format")
		}

		m.Add(s[:i], strings.TrimSpace(s[i+1:]))
	}

	return m, nil
}

func reconcileDeprecatedArguments(args *arguments) {
	if args.JSONOutput {
		args.Format = "json"
		args.Verbose = args.Verbose || args.VerboseJSON
	} else if args.JUnitOutput {
		args.Format = "junit"
	}
}
