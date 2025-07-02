package main

type ignoreNetworkErrors string

const (
	ignoreNetworkErrorsNone     ignoreNetworkErrors = "none"
	ignoreNetworkErrorsAll      ignoreNetworkErrors = "all"
	ignoreNetworkErrorsExternal ignoreNetworkErrors = "external"
)
