package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFetcher(t *testing.T) {
	newFetcher(1)
}

type handler struct{}

func (handler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("hello"))
}

func TestFetcherCache(t *testing.T) {
	go http.ListenAndServe(":8080", handler{})
	time.Sleep(time.Millisecond)

	f := newFetcher(1)

	p, err := f.Fetch("http://localhost:8080")

	assert.NotEqual(t, (*page)(nil), p)
	assert.Nil(t, err)

	p, err = f.Fetch("https://foo.com/non/existent/path")

	assert.Equal(t, (*page)(nil), p)
	assert.NotNil(t, err)

	p, err = f.Fetch("http://localhost:8080")

	assert.Equal(t, (*page)(nil), p)
	assert.Nil(t, err)

	p, err = f.Fetch("https://foo.com/non/existent/path")

	assert.Equal(t, (*page)(nil), p)
	assert.NotNil(t, err)
}

func TestFetcherFetchError(t *testing.T) {
	f := newFetcher(1)

	for _, s := range []string{"https://foo.com/non/existent/path", ":"} {
		p, err := f.Fetch(s)

		assert.Equal(t, (*page)(nil), p)
		assert.NotNil(t, err)
	}
}
