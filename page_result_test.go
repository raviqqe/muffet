package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPageResult(t *testing.T) {
	newPageResult("https://foo.com", nil, nil)
}

func TestPageResultOK(t *testing.T) {
	assert.True(t, newPageResult("https://foo.com", nil, nil).OK())
	assert.False(t, newPageResult("https://foo.com", nil, []string{"Oh, no!"}).OK())
}

func TestPageResultString(t *testing.T) {
	r := newPageResult("https://foo.com", []string{"foo"}, []string{"bar"})
	qs := r.String(false)
	vs := r.String(true)

	assert.Equal(t, 1, strings.Count(qs, "\n"))
	assert.Equal(t, 2, strings.Count(vs, "\n"))

	assert.True(t, strings.Contains(qs, "bar"))
	assert.True(t, strings.Contains(vs, "foo") && strings.Contains(vs, "bar"))
}
