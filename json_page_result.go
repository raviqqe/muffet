package main

type jsonPageResult struct {
	URL   string            `json:"url"`
	Links []*jsonLinkResult `json:"links"`
}

type jsonLinkResult struct {
	URL   string `json:"url"`
	Error string `json:"error"`
}

func newJSONPageResult(r *pageResult) *jsonPageResult {
	ls := make([]*jsonLinkResult, 0, len(r.ErrorLinkResults))

	for _, r := range r.ErrorLinkResults {
		ls = append(ls, &jsonLinkResult{r.URL, r.Error.Error()})
	}

	return &jsonPageResult{r.URL, ls}
}
