package main

import "strings"

type statusCodeSet map[statusCodeRange]struct{}

func parseStatusCodeSet(value string) (statusCodeSet, error) {
	rs := statusCodeSet{}

	for _, r := range strings.Split(value, ",") {
		r, err := parseStatusCodeRange(r)
		if err != nil {
			return nil, err
		}

		rs[*r] = struct{}{}
	}

	return rs, nil
}

func (s statusCodeSet) Contains(code int) bool {
	for r := range s {
		if r.Contains(code) {
			return true
		}
	}

	return false
}
