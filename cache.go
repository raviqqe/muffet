package main

import "sync"

type cache struct {
	locks  *sync.Map
	values *sync.Map
}

func newCache() cache {
	return cache{&sync.Map{}, &sync.Map{}}
}

func (c cache) LoadOrStore(key string) (any, func(any)) {
	if x, ok := c.values.Load(key); ok {
		return x, nil
	}

	g := &sync.WaitGroup{}
	g.Add(1)

	if g, ok := c.locks.LoadOrStore(key, g); ok {
		g.(*sync.WaitGroup).Wait()
		x, _ := c.values.Load(key)

		return x, nil
	}

	return nil, func(x any) {
		c.values.Store(key, x)
		g.Done()
	}
}
