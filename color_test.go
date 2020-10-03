package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsColorEnable(t *testing.T) {
	assert.True(t, isColorEnabled(auto, true))
	assert.False(t, isColorEnabled(auto, false))
	assert.True(t, isColorEnabled(always, true))
	assert.True(t, isColorEnabled(always, false))
	assert.False(t, isColorEnabled(never, true))
	assert.False(t, isColorEnabled(never, false))
}
