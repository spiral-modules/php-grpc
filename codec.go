package grpc

import (
	"google.golang.org/grpc/encoding"
)

type rawMessage []byte

type CodecWrapper struct {
	Base encoding.Codec
}

// Marshal returns the wire format of v. rawMessages would be returned without encoding.
func (c *CodecWrapper) Marshal(v interface{}) ([]byte, error) {
	if raw, ok := v.(*rawMessage); ok {
		return *raw, nil
	}

	return c.Base.Marshal(v)
}

// Unmarshal parses the wire format into v. rawMessages would not be unmarshalled.
func (c *CodecWrapper) Unmarshal(data []byte, v interface{}) error {
	if raw, ok := v.(*rawMessage); ok {
		*raw = data
		return nil
	}

	return c.Base.Unmarshal(data, v)
}

// String return codec name.
func (c *CodecWrapper) String() string {
	return "raw:" + c.Base.Name()
}
