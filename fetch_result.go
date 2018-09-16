package main

import (
	"bytes"
	"encoding/gob"
)

type fetchResult struct {
	statusCode int
	page       *page
}

func newFetchResult(s int, p *page) fetchResult {
	return fetchResult{s, p}
}

func (r fetchResult) StatusCode() int {
	return r.statusCode
}

func (r fetchResult) Page() (*page, bool) {
	return r.page, r.page != nil
}

type encodableFetchResult struct {
	StatusCode int
	Page       *page
}

func (r *fetchResult) MarshalBinary() ([]byte, error) {
	b := bytes.NewBuffer(nil)

	if err := gob.NewEncoder(b).Encode(encodableFetchResult{r.statusCode, r.page}); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (r *fetchResult) UnmarshalBinary(bs []byte) error {
	rr := encodableFetchResult{}

	if err := gob.NewDecoder(bytes.NewBuffer(bs)).Decode(&rr); err != nil {
		return err
	}

	*r = fetchResult{rr.StatusCode, rr.Page}

	return nil
}
