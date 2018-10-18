package grpc

import (
	"encoding/json"
	"github.com/spiral/roadrunner"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"sync"
)

type proxyService interface {
	// RegisterMethod registers new RPC method.
	RegisterMethod(method string)

	// ServiceDesc returns service description for the proxy.
	ServiceDesc() *grpc.ServiceDesc
}

type rpcContext struct {
	Service string                 `json:"service"`
	Method  string                 `json:"method"`
	Context map[string]interface{} `json:"context"`
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
		msgPool:  sync.Pool{New: func() interface{} { return rawMessage{} }},
		name:     name,
		metadata: metadata,
		methods:  make([]string, 0),
	}
}

// RegisterMethod registers new RPC method.
func (p *Proxy) RegisterMethod(method string) {
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
		msg := p.msgPool.Get().(rawMessage)

		if err := dec(&msg); err != nil {
			return nil, err
		}

		resp, err := p.rr.Exec(p.makePayload(ctx, method, msg))
		p.msgPool.Put(msg)

		// todo: wrap error (!)
		if err != nil {
			return nil, err
		}

		return rawMessage(resp.Body), nil
	}
}

// makePayload generates RoadRunner compatible payload based on GRPC message. todo: return error
func (p *Proxy) makePayload(ctx context.Context, method string, msg rawMessage) *roadrunner.Payload {
	ctxData, _ := json.Marshal(rpcContext{
		Service: p.name,
		Method:  method,
		// pack context
	})

	// 	log.Println(ctx)
	return &roadrunner.Payload{
		Context: ctxData,
		Body:    msg,
	}
}
