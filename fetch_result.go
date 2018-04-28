package main

type fetchResult struct {
	statusCode int
	page       page
}

func newFetchResult(s int) fetchResult {
	return fetchResult{s, page{}}
}

func newFetchResultWithPage(s int, p page) fetchResult {
	return fetchResult{s, p}
}

func (r fetchResult) StatusCode() int {
	return r.statusCode
}

func (r fetchResult) Page() (page, bool) {
	return r.page, r.page != page{}
}
