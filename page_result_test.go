package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPageResultOK(t *testing.T) {
	d, _ := time.ParseDuration("1s")
	assert.True(t, (&pageResult{"", nil, nil, d}).OK())
	assert.False(t, (&pageResult{"", nil, []*errorLinkResult{{}}, d}).OK())
}
