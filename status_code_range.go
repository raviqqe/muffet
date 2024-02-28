package main

import (
	"errors"
	"regexp"
	"strconv"
)

var fixedCodePattern = regexp.MustCompile(`^\s*(\d{3})\s*$`)
var rangeCodePattern = regexp.MustCompile(`^\s*(\d{3})\s*\.\.\s*(\d{3})\s*$`)

type statusCodeRange struct {
	start int
	end   int
}

func parseStatusCodeRange(value string) (*statusCodeRange, error) {
	fixedMatch := fixedCodePattern.FindAllStringSubmatch(value, -1)
	if len(fixedMatch) > 0 {
		code, _ := strconv.Atoi(fixedMatch[0][1])
		return &statusCodeRange{code, code + 1}, nil
	}

	rangeMatch := rangeCodePattern.FindAllStringSubmatch(value, -1)
	if len(rangeMatch) > 0 {
		start, _ := strconv.Atoi(rangeMatch[0][1])
		end, _ := strconv.Atoi(rangeMatch[0][2])
		return &statusCodeRange{start, end}, nil
	}

	return nil, errors.New("invalid HTTP response status code value")
}

func (r *statusCodeRange) isInRange(code int) bool {
	return code >= r.start && code < r.end
}
