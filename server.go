package grpc

import (
	"github.com/spiral/roadrunner"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server interface{}

type Server struct {
	Name     string
	Metadata string
}

func (s *Server) ServiceDesc(r *roadrunner.Server) *grpc.ServiceDesc {
	d := grpc.ServiceDesc{
		ServiceName: s.Name,
		HandlerType: (*server)(nil),
		Methods:     []grpc.MethodDesc{},
		Streams:     []grpc.StreamDesc{},
		Metadata:    s.Metadata,
	}

	m := grpc.MethodDesc{
		MethodName: "Ping",
		Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
			// from pool (not sure, to payload directly)?
			msg := new(rawMessage)
			dec(msg)

			resp, err := r.Exec(&roadrunner.Payload{Body: *msg})
			if err != nil {
				return nil, err
			}

			*msg = resp.Body

			return msg, nil
		},
	}

	d.Methods = append(d.Methods, m)

	return &d
}
