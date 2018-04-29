package main

type linkResult struct {
	statusCode int
	page       page
}

func newLinkResult(s int) linkResult {
	return linkResult{s, page{}}
}

func newLinkResultWithPage(s int, p page) linkResult {
	return linkResult{s, p}
}

func (r linkResult) StatusCode() int {
	return r.statusCode
}

func (r linkResult) Page() (page, bool) {
	return r.page, r.page != page{}
}
