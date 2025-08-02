package main

import (
	"testing"
)

type fakeNetError struct{}

func (fakeNetError) Error() string   { return "network error" }
func (fakeNetError) Timeout() bool   { return true }
func (fakeNetError) Temporary() bool { return true }

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
