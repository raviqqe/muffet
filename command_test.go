package main

import (
	"encoding/base64"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	for _, args := range [][]string{
		{"-x", rootURL},
		{"-j", authorizationHeader("me:password"), basicAuthURL},
		{"-e", ".*", erroneousURL},
	} {
		ok := newCommand(ioutil.Discard, ioutil.Discard).Run(args)

		assert.True(t, ok)
	}
}

func TestCommandErroneousResult(t *testing.T) {
	for _, args := range [][]string{
		{erroneousURL},
	} {
		ok := newCommand(ioutil.Discard, ioutil.Discard).Run(args)

		assert.False(t, ok)
	}
}

func TestCommandError(t *testing.T) {
	for _, args := range [][]string{
		{":"},
		{"-t", "foo", rootURL},
		{"-j", authorizationHeader("you:password"), basicAuthURL},
	} {
		ok := newCommand(ioutil.Discard, ioutil.Discard).Run(args)

		assert.False(t, ok)
	}
}

func authorizationHeader(s string) string {
	return "Authorization: Basic " + base64.StdEncoding.EncodeToString([]byte(s))
}
