package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCache(t *testing.T) {
	d := createCacheDirectory(t)
	defer os.Remove(d)

	_, err := newCache(d)
	assert.Nil(t, err)

	os.Remove(d)
}

func TestCacheAddAndGet(t *testing.T) {
	d := createCacheDirectory(t)
	defer os.Remove(d)

	c, err := newCache(d)
	assert.Nil(t, err)

	err = c.Add(rootURL, newFetchResult(200))
	assert.Nil(t, err)

	time.Sleep(time.Second)

	i, err := c.Get(rootURL)
	assert.Nil(t, err)
	assert.Equal(t, newFetchResult(200), i)
}

func createCacheDirectory(t *testing.T) string {
	d, err := ioutil.TempDir("", "")

	if err != nil {
		t.FailNow()
	}

	return d
}
