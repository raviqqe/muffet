package main

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func TestMarshalErrorJSONPageResult(t *testing.T) {
	bs, err := json.Marshal(newJSONPageResult(
		&pageResult{
			"http://foo.com",
			[]*successLinkResult{},
			[]*errorLinkResult{
				{"http://foo.com/bar", errors.New("baz")},
			},
		}, false))
	assert.Nil(t, err)
	cupaloy.SnapshotT(t, bs)
}

func TestMarshalSuccessJSONPageResult(t *testing.T) {
	bs, err := json.Marshal(newJSONPageResult(
		&pageResult{
			"http://foo.com",
			[]*successLinkResult{
				{"http://foo.com/foo", 200},
			},
			[]*errorLinkResult{},
		}, false))
	assert.Nil(t, err)
	cupaloy.SnapshotT(t, bs)
}

func TestMarshalVerboseSuccessJSONPageResult(t *testing.T) {
	bs, err := json.Marshal(newJSONPageResult(
		&pageResult{
			"http://foo.com",
			[]*successLinkResult{
				{"http://foo.com/foo", 200},
			},
			[]*errorLinkResult{},
		}, true))
	assert.Nil(t, err)
	cupaloy.SnapshotT(t, bs)
}
