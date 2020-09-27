package main

import "testing"

func TestNewThrottledHTTPClientHandleInvalidConnectionsOptionGracefully(t *testing.T) {
	newThrottledHTTPClient(nil, -1)
}
