package main

type jsonErrorPageResult struct {
	URL   string                 `json:"url"`
	Links []*jsonErrorLinkResult `json:"links"`
}

type jsonErrorLinkResult struct {
	URL   string `json:"url"`
	Error string `json:"error"`
}

type jsonSuccessPageResult struct {
	URL   string                   `json:"url"`
	Links []*jsonSuccessLinkResult `json:"links"`
}

type jsonSuccessLinkResult struct {
	URL    string `json:"url"`
	Status int    `json:"status"`
}

func newJSONErrorPageResult(r *pageResult) *jsonErrorPageResult {
	ls := make([]*jsonErrorLinkResult, 0, len(r.ErrorLinkResults))

	for _, r := range r.ErrorLinkResults {
		ls = append(ls, &jsonErrorLinkResult{r.URL, r.Error.Error()})
	}

	return &jsonErrorPageResult{r.URL, ls}
}

func newJSONSuccessPageResult(r *pageResult) *jsonSuccessPageResult {
	ls := make([]*jsonSuccessLinkResult, 0, len(r.SuccessLinkResults))

	for _, r := range r.SuccessLinkResults {
		ls = append(ls, &jsonSuccessLinkResult{r.URL, r.StatusCode})
	}

	return &jsonSuccessPageResult{r.URL, ls}
}
