package main

import (
	"strings"

	"github.com/fatih/color"
)

// Result represents a summarized result of a web page check.
type Result struct {
	url                            string
	successMessages, errorMessages []string
}

// NewResult creates a new result.
func NewResult(u string, ss, es []string) Result {
	return Result{u, ss, es}
}

// NewResultWithError creates a new result with a single error.
func NewResultWithError(u string, err error) Result {
	return NewResult(u, nil, []string{err.Error()})
}

// IsError returns true if a result contains some errors and false otherwise.
func (r Result) IsError() bool {
	return len(r.errorMessages) != 0
}

// String turns a result into informational string.
func (r Result) String(v bool) string {
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
