package main

import "sync"

type concurrentStringSet struct {
	set *sync.Map
}

func newConcurrentStringSet() concurrentStringSet {
	return concurrentStringSet{&sync.Map{}}
}

func (c concurrentStringSet) Add(s string) bool {
	_, exist := c.set.LoadOrStore(s, nil)
	return exist
}
