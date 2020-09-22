package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDaemonManager(t *testing.T) {
	newDaemonManager(1)
}

func TestDaemonsRun(t *testing.T) {
	x := 0

	m := newDaemonManager(42)
	m.Add(func() { x++ })
	m.Run()

	assert.Equal(t, 1, x)
	assert.Zero(t, len(m.daemons))
}
