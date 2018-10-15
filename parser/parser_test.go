package parser

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseFile(t *testing.T) {
	s, _ := ParseFile("test.proto")

	d, _ := json.Marshal(s)
	fmt.Printf("%s", d)
}
