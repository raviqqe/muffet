package main

import (
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestPageChecker(c *fakeHttpClient) *pageChecker {
	return newPageChecker(
		newLinkFetcher(
			c,
			[]pageParser{newHtmlPageParser(newLinkFinder(nil, nil))},
			linkFetcherOptions{},
		),
		newLinkValidator("foo.com", nil, nil),
		false,
	)
}

func newTestPage(t *testing.T, fragments map[string]struct{}, links map[string]error) page {
	u, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	return newHtmlPage(u, fragments, links)
}

func TestPageCheckerCheckOnePage(t *testing.T) {
	c := newTestPageChecker(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				return nil, errors.New("")
			},
		),
	)

	go c.Check(newTestPage(t, nil, nil))

	i := 0

	for r := range c.Results() {
		i++
		assert.True(t, r.OK())
	}

	assert.Equal(t, 1, i)
}

func TestPageCheckerCheckTwoPages(t *testing.T) {
	c := newTestPageChecker(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				s := "http://foo.com/foo"

				if u.String() != s {
					return nil, errors.New("")
				}

				return newFakeHtmlResponse(s, ""), nil
			},
		),
	)

	go c.Check(
		newTestPage(t, nil, map[string]error{"http://foo.com/foo": nil}),
	)

	i := 0

	for r := range c.Results() {
		i++
		assert.True(t, r.OK())
	}

	assert.Equal(t, 2, i)
}

func TestPageCheckerFailToCheckPage(t *testing.T) {
	c := newTestPageChecker(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				return nil, errors.New("")
			},
		),
	)

	go c.Check(
		newTestPage(t, nil, map[string]error{"http://foo.com/foo": nil}),
	)

	assert.False(t, (<-c.Results()).OK())
}

func TestPageCheckerDoNotCheckSamePageTwice(t *testing.T) {
	c := newTestPageChecker(
		newFakeHttpClient(
			func(u *url.URL) (*fakeHttpResponse, error) {
				return newFakeHtmlResponse("http://foo.com", ""), nil
			},
		),
	)

	go c.Check(newTestPage(t, nil, map[string]error{"http://foo.com": nil}))

	i := 0

	for range c.Results() {
		i++
	}

	assert.Equal(t, 1, i)
}
