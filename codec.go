package grpc

import (
	"google.golang.org/grpc/encoding"
)

type rawMessage []byte

const CodecName string = "proto"

func (r rawMessage) Reset()       {}
func (rawMessage) ProtoMessage()  {}
func (rawMessage) String() string { return "rawMessage" }

type Codec struct{ Base encoding.Codec }

func (c *Codec) Name() string {
	return CodecName
}

// Marshal returns the wire format of v. rawMessages would be returned without encoding.
func (c *Codec) Marshal(v interface{}) ([]byte, error) {
	if raw, ok := v.(rawMessage); ok {
		return raw, nil
	}

	return c.Base.Marshal(v)
}

// Unmarshal parses the wire format into v. rawMessages would not be unmarshalled.
func (c *Codec) Unmarshal(data []byte, v interface{}) error {
	if raw, ok := v.(*rawMessage); ok {
		*raw = data
		return nil
	}

	return c.Base.Unmarshal(data, v)
}

// String return Codec name.
func (c *Codec) String() string {
	return "raw:" + c.Base.Name()
}
