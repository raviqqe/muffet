package main

import (
	"sort"
	"strings"

	"github.com/fatih/color"
)

type pageResultFormatter struct {
	verbose bool
}

func newPageResultFormatter(verbose bool) *pageResultFormatter {
	return &pageResultFormatter{verbose}
}

func (f *pageResultFormatter) Format(r *pageResult) string {
	ss := []string(nil)

	if f.verbose {
		ss = append(ss, formatSuccessLinkResults(r.SuccessLinkResults)...)
	}

	ss = append(ss, formatErrorLinkResults(r.ErrorLinkResults)...)

	return strings.Join(
		append([]string{color.YellowString(r.URL)}, formatMessages(ss)...),
		"\n")
}

func formatSuccessLinkResults(rs []*successLinkResult) []string {
	ss := make([]string, 0, len(rs))

	for _, r := range rs {
		ss = append(ss, color.GreenString("%v", r.StatusCode)+"\t"+r.URL)
	}

	sort.Strings(ss)

	return ss
}

func formatErrorLinkResults(rs []*errorLinkResult) []string {
	ss := make([]string, 0, len(rs))

	for _, r := range rs {
		ss = append(ss, color.RedString("%v", r.Error)+"\t"+r.URL)
	}

	sort.Strings(ss)

	return ss
}

func formatMessages(ss []string) []string {
	ts := make([]string, 0, len(ss))

	for _, s := range ss {
		ts = append(ts, "\t"+s)
	}

	return ts
}
