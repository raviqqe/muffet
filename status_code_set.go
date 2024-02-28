package main

import (
	"fmt"
	"strings"
)

type statusCodeSet map[statusCodeRange]struct{}

func parseStatusCodeSet(value string) (statusCodeSet, error) {
	rs := statusCodeSet{}

	for _, r := range strings.Split(value, ",") {
		if len(value) == 0 {
			return nil, fmt.Errorf("invalid status code range: %s", value)
		}

		r, err := parseStatusCodeRange(r)
		if err != nil {
			return nil, err
		}

		rs[*r] = struct{}{}
	}

	return rs, nil
}

func (c statusCodeSet) isInSet(code int) bool {
	for r := range c {
		if r.isInRange(code) {
			return true
		}
	}

	return false
}
