package main

import (
	"bytes"
	"errors"
	"regexp"
	"strings"

	"github.com/jessevdk/go-flags"
)

type arguments struct {
	BufferSize            int      `short:"b" long:"buffer-size" value-name:"<size>" default:"4096" description:"HTTP response buffer size in bytes"`
	MaxConnections        int      `short:"c" long:"max-connections" value-name:"<count>" default:"512" description:"Maximum number of HTTP connections"`
	MaxConnectionsPerHost int      `long:"max-connections-per-host" value-name:"<count>" default:"512" description:"Maximum number of HTTP connections per host"`
	RawExcludedPatterns   []string `short:"e" long:"exclude" value-name:"<pattern>..." description:"Exclude URLs matched with given regular expressions"`
	FollowRobotsTxt       bool     `long:"follow-robots-txt" description:"Follow robots.txt when scraping pages"`
	FollowSitemapXML      bool     `long:"follow-sitemap-xml" description:"Scrape only pages listed in sitemap.xml"`
	RawHeaders            []string `long:"header" value-name:"<header>..." description:"Custom headers"`
	IgnoreFragments       bool     `short:"f" long:"ignore-fragments" description:"Ignore URL fragments"`
	JSONOutput            bool     `long:"json" description:"Output results in JSON"`
	MaxRedirections       int      `short:"r" long:"max-redirections" value-name:"<count>" default:"64" description:"Maximum number of redirections"`
	Timeout               int      `short:"t" long:"timeout" value-name:"<seconds>" default:"10" description:"Timeout for HTTP requests in seconds"`
	Verbose               bool     `short:"v" long:"verbose" description:"Show successful results too"`
	SkipTLSVerification   bool     `long:"skip-tls-verification" description:"Skip TLS certificate verification"`
	OnePageOnly           bool     `long:"one-page-only" description:"Only check links found in the given URL"`
	Color                 color    `long:"color" description:"Color output." choice:"auto" choice:"always" choice:"never" default:"auto"`
	Help                  bool     `short:"h" long:"help" description:"Show this help"`
	Version               bool     `long:"version" description:"Show version"`
	URL                   string
	ExcludedPatterns      []*regexp.Regexp
	Headers               map[string]string
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

	args.URL = ss[0]

	rs, err := compileRegexps(args.RawExcludedPatterns)
	if err != nil {
		return nil, err
	}

	args.ExcludedPatterns = rs

	hs, err := parseHeaders(args.RawHeaders)
	if err != nil {
		return nil, err
	}

	args.Headers = hs

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
