package main

type networkError int

const (
	networkErrorNone networkError = iota
	networkErrorAll
	networkErrorExternal
)
