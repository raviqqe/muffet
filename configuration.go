package main

import "time"

const (
	version           = "2.10.9"
	agentName         = "muffet"
	concurrency       = 1024
	tcpTimeout        = 5 * time.Second
	initialRetryDelay = 500 * time.Millisecond
)
