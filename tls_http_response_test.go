package main

import (
	"bytes"
	"io"
	"testing"

	fh "github.com/bogdanfinn/fhttp"
	"github.com/stretchr/testify/assert"
)

func TestTlsHttpResponseBody(t *testing.T) {
	b := bytes.Buffer{}
	_, err := b.Write([]byte("foo"))
	assert.Nil(t, err)

	r := fh.Response{}
	r.Body = io.NopCloser(bytes.NewReader(b.Bytes()))

	bs, err := newTlsHttpResponse(nil, &r).Body()

	assert.Nil(t, err)
	assert.Equal(t, "foo", string(bs))
}
