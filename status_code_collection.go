package main

import "strings"

type statusCodeCollection struct {
	elements []statusCodeRange
}

func parseStatusCodeCollection(value string) (statusCodeCollection, error) {
	statusCodeRanges := []statusCodeRange{}

	for _, partial := range strings.Split(value, ",") {
		if len(value) == 0 {
			continue
		}

		statusCodeRange, err := parseStatusCodeRange(partial)

		if err != nil {
			return statusCodeCollection{}, err
		}

		statusCodeRanges = append(statusCodeRanges, *statusCodeRange)
	}

	if len(statusCodeRanges) == 0 {
		statusCodeRanges = append(statusCodeRanges, statusCodeRange{200, 300})
	}

	return statusCodeCollection{statusCodeRanges}, nil
}

func (c *statusCodeCollection) isInCollection(code int) bool {
	for _, element := range c.elements {
		if element.isInRange(code) {
			return true
		}
	}

	return false
}
