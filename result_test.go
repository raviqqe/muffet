package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResult(t *testing.T) {
	newResult("https://foo.com", nil, nil)
}

func TestNewResultWithError(t *testing.T) {
	newResultWithError("https://foo.com", errors.New(""))
}

func TestResultOK(t *testing.T) {
	assert.True(t, newResult("https://foo.com", nil, nil).OK())
	assert.False(t, newResult("https://foo.com", nil, []string{"Oh, no!"}).OK())
}

func TestResultString(t *testing.T) {
	r := newResult("https://foo.com", []string{"foo"}, []string{"bar"})
	qs := r.String(false)
	vs := r.String(true)

	assert.Equal(t, 1, strings.Count(qs, "\n"))
	assert.Equal(t, 2, strings.Count(vs, "\n"))

	assert.True(t, strings.Contains(qs, "bar"))
	assert.True(t, strings.Contains(vs, "foo") && strings.Contains(vs, "bar"))
}
