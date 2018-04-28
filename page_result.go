package main

import (
	"sort"
	"strings"

	"github.com/fatih/color"
)

type pageResult struct {
	url                            string
	successMessages, errorMessages []string
}

func newPageResult(u string, ss, es []string) pageResult {
	return pageResult{u, ss, es}
}

func (r pageResult) OK() bool {
	return len(r.errorMessages) == 0
}

func (r pageResult) String(v bool) string {
	ss := []string(nil)

	if v {
		ss = formatMessages(color.GreenString("OK"), r.successMessages)
	}

	return strings.Join(
		append(append([]string{color.YellowString(r.url)},
			ss...),
			formatMessages(color.RedString("ERROR"), r.errorMessages)...),
		"\n")
}

func formatMessages(prefix string, ss []string) []string {
	ts := make([]string, 0, len(ss))

	for _, s := range ss {
		ts = append(ts, strings.Join([]string{"\t", prefix, "\t", s}, ""))
	}

	sort.Strings(ts)

	return ts
}
