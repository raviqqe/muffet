package main

import "time"

const (
	version           = "2.11.1"
	agentName         = "muffet"
	concurrency       = 1024
	tcpTimeout        = 5 * time.Second
	initialRetryDelay = 500 * time.Millisecond
	maxRetryDelay     = 10 * time.Second
	retryBackoff      = 2
)
