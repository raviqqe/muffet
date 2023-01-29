package main

type xmlPageResult struct {
	Url      string           `xml:"name,attr"`
	Total    int              `xml:"tests,attr"`
	Failures int              `xml:"failures,attr"`
	Skipped  int              `xml:"skipped,attr"`
	Time     float64          `xml:"time,attr"`
	Links    []*xmlLinkResult `xml:"testcase"`
}

type xmlLinkResult struct {
	Url     string          `xml:"name,attr"`
	Time    float64         `xml:"time,attr"`
	Source  string          `xml:"classname,attr"`
	Failure *xmlLinkFailure `xml:"failure"`
}

type xmlLinkFailure struct {
	Message string `xml:"message,attr"`
}

func newXMLPageResult(pr *pageResult) *xmlPageResult {
	ls := make([]*xmlLinkResult, 0, len(pr.ErrorLinkResults)+len(pr.SuccessLinkResults))

	// TODO: Consider adding information skipped links, if that can be
	// tracked.
	for _, r := range pr.ErrorLinkResults {
		failure := &xmlLinkFailure{Message: r.Error.Error()}
		l := &xmlLinkResult{Url: r.URL, Source: pr.URL, Time: r.Elapsed.Seconds(), Failure: failure}
		ls = append(ls, l)
	}

	for _, r := range pr.SuccessLinkResults {
		l := &xmlLinkResult{Url: r.URL, Source: pr.URL, Time: r.Elapsed.Seconds()}
		ls = append(ls, l)
	}

	return &xmlPageResult{Url: pr.URL, Time: pr.Elapsed.Seconds(), Skipped: 0, Total: len(ls), Failures: len(pr.ErrorLinkResults), Links: ls}
}
