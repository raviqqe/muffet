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
		return err == nil, err
	} else if args.JSONOutput && args.Verbose {
		return false, errors.New("verbose option not supported for JSON output")
	} else if args.JUnitOutput && args.Verbose {
		return false, errors.New("verbose option not supported for JUnit output")
	} else if args.JSONOutput && args.JUnitOutput {
		return false, errors.New("JSON and JUnit output are mutually exclusive")
	}

	client := newRedirectHttpClient(
		newThrottledHttpClient(
			c.httpClientFactory.Create(
				httpClientOptions{
					MaxConnectionsPerHost: args.MaxConnectionsPerHost,
					MaxResponseBodySize:   args.MaxResponseBodySize,
					BufferSize:            args.BufferSize,
					Proxy:                 args.Proxy,
					SkipTLSVerification:   args.SkipTLSVerification,
					Timeout:               time.Duration(args.Timeout) * time.Second,
					Headers:               args.Headers,
				},
			),
			args.RateLimit,
			args.MaxConnections,
			args.MaxConnectionsPerHost,
		),
		args.MaxRedirections,
	)

	pp := newPageParser(newLinkFinder(args.ExcludedPatterns, args.IncludePatterns))

	f := newLinkFetcher(
		client,
		pp,
		linkFetcherOptions{
			args.IgnoreFragments,
		},
	)

	_, p, err := f.Fetch(args.URL)
	if err != nil {
		return false, fmt.Errorf("failed to fetch root page: %v", err)
	} else if p == nil {
		return false, errors.New("root page is not HTML")
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

	if args.JSONOutput {
		return c.printResultsInJSON(checker.Results(), args.IncludeSuccessInJSONOutput)
	} else if args.JUnitOutput {
		return c.printResultsAsJUnitXML(checker.Results())
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

		if !r.OK() {
			ok = false
		}
	}

	return ok, nil
}

func (c *command) printResultsInJSON(rc <-chan *pageResult, includeSuccess bool) (bool, error) {
	rs := []interface{}{}
	ok := true

	for r := range rc {
		if r.OK() && includeSuccess {
			rs = append(rs, newJSONSuccessPageResult(r))
		} else if !r.OK() {
			rs = append(rs, newJSONErrorPageResult(r))
			ok = false
		}
	}

	bs, err := json.Marshal(rs)

	if err != nil {
		return false, err
	}

	c.print(string(bs))

	return ok, nil
}

func (c *command) printResultsAsJUnitXML(rc <-chan *pageResult) (bool, error) {
	rs := []*xmlPageResult{}
	ok := true

	for r := range rc {
		rs = append(rs, newXMLPageResult(r))

		if !r.OK() {
			ok = false
		}
	}

	results := &struct {
		// spell-checker: disable-next-line
		XMLName xml.Name `xml:"testsuites"`
		// spell-checker: disable-next-line
		PageResults []*xmlPageResult `xml:"testsuite"`
	}{PageResults: rs}

	data, err := xml.MarshalIndent(results, "", "  ")

	if err != nil {
		return ok, err
	}

	c.print(xml.Header)
	c.print(string(data))

	return ok, nil
}

func (c *command) print(xs ...interface{}) {
	if _, err := fmt.Fprintln(c.stdout, strings.TrimSpace(fmt.Sprint(xs...))); err != nil {
		panic(err)
	}
}

func (c *command) printError(xs ...interface{}) {
	s := fmt.Sprint(xs...)

	// Do not check --color option here because this can be used on argument parsing errors.
	if c.terminal {
		s = aurora.Red(s).String()
	}

	if _, err := fmt.Fprintln(c.stderr, s); err != nil {
		panic(err)
	}
}
