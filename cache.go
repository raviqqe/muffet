package main

import "sync"

type cache struct {
	locks  *sync.Map
	values *sync.Map
}

func newCache() cache {
	return cache{&sync.Map{}, &sync.Map{}}
}

func (c cache) LoadOrStore(s string) (interface{}, func(interface{}), bool) {
	g := &sync.WaitGroup{}
	g.Add(1)

	if g, ok := c.locks.LoadOrStore(s, g); ok {
		g.(*sync.WaitGroup).Wait()
		x, _ := c.values.Load(s)

		return x, nil, true
	}

	return nil, func(x interface{}) {
		c.values.Store(s, x)
		g.Done()
	}, false
}
