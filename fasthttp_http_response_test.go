package main

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestFastHttpResponseDecodeGzipBody(t *testing.T) {
	b := bytes.Buffer{}
	w := gzip.NewWriter(&b)
	_, err := w.Write([]byte("foo"))
	assert.Nil(t, err)
	err = w.Close()
	assert.Nil(t, err)

	r := fasthttp.Response{}
	r.Header.Add("Content-Encoding", "gzip")
	r.SetBody(b.Bytes())

	bs, err := newFasthttpHttpResponse(nil, &r).Body()

	assert.Nil(t, err)
	assert.Equal(t, "foo", string(bs))
}

func TestFastHttpResponseDecodeDeflateBody(t *testing.T) {
	b := bytes.Buffer{}
	w := zlib.NewWriter(&b)
	_, err := w.Write([]byte("foo"))
	assert.Nil(t, err)
	err = w.Close()
	assert.Nil(t, err)

	r := fasthttp.Response{}
	r.Header.Add("Content-Encoding", "deflate")
	r.SetBody(b.Bytes())

	bs, err := newFasthttpHttpResponse(nil, &r).Body()

	assert.Nil(t, err)
	assert.Equal(t, "foo", string(bs))
}

func TestFastHttpResponseDecodeBrotliBody(t *testing.T) {
	b := bytes.Buffer{}
	w := brotli.NewWriter(&b)
	_, err := w.Write([]byte("foo"))
	assert.Nil(t, err)
	err = w.Close()
	assert.Nil(t, err)

	r := fasthttp.Response{}
	r.Header.Add("Content-Encoding", "br")
	r.SetBody(b.Bytes())

	bs, err := newFasthttpHttpResponse(nil, &r).Body()

	assert.Nil(t, err)
	assert.Equal(t, "foo", string(bs))
}
