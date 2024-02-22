package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsingFixedStatusCode(t *testing.T) {
	code, err := parseStatusCodeRange("403")

	assert.Nil(t, err)
	assert.Equal(t, 403, code.start)
	assert.Equal(t, 403, code.end)
}

func TestParsingStatusCodeRange(t *testing.T) {
	code, err := parseStatusCodeRange("200..299")

	assert.Nil(t, err)
	assert.Equal(t, 200, code.start)
	assert.Equal(t, 299, code.end)
}

func TestParsingInvalidStatusCode(t *testing.T) {
	code, err := parseStatusCodeRange("foo")

	assert.NotNil(t, err)
	assert.Nil(t, code)
}

func TestInRangeOfStatusCode(t *testing.T) {
	code := statusCodeRange{200, 299}

	assert.False(t, code.isInRange(199))
	assert.True(t, code.isInRange(200))
	assert.True(t, code.isInRange(201))

	assert.True(t, code.isInRange(298))
	assert.True(t, code.isInRange(299))
	assert.False(t, code.isInRange(300))
}
