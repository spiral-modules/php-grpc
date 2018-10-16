package grpc

import (
	"google.golang.org/grpc"
)

type Proxy struct {
	Name  string
	Proto string
}

func (p *Proxy) ServiceDesc() *grpc.ServiceDesc {
	return nil
	//d := grpc.ServiceDesc{
	//		ServiceName: s.Name,
	//		HandlerType: (*server)(nil),
	//		Methods:     []grpc.MethodDesc{},
	//		Streams:     []grpc.StreamDesc{},
	//	}

	// todo: multiple services?

	//m := grpc.MethodDesc{
	//	MethodName: "Ping",
	//	Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	//
	//		log.Println(ctx)
	//
	//		// from pool (not sure, to payload directly)?
	//		msg := new(rawMessage)
	//		dec(msg)
	//
	//		resp, err := r.Exec(&roadrunner.Payload{Body: *msg})
	//		if err != nil {
	//			return nil, err
	//		}
	//
	//		*msg = resp.Body
	//
	//		return msg, nil
	//	},
	//}

	//d.Methods = append(d.Methods, m)

	//	return &d
}
