package grpc

import (
	"github.com/spiral/roadrunner"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"sync"
)

type proxyService interface {
	// AddMethod registers new RPC method.
	AddMethod(method string)

	// ServiceDesc returns service description for the proxy.
	ServiceDesc() *grpc.ServiceDesc
}

// Proxy manages GRPC/RoadRunner bridge.
type Proxy struct {
	rr       *roadrunner.Server
	msgPool  sync.Pool
	name     string
	metadata string
	methods  []string
}

// NewProxy creates new service proxy object.
func NewProxy(name string, metadata string, rr *roadrunner.Server) *Proxy {
	return &Proxy{
		rr:       rr,
		msgPool:  sync.Pool{New: func() interface{} { return &rawMessage{} }},
		name:     name,
		metadata: metadata,
		methods:  make([]string, 0),
	}
}

// Attach attaches proxy to the GRPC server.
func (p *Proxy) Attach(server *grpc.Server) {
	server.RegisterService(p.ServiceDesc(), p)
}

// AddMethod registers new RPC method.
func (p *Proxy) AddMethod(method string) {
	p.methods = append(p.methods, method)
}

// ServiceDesc returns service description for the proxy.
func (p *Proxy) ServiceDesc() *grpc.ServiceDesc {
	desc := &grpc.ServiceDesc{
		ServiceName: p.name,
		Metadata:    p.metadata,
		HandlerType: (*proxyService)(nil),
		Methods:     []grpc.MethodDesc{},
		Streams:     []grpc.StreamDesc{},
	}

	// Registering methods
	for _, m := range p.methods {
		desc.Methods = append(desc.Methods, grpc.MethodDesc{
			MethodName: m,
			Handler:    p.methodHandler(m),
		})
	}

	return desc
}

// Generate method handler proxy.
func (p *Proxy) methodHandler(method string) func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		msg := p.msgPool.Get().(*rawMessage)
		defer p.msgPool.Put(msg)

		// decode incoming message
		dec(msg)

		resp, err := p.rr.Exec(p.makePayload(ctx, msg))

		// todo: wrap error (!)
		if err != nil {
			return nil, err
		}

		return resp.Body, nil
	}
}

// makePayload generates RoadRunner compatible payload based on GRPC message.
func (p *Proxy) makePayload(ctx context.Context, msg *rawMessage) *roadrunner.Payload {
	log.Println(ctx)
	return &roadrunner.Payload{
		Body: *msg,
	}
}
