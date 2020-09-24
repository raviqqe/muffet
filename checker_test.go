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
			newPageParser(newLinkFinder(nil)),
			linkFetcherOptions{},
		),
		newLinkValidator("foo.com", nil, nil),
		512,
		false,
	)
}

func createTestPage(t *testing.T, fragments map[string]struct{}, links map[string]error) *page {
	u, err := url.Parse("http://foo.com")
	assert.Nil(t, err)

	return newPage(u, fragments, links)
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
