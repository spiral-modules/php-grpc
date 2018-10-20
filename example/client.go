package main

import (
	"./proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
)

func main() {
	conn := makeConn()
	defer conn.Close()

	client := proto.NewEchoClient(conn)

	log.Println(
		client.Ping(
			metadata.AppendToOutgoingContext(context.Background(), "key", "value"),
			&proto.Message{Msg: "hello world"},
		),
	)
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
