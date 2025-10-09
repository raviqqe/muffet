package main

import "testing"

func TestTlsHttpClientFactoryCreate(t *testing.T) {
	newTlsHttpClientFactory().Create(httpClientOptions{})
}
