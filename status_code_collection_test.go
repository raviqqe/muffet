package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsingEmptyStatusCodeCollection(t *testing.T) {
	collection, err := parseStatusCodeCollection("")

	assert.Nil(t, err)

	assert.False(t, collection.isInCollection(199))
	assert.True(t, collection.isInCollection(200))
	assert.True(t, collection.isInCollection(201))

	assert.True(t, collection.isInCollection(298))
	assert.True(t, collection.isInCollection(299))
	assert.False(t, collection.isInCollection(300))
}

func TestParsingValidStatusCodeCollection(t *testing.T) {
	collection, err := parseStatusCodeCollection("200..206,403")

	assert.Nil(t, err)

	assert.False(t, collection.isInCollection(199))
	assert.True(t, collection.isInCollection(200))
	assert.True(t, collection.isInCollection(201))

	assert.True(t, collection.isInCollection(205))
	assert.True(t, collection.isInCollection(206))
	assert.False(t, collection.isInCollection(207))

	assert.False(t, collection.isInCollection(402))
	assert.True(t, collection.isInCollection(403))
	assert.False(t, collection.isInCollection(404))
}

func TestParsingInvalidStatusCodeCollection(t *testing.T) {
	collection, err := parseStatusCodeCollection("200,foo")

	assert.NotNil(t, err)
	assert.Nil(t, collection)
}
