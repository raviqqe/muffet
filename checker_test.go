package main

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChecker(t *testing.T) {
	_, err := newChecker(rootURL, checkerOptions{})
	assert.Nil(t, err)
}

func TestNewCheckerError(t *testing.T) {
	for _, s := range []string{":", invalidBaseURL} {
		_, err := newChecker(s, checkerOptions{})
		assert.NotNil(t, err)
	}
}

func TestNewCheckerWithNonHTMLPage(t *testing.T) {
	_, err := newChecker(robotsTxtURL, checkerOptions{})
	assert.Equal(t, "non-HTML page", err.Error())
}

func TestNewCheckerWithMissingSitemapXML(t *testing.T) {
	_, err := newChecker(missingMetadataURL, checkerOptions{FollowSitemapXML: true})
	assert.Equal(t, "sitemap not found", err.Error())
}

func TestCheckerCheck(t *testing.T) {
	for _, s := range []string{rootURL, fragmentURL, baseURL, redirectURL} {
		c, err := newChecker(s, checkerOptions{})
		assert.Nil(t, err)

		go c.Check()

		for r := range c.Results() {
			assert.True(t, r.OK())
		}
	}
}

func TestCheckerCheckMultiplePages(t *testing.T) {
	c, _ := newChecker(rootURL, checkerOptions{})

	go c.Check()

	i := 0

	for r := range c.Results() {
		i += strings.Count(r.String(true), "\n") + 1
	}

	assert.Equal(t, 4, i)
}

func TestCheckerCheckPage(t *testing.T) {
	c, _ := newChecker(rootURL, checkerOptions{})

	r, err := c.fetcher.Fetch(existentURL)
	assert.Nil(t, err)

	p, ok := r.Page()
	assert.True(t, ok)

	go c.checkPage(p)

	assert.True(t, (<-c.Results()).OK())
}

func TestCheckerCheckWithExcludedURLs(t *testing.T) {
	r, err := regexp.Compile("bar")
	assert.Nil(t, err)

	c, _ := newChecker(erroneousURL, checkerOptions{
		fetcherOptions: fetcherOptions{ExcludedPatterns: []*regexp.Regexp{r}},
	})

	go c.Check()

	assert.Equal(t, 2, strings.Count((<-c.Results()).String(true), "\n"))
}

func TestCheckerCheckPageError(t *testing.T) {
	for _, s := range []string{erroneousURL} {
		c, _ := newChecker(rootURL, checkerOptions{})

		r, err := c.fetcher.Fetch(s)
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
		close(c.channel)

		assert.Equal(t, c.slice, stringChannelToSlice(c.channel))
	}
}
