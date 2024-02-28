package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsingFixedStatusCode(t *testing.T) {
	r, err := parseStatusCodeRange("403")

	assert.Nil(t, err)
	assert.Equal(t, 403, r.start)
	assert.Equal(t, 404, r.end)
}

func TestParsingStatusCodeRange(t *testing.T) {
	r, err := parseStatusCodeRange("200..300")

	assert.Nil(t, err)
	assert.Equal(t, 200, r.start)
	assert.Equal(t, 300, r.end)
}

func TestParsingInvalidStatusCode(t *testing.T) {
	_, err := parseStatusCodeRange("foo")

	assert.NotNil(t, err)
}

func TestInRangeOfStatusCode(t *testing.T) {
	r := statusCodeRange{200, 300}

	assert.False(t, r.isInRange(199))
	assert.True(t, r.isInRange(200))
	assert.True(t, r.isInRange(201))

	assert.True(t, r.isInRange(298))
	assert.True(t, r.isInRange(299))
	assert.False(t, r.isInRange(300))
}
