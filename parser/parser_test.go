package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseFile(t *testing.T) {
	services, err := ParseFile("test.proto")
	assert.NoError(t, err)
	assert.Len(t, services, 2)
}

func TestParseNotFound(t *testing.T) {
	_, err := ParseFile("test2.proto")
	assert.Error(t, err)
}
