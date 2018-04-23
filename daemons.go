package main

import "sync"

type daemons struct {
	daemons   chan func()
	waitGroup *sync.WaitGroup
}

func newDaemons(c int) daemons {
	return daemons{make(chan func(), c), &sync.WaitGroup{}}
}

func (ds daemons) Add(f func()) {
	ds.waitGroup.Add(1)

	ds.daemons <- func() {
		f()
		ds.waitGroup.Done()
	}
}

func (ds daemons) Run() {
	go func() {
		for f := range ds.daemons {
			go f()
		}
	}()

	ds.waitGroup.Wait()
}
