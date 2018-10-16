package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseFile(t *testing.T) {
	services, err := ParseFile("test.proto")
	assert.NoError(t, err)
	assert.Len(t, services, 2)

	assert.Equal(t, "app.namespace", services[0].Package)
}

func TestParseNotFound(t *testing.T) {
	_, err := ParseFile("test2.proto")
	assert.Error(t, err)
}

func TestParseBytes(t *testing.T) {
	services, err := ParseBytes([]byte{})
	assert.NoError(t, err)
	assert.Len(t, services, 0)
}

func TestParseString(t *testing.T) {
	services, err := ParseBytes([]byte(`
syntax = "proto3";
package app.namespace;

// Ping Service.
service PingService {
    // Ping Method.
    rpc Ping (Message) returns (Message) {
    }
}

// Pong service.
service PongService {
    rpc Pong (stream Message) returns (stream Message) {
    }
}

message Message {
    string msg = 1;
    int64 value = 2;
}
`))
	assert.NoError(t, err)
	assert.Len(t, services, 2)

	assert.Equal(t, "app.namespace", services[0].Package)
}
