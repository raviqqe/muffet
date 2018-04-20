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

// String is a method required for Stringer interface.
func (r Result) String() string {
	ss := make([]string, 0, len(r.successMessages))

	for _, s := range r.successMessages {
		ss = append(ss, formatMessage(color.GreenString("OK"), s))
	}

	es := make([]string, 0, len(r.errorMessages))

	for _, s := range r.errorMessages {
		es = append(es, formatMessage(color.RedString("ERROR"), s))
	}

	return strings.Join(append(append([]string{color.YellowString(r.url)}, ss...), es...), "\n")
}

func formatMessage(s, t string) string {
	return strings.Join([]string{"\t", s, "\t", t}, "")
}
