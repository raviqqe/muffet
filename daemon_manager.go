package main

import "sync"

type daemonManager struct {
	daemons   chan func()
	waitGroup *sync.WaitGroup
}

func newDaemonManager(capacity int) daemonManager {
	return daemonManager{make(chan func(), capacity), &sync.WaitGroup{}}
}

func (m daemonManager) Add(f func()) {
	m.waitGroup.Add(1)

	m.daemons <- func() {
		f()
		m.waitGroup.Done()
	}
}

func (m daemonManager) Run() {
	go func() {
		for f := range m.daemons {
			go f()
		}
	}()

	m.waitGroup.Wait()
}
