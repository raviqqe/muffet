package main

import (
	"bytes"
	"encoding/gob"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestNewFetchResult(t *testing.T) {
	newFetchResult(200, nil)
}

func TestNewFetchResultWithPage(t *testing.T) {
	p, err := newPage("", dummyHTML(t), newScraper(nil))
	assert.Nil(t, err)

	newFetchResult(200, p)
}

func TestFetchResultStatusCode(t *testing.T) {
	assert.Equal(t, 200, newFetchResult(200, nil).StatusCode())
}

func TestFetchResultPage(t *testing.T) {
	p, ok := newFetchResult(200, nil).Page()

	assert.False(t, ok)
	assert.Equal(t, (*page)(nil), p)

	q, err := newPage("", dummyHTML(t), newScraper(nil))
	assert.Nil(t, err)

	p, ok = newFetchResult(200, q).Page()

	assert.True(t, ok)
	assert.Equal(t, q, p)
}

func TestFetchResultEncodeDecode(t *testing.T) {
	for _, r := range []fetchResult{
		fetchResult{200, nil},
		func() fetchResult {
			n, err := html.Parse(strings.NewReader(htmlWithBody(`
				<a href=":" />
				<a href="mailto:me@right.here" />
				<a href="/bar" />
				<a href="#foo" />
			`)))
			assert.Nil(t, err)

			p, err := newPage("https://foo.com", n, newScraper(nil))
			assert.Nil(t, err)

			return newFetchResult(200, p)
		}(),
	} {
		b := bytes.Buffer{}

		assert.Nil(t, gob.NewEncoder(&b).Encode(&r))

		rr := fetchResult{}
		assert.Nil(t, gob.NewDecoder(&b).Decode(&rr))

		assert.Equal(t, r.statusCode, rr.statusCode)
		assertPagesEqual(t, r.page, rr.page)
	}
}

func TestFetchResultUnmarshalError(t *testing.T) {
	r := fetchResult{}
	assert.NotNil(t, r.UnmarshalBinary(nil))
}
