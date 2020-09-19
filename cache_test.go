package main

import (
	"sync"
	"sync/atomic"
	"testing"

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

	l, s := int32(0), int32(0)
	g := &sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		g.Add(1)

		go func() {
			x, f := c.LoadOrStore("https://foo.com")

			if f == nil {
				assert.Equal(t, 42, x)
				atomic.AddInt32(&l, 1)
			} else {
				assert.Nil(t, x)
				f(42)
				atomic.AddInt32(&s, 1)
			}

			g.Done()
		}()
	}

	g.Wait()

	assert.Equal(t, int32(999), l)
	assert.Equal(t, int32(1), s)
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
