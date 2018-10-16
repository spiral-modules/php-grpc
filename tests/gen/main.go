package main

import (
	"fmt"
	parser "github.com/emicklei/proto"
	"os"
)

func main() {
	reader, _ := os.Open("test.proto")
	defer reader.Close()

	p := parser.NewParser(reader)
	definition, _ := p.Parse()

	parser.Walk(definition, parser.WithService(handleService))
}

func handleService(s *parser.Service) {
	fmt.Println(s.Name)

	for _, e := range s.Elements {
		if m, ok := e.(*parser.RPC); ok {
			fmt.Println(m.Name)
		}
	}
}
