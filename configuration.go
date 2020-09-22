package main

import "time"

const (
	version                = "1.5.3"
	defaultBufferSize      = 4096
	defaultConcurrency     = 512
	defaultMaxRedirections = 64
	defaultTimeout         = 10 * time.Second
	agentName              = "muffet"
	tcpTimeout             = time.Minute
)
