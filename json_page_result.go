package main

type jsonPageResult struct {
	URL   string `json:"url"`
	Links []any  `json:"links"`
}

type jsonSuccessLinkResult struct {
	URL    string `json:"url"`
	Status int    `json:"status"`
}

type jsonErrorLinkResult struct {
	URL   string `json:"url"`
	Error string `json:"error"`
}

func newJSONPageResult(r *pageResult, includeSuccess bool) *jsonPageResult {
	c := len(r.ErrorLinkResults)

	if includeSuccess {
		c += len(r.SuccessLinkResults)
	}

	ls := make([]any, 0, c)

	if includeSuccess {
		for _, r := range r.SuccessLinkResults {
			ls = append(ls, &jsonSuccessLinkResult{r.URL, r.StatusCode})
		}
	}

	for _, r := range r.ErrorLinkResults {
		ls = append(ls, &jsonErrorLinkResult{r.URL, r.Error.Error()})
	}

	return &jsonPageResult{r.URL, ls}
}
