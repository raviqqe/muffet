package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/logrusorgru/aurora/v3"
)

type pageResultFormatter struct {
	verbose bool
	aurora  aurora.Aurora
}

func newPageResultFormatter(verbose bool, color bool) *pageResultFormatter {
	return &pageResultFormatter{verbose, aurora.NewAurora(color)}
}

func (f *pageResultFormatter) Format(r *pageResult) string {
	ss := []string(nil)

	if f.verbose {
		ss = append(ss, f.formatSuccessLinkResults(r.SuccessLinkResults)...)
	}

	ss = append(ss, f.formatErrorLinkResults(r.ErrorLinkResults)...)

	return strings.Join(
		append([]string{fmt.Sprint(f.aurora.Yellow(r.URL))}, formatMessages(ss)...),
		"\n",
	)
}

func (f *pageResultFormatter) formatSuccessLinkResults(rs []*successLinkResult) []string {
	ss := make([]string, 0, len(rs))

	for _, r := range rs {
		ss = append(ss, fmt.Sprintf("%v", f.aurora.Green(r.StatusCode))+"\t"+r.URL)
	}

	sort.Strings(ss)

	return ss
}

func (f *pageResultFormatter) formatErrorLinkResults(rs []*errorLinkResult) []string {
	ss := make([]string, 0, len(rs))

	for _, r := range rs {
		ss = append(ss, fmt.Sprintf("%v", f.aurora.Red(r.Error))+"\t"+r.URL)
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
