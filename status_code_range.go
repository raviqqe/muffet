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

	start, err := strconv.Atoi(ss[0])
	if err != nil {
		return nil, fmt.Errorf("invalid status code: %v", ss[0])
	}

	end, err := strconv.Atoi(ss[1])
	if err != nil {
		return nil, fmt.Errorf("invalid status code: %v", ss[1])
	}

	return &statusCodeRange{start, end}, nil
}

func (r statusCodeRange) isInRange(code int) bool {
	return code >= r.start && code < r.end
}
