package main

import (
	"fmt"
	"strconv"
	"strings"
)

type statusCodeRange struct {
	start int
	end   int
}

func parseStatusCodeRange(s string) (*statusCodeRange, error) {
	if c, err := strconv.Atoi(s); err == nil {
		return &statusCodeRange{c, c + 1}, nil
	}

	ss := strings.Split(s, "..")
	if len(ss) != 2 {
		return nil, fmt.Errorf("invalid status code range: %v", s)
	}

	cs := []int{0, 0}

	for i, s := range ss {
		c, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid status code: %v", s)
		}

		cs[i] = c
	}

	return &statusCodeRange{c[0], c[1]}, nil
}

func (r statusCodeRange) isInRange(code int) bool {
	return code >= r.start && code < r.end
}
