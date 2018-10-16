package parser

import (
	pp "github.com/emicklei/proto"
	"os"
	"strings"
)

// Service contains information about singular GRPC service.
type Service struct {
	// Package defines service namespace.
	Package string

	// Name defines service name.
	Name string

	// Comment associated with the service.
	Comment string

	// Methods list.
	Methods []Method
}

// Method describes singular RPC method.
type Method struct {
	// Name is method name.
	Name string

	// Comment associated with method.
	Comment string

	// StreamsRequest defines if method accept stream input.
	StreamsRequest bool

	// RequestType defines message name (from the same package) of method input.
	RequestType string

	// StreamsReturns defines if method streams result.
	StreamsReturns bool

	// ReturnsType defines message name (from the same package) of method return value.
	ReturnsType string
}

// ParseFile parses given proto file or returns error.
func ParseFile(file string) ([]Service, error) {
	reader, _ := os.Open(file)
	defer reader.Close()

	proto, err := pp.NewParser(reader).Parse()
	if err != nil {
		return nil, err
	}

	var pkg string
	for _, e := range proto.Elements {
		if p, ok := e.(*pp.Package); ok {
			pkg = p.Name
		}
	}

	return parseServices(proto, pkg)
}

func parseServices(proto *pp.Proto, pkg string) ([]Service, error) {
	services := make([]Service, 0)
	pp.Walk(proto, pp.WithService(func(service *pp.Service) {
		services = append(services, Service{
			Package: pkg,
			Name:    service.Name,
			Comment: parseComment(service.Comment),
			Methods: parseMethods(service),
		})
	}))

	return services, nil
}

func parseMethods(s *pp.Service) []Method {
	methods := make([]Method, 0)
	for _, e := range s.Elements {
		if m, ok := e.(*pp.RPC); ok {
			methods = append(methods, Method{
				Name:           m.Name,
				Comment:        parseComment(m.Comment),
				StreamsRequest: m.StreamsRequest,
				RequestType:    m.RequestType,
				StreamsReturns: m.StreamsReturns,
				ReturnsType:    m.ReturnsType,
			})
		}
	}

	return methods
}

func parseComment(comment *pp.Comment) string {
	if comment == nil {
		return ""
	}

	return strings.Trim(strings.Join(comment.Lines, "\n"), "\r \n")
}
