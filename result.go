package main

import (
	"strings"

	"github.com/fatih/color"
)

type result struct {
	url                            string
	successMessages, errorMessages []string
}

func newResult(u string, ss, es []string) result {
	return result{u, ss, es}
}

func newResultWithError(u string, err error) result {
	return newResult(u, nil, []string{err.Error()})
}

func (r result) OK() bool {
	return len(r.errorMessages) == 0
}

func (r result) String(v bool) string {
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

	return ts
}
