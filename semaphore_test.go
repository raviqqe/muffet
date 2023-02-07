package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSemaphoreRequest(t *testing.T) {
	s := newSemaphore(1)

	s.Request()

	assert.Equal(t, 1, len(s.channel))
}

func TestSemaphoreRelease(t *testing.T) {
	s := newSemaphore(1)

	s.Request()
	s.Release()

	assert.Equal(t, 0, len(s.channel))
}
