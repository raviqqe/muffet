package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHostThrottlerPool(t *testing.T) {
	newHostThrottlerPool(1, 1)
}

func TestHostThrottlerPoolGet(t *testing.T) {
	c := make(chan struct{}, 100)
	s := newHostThrottlerPool(1000000, 1)

	for i := 0; i < 2; i++ {
		go func() {
			s.Get("foo").Request()
			c <- struct{}{}
		}()
	}

	<-c

	assert.Equal(t, 0, len(c))

	s.Get("foo").Release()
	<-c
}
