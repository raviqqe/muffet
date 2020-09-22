package main

import "net/url"

type fakeHTTPClientFactory struct {
	handler func(*url.URL) (*fakeHTTPResponse, error)
}

func newFakeHTTPClientFactory(h func(*url.URL) (*fakeHTTPResponse, error)) httpClientFactory {
	return &fakeHTTPClientFactory{h}
}

func (f *fakeHTTPClientFactory) Create(httpClientOptions) httpClient {
	return newFakeHTTPClient(f.handler)
}
