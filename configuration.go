package main

import "time"

const (
	version               = "2.0.1"
	agentName             = "muffet"
	concurrency           = 1024
	tcpTimeout            = time.Minute
	defaultMaxConnections = 512
)
