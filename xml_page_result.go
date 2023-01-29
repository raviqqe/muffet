package main

type xmlPageResult struct {
	Url      string `xml:"name,attr"`
	Total    int    `xml:"tests,attr"`
	Failures int    `xml:"failures,attr"`
	Skipped  int    `xml:"skipped,attr"`
	// spell-checker: disable-next-line
	Links []*xmlLinkResult `xml:"testcase"`
}

type xmlLinkResult struct {
	Url string `xml:"name,attr"`
	// spell-checker: disable-next-line
	Source  string          `xml:"classname,attr"`
	Failure *xmlLinkFailure `xml:"failure"`
}

type xmlLinkFailure struct {
	Message string `xml:"message,attr"`
}

func newXMLPageResult(pr *pageResult) *xmlPageResult {
	ls := make([]*xmlLinkResult, 0, len(pr.SuccessLinkResults)+len(pr.ErrorLinkResults))

	for _, r := range pr.SuccessLinkResults {
		ls = append(
			ls,
			&xmlLinkResult{
				Url:    r.URL,
				Source: pr.URL,
			},
		)
	}

	for _, r := range pr.ErrorLinkResults {
		ls = append(
			ls,
			&xmlLinkResult{
				Url:     r.URL,
				Source:  pr.URL,
				Failure: &xmlLinkFailure{Message: r.Error.Error()},
			},
		)
	}

	return &xmlPageResult{
		Url: pr.URL,
		// TODO: Consider adding information skipped links, if that can be tracked.
		Skipped:  0,
		Total:    len(ls),
		Failures: len(pr.ErrorLinkResults),
		Links:    ls,
	}
}
