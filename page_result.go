package main

type pageResult struct {
	URL                string
	SuccessLinkResults []*successLinkResult
	ErrorLinkResults   []*errorLinkResult
}

type successLinkResult struct {
	URL        string
	StatusCode int
}

type errorLinkResult struct {
	URL   string
	Error error
}

func (r *pageResult) OK() bool {
	return len(r.ErrorLinkResults) == 0
}
