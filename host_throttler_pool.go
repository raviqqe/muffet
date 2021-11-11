package main

import "sync"

type hostThrottlerPool struct {
	requestPerSecond, maxConnectionsPerHost int
	hostMap                                 sync.Map
}

func newHostThrottlerPool(requestPerSecond, maxConnectionsPerHost int) *hostThrottlerPool {
	return &hostThrottlerPool{requestPerSecond, maxConnectionsPerHost, sync.Map{}}
}

func (p *hostThrottlerPool) Get(name string) *hostThrottler {
	t := newHostThrottler(p.requestPerSecond, p.maxConnectionsPerHost)
	x, ok := p.hostMap.LoadOrStore(name, t)

	if ok {
		t = x.(*hostThrottler)
	}

	return t
}
