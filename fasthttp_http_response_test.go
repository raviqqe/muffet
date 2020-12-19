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

func TestFastHTTPResponseDecodeGzipBody(t *testing.T) {
	b := bytes.Buffer{}
	w := gzip.NewWriter(&b)
	w.Write([]byte("foo"))
	w.Close()

	r := fasthttp.Response{}
	r.Header.Add("Content-Encoding", "gzip")
	r.SetBody(b.Bytes())

	bs, err := newFasthttpHTTPResponse(nil, &r).Body()

	assert.Nil(t, err)
	assert.Equal(t, "foo", string(bs))
}

func TestFastHTTPResponseDecodeDeflateBody(t *testing.T) {
	b := bytes.Buffer{}
	w := zlib.NewWriter(&b)
	w.Write([]byte("foo"))
	w.Close()

	r := fasthttp.Response{}
	r.Header.Add("Content-Encoding", "deflate")
	r.SetBody(b.Bytes())

	bs, err := newFasthttpHTTPResponse(nil, &r).Body()

	assert.Nil(t, err)
	assert.Equal(t, "foo", string(bs))
}

func TestFastHTTPResponseDecodeBrotliBody(t *testing.T) {
	b := bytes.Buffer{}
	w := brotli.NewWriter(&b)
	w.Write([]byte("foo"))
	w.Close()

	r := fasthttp.Response{}
	r.Header.Add("Content-Encoding", "br")
	r.SetBody(b.Bytes())

	bs, err := newFasthttpHTTPResponse(nil, &r).Body()

	assert.Nil(t, err)
	assert.Equal(t, "foo", string(bs))
}
