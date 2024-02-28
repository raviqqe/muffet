package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsingValidStatusCodeSet(t *testing.T) {
	for s, r := range map[string]statusCodeSet{
		"200":          {{200, 201}: {}},
		"200..300":     {{200, 300}: {}},
		"200..207,403": {{200, 207}: {}, {403, 404}: {}},
	} {
		t.Run(s, func(t *testing.T) {
			s, err := parseStatusCodeSet(s)

			assert.Nil(t, err)
			assert.Equal(t, s, r)
		})
	}
}

func TestParsingInvalidStatusCodeSet(t *testing.T) {
	_, err := parseStatusCodeSet("200,foo")

	assert.NotNil(t, err)
}
