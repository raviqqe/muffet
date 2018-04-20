package main

import "strings"

// Result represents a summarized result of web page check.
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

// String is a method required for Stringer interface.
func (r Result) String() string {
	ss := make([]string, 0, len(r.successMessages))

	for _, s := range r.successMessages {
		ss = append(ss, "\tOK:\t"+s)
	}

	es := make([]string, 0, len(r.errorMessages))

	for _, e := range r.errorMessages {
		es = append(es, "\tERROR:\t"+e)
	}

	return strings.Join(append(append([]string{r.url}, ss...), es...), "\n")
}
