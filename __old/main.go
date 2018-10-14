package __old

//
//import (
//	"./test"
//	"context"
//	"fmt"
//	"github.com/emicklei/proto"
//	"github.com/sirupsen/logrus"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/encoding"
//	"net"
//	"os"
//)
//
//type PingServer struct {
//}
//
//func (p *PingServer) Ping(ctx context.Context, msg *test.Message) (*test.Message, error) {
//	logrus.Print("got message")
//	return msg, nil
//}
//
//type CodecWrapper struct {
//	Base encoding.Codec
//}
//
//var ServiceDesc = grpc.ServiceDesc{
//	ServiceName: "test.Ping",
//	HandlerType: (*interface{})(nil),
//	Methods: []grpc.MethodDesc{
//		{
//			MethodName: "Ping",
//			Handler:    _Ping_Handler,
//		},
//	},
//	Streams:  []grpc.StreamDesc{},
//	Metadata: "test.proto",
//}
//
//func _Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
//
//	logrus.Println(ctx)
//
//	in := new(test.Message)
//	if err := dec(in); err != nil {
//		return nil, err
//	}
//
//	if interceptor == nil {
//		return srv.(*PingServer).Ping(ctx, in)
//	}
//
//	info := &grpc.UnaryServerInfo{
//		Server:     srv,
//		FullMethod: "/test.Ping/Ping",
//	}
//
//	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
//		return srv.(*PingServer).Ping(ctx, req.(*test.Message))
//	}
//
//	return interceptor(ctx, in, info, handler)
//}
//
//// Marshal returns the wire format of v.
//func (c *CodecWrapper) Marshal(v interface{}) ([]byte, error) {
//	logrus.Warningf("got output: %v", v)
//
//	return c.Base.Marshal(v)
//}
//
//func (c *CodecWrapper) Unmarshal(data []byte, v interface{}) error {
//	logrus.Warningf("got input: %s", data)
//
//	return c.Base.Unmarshal(data, v)
//}
//
//func (c *CodecWrapper) String() string {
//	return "CodecWrapper:" + c.Base.Name()
//}
//func handleService(s *proto.Service) {
//	fmt.Println(s.Name)
//}
//
//func handleMessage(m *proto.Message) {
//	fmt.Println(m.Name)
//}
//
//func handleRPC(m *proto.RPC) {
//	fmt.Println(m.Name)
//}
//
//func main() {
//
//	reader, _ := os.Open("test.proto")
//	defer reader.Close()
//
//	parser := proto.NewParser(reader)
//	definition, _ := parser.Parse()
//
//	proto.Walk(definition,
//		proto.WithService(handleService),
//		proto.WithMessage(handleMessage),
//		proto.WithRPC(handleRPC))
//
//	return
//	mode := os.Args[1]
//	logrus.Println("working in mode: " + mode)
//
//	if mode == "server" {
//		server := grpc.NewServer(grpc.CustomCodec(&CodecWrapper{
//			Base: encoding.GetCodec("proto"),
//		}))
//
//		//Adding unit
//		server.RegisterService(&ServiceDesc, &PingServer{})
//
//		//Open connection
//		listener, err := net.Listen("tcp", "localhost:8000")
//		if err != nil {
//			panic(err)
//		}
//
//		server.Serve(listener)
//		return
//	}
//
//	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
//	if err != nil {
//		panic(err)
//	}
//
//	client := test.NewPingClient(conn)
//	resp, err := client.Ping(context.Background(), &test.Message{
//		Msg: "hello world",
//	})
//
//	if err != nil {
//		panic(err)
//	}
//
//	logrus.Print(resp)
//}
