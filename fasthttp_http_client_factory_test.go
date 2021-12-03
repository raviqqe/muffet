package main

import "testing"

func TestFasthttpHttpClientFactoryCreate(t *testing.T) {
	newFasthttpHttpClientFactory().Create(httpClientOptions{})
}

func TestFasthttpHttpClientFactoryCreateWithProxy(t *testing.T) {
	newFasthttpHttpClientFactory().Create(httpClientOptions{Proxy: "foo"})
}
