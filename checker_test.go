package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChecker(t *testing.T) {
	_, err := newChecker(rootURL, 1, false, false)
	assert.Nil(t, err)
}

func TestNewCheckerError(t *testing.T) {
	_, err := newChecker(":", 1, false, false)
	assert.NotNil(t, err)
}

func TestNewCheckerWithMissingSitemap(t *testing.T) {
	_, err := newChecker("http://localhost:8081", 1, false, true)
	assert.NotNil(t, err)
}

func TestCheckerCheck(t *testing.T) {
	for _, s := range []string{rootURL, fragmentURL, baseURL, redirectURL} {
		c, _ := newChecker(s, 1, false, false)

		go c.Check()

		for r := range c.Results() {
			assert.True(t, r.OK())
		}
	}
}

func TestCheckerCheckMultiplePages(t *testing.T) {
	c, _ := newChecker(rootURL, 1, false, false)

	go c.Check()

	i := 0

	for r := range c.Results() {
		i += strings.Count(r.String(true), "\n") + 1
	}

	assert.Equal(t, 4, i)
}

func TestCheckerCheckPage(t *testing.T) {
	c, _ := newChecker(rootURL, 256, false, false)

	r, err := c.fetcher.FetchLink(existentURL)
	assert.Nil(t, err)

	p, ok := r.Page()
	assert.True(t, ok)

	go c.checkPage(p)

	assert.True(t, (<-c.Results()).OK())
}

func TestCheckerCheckPageError(t *testing.T) {
	for _, s := range []string{erroneousURL, invalidBaseURL} {
		c, _ := newChecker(rootURL, 256, false, false)

		r, err := c.fetcher.FetchLink(s)
		assert.Nil(t, err)

		p, ok := r.Page()
		assert.True(t, ok)

		go c.checkPage(p)

		assert.False(t, (<-c.Results()).OK())
	}
}

func TestStringChannelToSlice(t *testing.T) {
	for _, c := range []struct {
		channel chan string
		slice   []string
	}{
		{
			make(chan string, 1),
			[]string{},
		},
		{
			func() chan string {
				c := make(chan string, 1)
				c <- "foo"
				return c
			}(),
			[]string{"foo"},
		},
		{
			func() chan string {
				c := make(chan string, 2)
				c <- "foo"
				c <- "bar"
				return c
			}(),
			[]string{"foo", "bar"},
		},
		{
			func() chan string {
				c := make(chan string, 3)
				c <- "foo"
				c <- "bar"
				c <- "baz"
				return c
			}(),
			[]string{"foo", "bar", "baz"},
		},
	} {
		assert.Equal(t, c.slice, stringChannelToSlice(c.channel))
	}
}
