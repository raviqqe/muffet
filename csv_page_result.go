package main

import (
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

	buf.WriteString(`"Page URL","Link URL",Status` + "\n")

	for _, link := range r.Links {
		buf.WriteString(`"` + strings.ReplaceAll(r.URL, `"`, `""`) + `","` +
			strings.ReplaceAll(link.URL, `"`, `""`) + `",` +
			link.Status + "\n")
	}

	return buf.String()
}
