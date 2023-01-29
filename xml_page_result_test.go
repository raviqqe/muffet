package main

import (
	"encoding/xml"
	"errors"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func TestMarshalErrorXMLPageResult(t *testing.T) {
	bs, err := marshalXML(newXMLPageResult(
		&pageResult{
			"http://foo.com",
			[]*successLinkResult{},
			[]*errorLinkResult{
				{"http://foo.com/bar", errors.New("baz")},
			},
		}))
	assert.Nil(t, err)
	cupaloy.SnapshotT(t, bs)
}

func TestMarshalSuccessXMLPageResult(t *testing.T) {
	bs, err := marshalXML(newXMLPageResult(
		&pageResult{
			"http://foo.com",
			[]*successLinkResult{
				{"http://foo.com/bar", 200},
			},
			[]*errorLinkResult{},
		}))
	assert.Nil(t, err)
	cupaloy.SnapshotT(t, bs)
}

func marshalXML(x any) ([]byte, error) {
	return xml.MarshalIndent(x, "", "  ")
}
