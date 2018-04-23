package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const existentURL = "http://localhost:8080"
const nonExistentURL = "http://localhost:8080/hello"

type handler struct{}

func (handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "", "/":
		w.Write([]byte("hello"))
	default:
		w.WriteHeader(404)
	}
}

func TestNewFetcher(t *testing.T) {
	newFetcher(1)
}

func TestFetcherCache(t *testing.T) {
	go http.ListenAndServe(":8080", handler{})
	time.Sleep(time.Millisecond)

	f := newFetcher(1)

	p, err := f.Fetch(existentURL)

	assert.NotEqual(t, (*page)(nil), p)
	assert.Nil(t, err)

	p, err = f.Fetch(nonExistentURL)

	assert.Equal(t, (*page)(nil), p)
	assert.NotNil(t, err)

	p, err = f.Fetch(existentURL)

	assert.Equal(t, (*page)(nil), p)
	assert.Nil(t, err)

	p, err = f.Fetch(nonExistentURL)

	assert.Equal(t, (*page)(nil), p)
	assert.NotNil(t, err)
}

func TestFetcherFetchError(t *testing.T) {
	go http.ListenAndServe(":8080", handler{})
	time.Sleep(time.Millisecond)

	f := newFetcher(1)

	for _, s := range []string{nonExistentURL, ":"} {
		p, err := f.Fetch(s)

		assert.Equal(t, (*page)(nil), p)
		assert.NotNil(t, err)
	}
}
