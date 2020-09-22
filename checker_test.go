package main

import (
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestChecker(c *fakeHTTPClient) *checker {
	return newChecker(
		newLinkFetcher(
			c,
			newPageParser(newLinkFinder(nil), false),
			linkFetcherOptions{},
		),
		newLinkValidator("foo.com", nil, nil),
		512,
		false,
	)
}

func createTestPage(t *testing.T, ids map[string]struct{}, links map[string]error) *page {
	u, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	return newPage(u, ids, links)
}

func TestCheckerCheckOnePage(t *testing.T) {
	c := createTestChecker(
		newFakeHTTPClient(
			func(u *url.URL) (*fakeHTTPResponse, error) {
				return nil, errors.New("")
			},
		),
	)

	go c.Check(createTestPage(t, nil, nil))

	i := 0

	for r := range c.Results() {
		i++
		assert.True(t, r.OK())
	}

	assert.Equal(t, 1, i)
}

func TestCheckerCheckTwoPages(t *testing.T) {
	c := createTestChecker(
		newFakeHTTPClient(
			func(u *url.URL) (*fakeHTTPResponse, error) {
				s := "http://foo.com/foo"

				if u.String() != s {
					return nil, errors.New("")
				}

				return newFakeHTTPResponse(200, s, "text/html", nil), nil
			},
		),
	)

	go c.Check(
		createTestPage(t, nil, map[string]error{"http://foo.com/foo": nil}),
	)

	i := 0

	for r := range c.Results() {
		i++
		assert.True(t, r.OK())
	}

	assert.Equal(t, 2, i)
}

func TestCheckerFailToCheckPage(t *testing.T) {
	c := createTestChecker(
		newFakeHTTPClient(
			func(u *url.URL) (*fakeHTTPResponse, error) {
				return nil, errors.New("")
			},
		),
	)

	go c.Check(
		createTestPage(t, nil, map[string]error{"http://foo.com/foo": nil}),
	)

	assert.False(t, (<-c.Results()).OK())
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
