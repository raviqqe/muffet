package main

import (
	"encoding/base64"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	for _, ss := range [][]string{
		{"-x", rootURL},
		{"-j", authorizationHeader("me:password"), basicAuthURL},
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
		{"-j", authorizationHeader("you:password"), basicAuthURL},
	} {
		_, err := command(ss, ioutil.Discard)

		assert.NotNil(t, err)
	}
}

func authorizationHeader(s string) string {
	return "Authorization: Basic " + base64.StdEncoding.EncodeToString([]byte(s))
}
