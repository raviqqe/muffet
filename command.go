package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/logrusorgru/aurora/v3"
	"github.com/temoto/robotstxt"
)

type command struct {
	stdout, stderr    io.Writer
	terminal          bool
	httpClientFactory httpClientFactory
}

func newCommand(stdout, stderr io.Writer, terminal bool, f httpClientFactory) *command {
	return &command{stdout, stderr, terminal, f}
}

func (c *command) Run(args []string) bool {
	ok, err := c.runWithError(args)
	if err != nil {
		c.printError(err)
	}

	return ok
}

func (c *command) runWithError(ss []string) (bool, error) {
	args, err := getArguments(ss)
	if err != nil {
		return false, err
	} else if args.Help {
		c.print(help())
		return true, nil
	} else if args.Version {
		c.print(version)
		return true, nil
	}

	client := newCheckedHttpClient(
		newRedirectHttpClient(
			newThrottledHttpClient(
				c.httpClientFactory.Create(
					httpClientOptions{
						MaxConnectionsPerHost: args.MaxConnectionsPerHost,
						MaxResponseBodySize:   args.MaxResponseBodySize,
						BufferSize:            args.BufferSize,
						Proxy:                 args.Proxy,
						SkipTLSVerification:   args.SkipTLSVerification,
						Timeout:               time.Duration(args.Timeout) * time.Second,
						Header:                args.Header,
					},
				),
				args.RateLimit,
				args.MaxConnections,
				args.MaxConnectionsPerHost,
			),
			args.MaxRedirections,
		),
		args.AcceptedStatusCodes,
	)

	fl := newLinkFilterer(args.ExcludedPatterns, args.IncludePatterns)

	f := newLinkFetcher(
		client,
		[]pageParser{
			newSitemapPageParser(fl),
			newHtmlPageParser(newLinkFinder(fl)),
		},
		linkFetcherOptions{
			args.IgnoreFragments,
		},
	)

	_, p, err := f.Fetch(args.URL)
	if err != nil {
		return false, fmt.Errorf("failed to fetch root page: %v", err)
	} else if p == nil {
		return false, errors.New("root page has invalid content type")
	}

	rd := (*robotstxt.RobotsData)(nil)

	if args.FollowRobotsTxt {
		rd, err = newRobotsTxtFetcher(client).Fetch(p.URL())

		if err != nil {
			return false, err
		}
	}

	sm := (map[string]struct{})(nil)

	if args.FollowSitemapXML {
		sm, err = newSitemapFetcher(client).Fetch(p.URL())

		if err != nil {
			return false, err
		}
	}

	checker := newPageChecker(
		f,
		newLinkValidator(p.URL().Hostname(), rd, sm),
		args.OnePageOnly,
	)

	go checker.Check(p)

	switch args.Format {
	case "json":
		return c.printResultsInJSON(checker.Results(), args.Verbose)
	case "junit":
		return c.printResultsInJUnitXML(checker.Results())
	}

	formatter := newPageResultFormatter(
		args.Verbose,
		isColorEnabled(args.Color, c.terminal),
	)

	ok := true

	for r := range checker.Results() {
		if !r.OK() || args.Verbose {
			c.print(formatter.Format(r))
		}

		ok = ok && r.OK()
	}

	return ok, nil
}

func (c *command) printResultsInJSON(rc <-chan *pageResult, verbose bool) (bool, error) {
	rs := []any{}
	ok := true

	for r := range rc {
		if !r.OK() || verbose {
			rs = append(rs, newJSONPageResult(r, verbose))
		}

		ok = ok && r.OK()
	}

	bs, err := json.Marshal(rs)

	if err != nil {
		return false, err
	}

	c.print(string(bs))

	return ok, nil
}

func (c *command) printResultsInJUnitXML(rc <-chan *pageResult) (bool, error) {
	rs := []*xmlPageResult{}
	ok := true

	for r := range rc {
		rs = append(rs, newXMLPageResult(r))
		ok = ok && r.OK()
	}

	bs, err := xml.MarshalIndent(
		struct {
			// spell-checker: disable-next-line
			XMLName xml.Name `xml:"testsuites"`
			// spell-checker: disable-next-line
			PageResults []*xmlPageResult `xml:"testsuite"`
		}{
			PageResults: rs,
		},
		"",
		"  ",
	)

	if err != nil {
		return false, err
	}

	c.print(xml.Header)
	c.print(string(bs))

	return ok, nil
}

func (c *command) print(xs ...any) {
	if _, err := fmt.Fprintln(c.stdout, strings.TrimSpace(fmt.Sprint(xs...))); err != nil {
		panic(err)
	}
}

func (c *command) printError(xs ...any) {
	s := fmt.Sprint(xs...)

	// Do not check --color option here because this can be used on argument parsing errors.
	if c.terminal {
		s = aurora.Red(s).String()
	}

	if _, err := fmt.Fprintln(c.stderr, s); err != nil {
		panic(err)
	}
}
