package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConcurrentStringSet(t *testing.T) {
	newConcurrentStringSet()
}

func TestConcurrentStringSetAdd(t *testing.T) {
	s := newConcurrentStringSet()
	assert.False(t, s.Add("foo"))
	assert.True(t, s.Add("foo"))
}
