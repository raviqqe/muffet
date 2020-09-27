package main

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCache(t *testing.T) {
	newCache()
}

func TestCacheLoadOrStore(t *testing.T) {
	c := newCache()

	x, f := c.LoadOrStore("https://foo.com")

	assert.Nil(t, x)
	assert.NotNil(t, f)

	f(42)

	x, f = c.LoadOrStore("https://foo.com")

	assert.Equal(t, 42, x)
	assert.Nil(t, f)
}

func TestCacheLoadOrStoreConcurrency(t *testing.T) {
	c := newCache()

	x, f := c.LoadOrStore("key")
	assert.Nil(t, x)
	assert.NotNil(t, f)

	go func() {
		x, f := c.LoadOrStore("key")
		assert.Equal(t, 42, x)
		assert.Nil(t, f)
	}()

	time.Sleep(time.Millisecond)

	f(42)
}

func BenchmarkCacheLoadOrStore(b *testing.B) {
	c := newCache()
	g := &sync.WaitGroup{}

	_, f := c.LoadOrStore("https://foo.com")
	f(42)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.Add(1)

		go func() {
			c.LoadOrStore("https://foo.com")
			g.Done()
		}()
	}

	g.Wait()
}
