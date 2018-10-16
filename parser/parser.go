package parser

import (
	pp "github.com/emicklei/proto"
	"os"
	"strings"
)

type Service struct {
	Package string
	Name    string
	Comment string
	Methods []Method
}

type Method struct {
	Name           string
	Comment        string
	StreamsRequest bool
	RequestType    string
	StreamsReturns bool
	ReturnsType    string
}

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

	return fetchServices(proto, pkg)
}

func fetchServices(proto *pp.Proto, pkg string) ([]Service, error) {
	services := make([]Service, 0)
	pp.Walk(proto, pp.WithService(func(service *pp.Service) {
		svc := handleService(service)
		svc.Package = pkg

		services = append(services, svc)
	}))

	return services, nil
}

func handleService(s *pp.Service) Service {
	return Service{
		Name:    s.Name,
		Comment: comment(s.Comment),
		Methods: methods(s),
	}
}

func comment(comment *pp.Comment) string {
	if comment == nil {
		return ""
	}

	return strings.Trim(strings.Join(comment.Lines, "\n"), "\r \n")
}

func methods(s *pp.Service) []Method {
	methods := make([]Method, 0)
	for _, e := range s.Elements {
		if m, ok := e.(*pp.RPC); ok {
			methods = append(methods, Method{
				Name:           m.Name,
				Comment:        comment(m.Comment),
				StreamsRequest: m.StreamsRequest,
				RequestType:    m.RequestType,
				StreamsReturns: m.StreamsReturns,
				ReturnsType:    m.ReturnsType,
			})
		}
	}

	return methods
}
