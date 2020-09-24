package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPageResultOK(t *testing.T) {
	assert.True(t, (&pageResult{"", nil, nil}).OK())
	assert.False(t, (&pageResult{"", nil, []*errorLinkResult{{}}}).OK())
}
