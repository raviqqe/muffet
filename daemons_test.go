package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDaemons(t *testing.T) {
	newDaemons(1)
}

func TestDaemonsAdd(t *testing.T) {
	ds := newDaemons(42)
	ds.Add(func() {})

	assert.Equal(t, 1, len(ds.daemons))
}

func TestDaemonsRun(t *testing.T) {
	x := 0

	ds := newDaemons(42)
	ds.Add(func() { x++ })
	ds.Run()

	assert.Equal(t, 1, x)
	assert.Zero(t, len(ds.daemons))
}
