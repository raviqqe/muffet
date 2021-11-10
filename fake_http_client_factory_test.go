package main

import "net/url"

type fakeHttpClientFactory struct {
	handler func(*url.URL) (*fakeHttpResponse, error)
}

func newFakeHttpClientFactory(h func(*url.URL) (*fakeHttpResponse, error)) httpClientFactory {
	return &fakeHttpClientFactory{h}
}

func (f *fakeHttpClientFactory) Create(httpClientOptions) httpClient {
	return newFakeHttpClient(f.handler)
}
