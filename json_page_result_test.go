package main

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func TestMarshalErrorJSONPageResult(t *testing.T) {
	bs, err := json.Marshal(newJSONErrorPageResult(
		&pageResult{
			"http://foo.com",
			[]*successLinkResult{},
			[]*errorLinkResult{
				{"http://foo.com/bar", errors.New("baz"), 0},
			},
			0,
		}))
	assert.Nil(t, err)
	cupaloy.SnapshotT(t, bs)
}

func TestMarshalSuccessJSONPageResult(t *testing.T) {
	bs, err := json.Marshal(newJSONSuccessPageResult(
		&pageResult{
			"http://foo.com",
			[]*successLinkResult{
				{"http://foo.com/foo", 200, 0},
			},
			[]*errorLinkResult{},
			0,
		}))
	assert.Nil(t, err)
	cupaloy.SnapshotT(t, bs)
}
