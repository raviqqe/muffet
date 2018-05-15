package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	for _, ss := range [][]string{
		{"-x", rootURL},
		{"-a", "me:password", basicAuthURL},
	} {
		s, err := command(ss, ioutil.Discard)

		assert.Zero(t, s)
		assert.Nil(t, err)
	}
}

func TestCommandErroneousResult(t *testing.T) {
	for _, ss := range [][]string{
		{erroneousURL},
	} {
		s, err := command(ss, ioutil.Discard)

		assert.Equal(t, 1, s)
		assert.Nil(t, err)
	}
}

func TestCommandError(t *testing.T) {
	for _, ss := range [][]string{
		{":"},
		{"-t", "foo", rootURL},
		{"-a", "you:password", basicAuthURL},
	} {
		_, err := command(ss, ioutil.Discard)

		assert.NotNil(t, err)
	}
}
