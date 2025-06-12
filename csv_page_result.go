package main

import (
	"encoding/csv"
	"strconv"
	"strings"
)

type csvPageResult struct {
	URL   string
	Links []csvLinkResult
}

type csvLinkResult struct {
	URL    string
	Status string
}

func newCSVPageResult(r *pageResult, verbose bool) *csvPageResult {
	c := len(r.ErrorLinkResults)

	if verbose {
		c += len(r.SuccessLinkResults)
	}

	ls := make([]csvLinkResult, 0, c)

	if verbose {
		for _, r := range r.SuccessLinkResults {
			ls = append(ls, csvLinkResult{
				URL:    r.URL,
				Status: strconv.Itoa(r.StatusCode),
			})
		}
	}

	for _, r := range r.ErrorLinkResults {
		ls = append(ls, csvLinkResult{
			URL:    r.URL,
			Status: r.Error.Error(),
		})
	}

	return &csvPageResult{r.URL, ls}
}

func (r *csvPageResult) String() string {
	var buf strings.Builder
	writer := csv.NewWriter(&buf)

	// Write header
	writer.Write([]string{"Page URL", "Link URL", "Status"})

	// Write data rows
	for _, link := range r.Links {
		writer.Write([]string{r.URL, link.URL, link.Status})
	}

	writer.Flush()
	return buf.String()
}
