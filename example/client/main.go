package main

import (
	"fmt"
	"github.com/spiral/php-grpc/example/client/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("message text is required (go run client.go \"hello world\")")
	}

	conn := makeConn()
	defer conn.Close()

	client := proto.NewEchoClient(conn)

	resp, err := client.Ping(
		metadata.AppendToOutgoingContext(context.Background(), "key", "value"),
		&proto.Message{Msg: os.Args[1]},
	)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Response: %s\n", resp.Msg)
}

func makeConn() *grpc.ClientConn {
	creds, err := credentials.NewClientTLSFromFile("server.crt", "")
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial("localhost:9001", grpc.WithTransportCredentials(creds))
	if err != nil {
		panic(err)
	}

	return conn
}
