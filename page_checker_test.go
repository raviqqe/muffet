package main

import (
	"errors"
	"net/url"
	"testing"
	"sync/atomic"

	"github.com/stretchr/testify/assert"
)

func newTestPageChecker(c *fakeHttpClient) *pageChecker {
	return newTestPageCheckerWithRetries(c, 0)
}

func newTestPageCheckerWithRetries(c *fakeHttpClient, retries int) *pageChecker {
	return newPageChecker(
		newLinkFetcher(
			c,
			[]pageParser{newHtmlPageParser(newTestLinkFinder())},
			linkFetcherOptions{
				Retries: retries,
			},
		),
		newLinkValidator("foo.com", nil, nil),
		false,
	)
}

type fakeNetError struct{}

func (fakeNetError) Error() string   { return "network error" }
func (fakeNetError) Timeout() bool   { return true }
func (fakeNetError) Temporary() bool { return true }

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

func TestPageCheckerCheckPageRetry(t *testing.T) {
	for _, tt := range []struct {
		name             string
		errCnt           int
		expectedRequests int
		success          bool
	}{
		{name: "no errors", errCnt: 0, expectedRequests: 1, success: true},
		{name: "2 errors", errCnt: 2, expectedRequests: 3, success: true},
		{name: "3 errors", errCnt: 3, expectedRequests: 3, success: false},
	} {
		t.Run(
			tt.name,
			func(t *testing.T) {
				var reqCnt atomic.Int32
				c := newTestPageCheckerWithRetries(
					newFakeHttpClient(
						func(u *url.URL) (*fakeHttpResponse, error) {
							if u.String() == "http://foo.com/foo" {
								if reqCnt.Add(1) <= int32(tt.errCnt) {
									return nil, &fakeNetError{}
								}
								return newFakeHtmlResponse("http://foo.com/foo", ""), nil
							}
							return newFakeHtmlResponse("http://foo.com/", ""), nil
						},
					), 3,
				)

				go c.Check(
					newTestPage(t, nil, map[string]error{"http://foo.com/foo": nil}),
				)

				i := 0

				for r := range c.Results() {
					i++
					if tt.success {
						assert.True(t, r.OK())
					} else {
						assert.False(t, r.OK())
						assert.Len(t, r.ErrorLinkResults, 1)
						assert.Len(t, r.SuccessLinkResults, 0)
						assert.Equal(t, "http://foo.com/foo", r.ErrorLinkResults[0].URL)
					}
				}

				if tt.success {
					// initial page + 1 crawled page
					assert.Equal(t, 2, i)
				} else {
					// the crawled page failed
					assert.Equal(t, 1, i)
				}
				assert.Equal(t, int32(tt.expectedRequests), reqCnt.Load())
			},
		)
	}
}
