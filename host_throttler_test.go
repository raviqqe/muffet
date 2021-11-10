package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHostThrottler(t *testing.T) {
	newHostThrottler(1, 1)
}

func TestHostThrottlerRequest(t *testing.T) {
	c := make(chan struct{}, 100)
	s := newHostThrottler(1000000, 1)

	for i := 0; i < 2; i++ {
		go func() {
			s.Request()
			c <- struct{}{}
		}()
	}

	<-c

	assert.Equal(t, 0, len(c))

	s.Release()
	<-c
}
