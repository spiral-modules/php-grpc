package main

import (
	"./service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("server.crt", "")
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial("localhost:9001", grpc.WithTransportCredentials(creds))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := service.NewEchoClient(conn)

	log.Println(client.Ping(metadata.AppendToOutgoingContext(context.Background(), "key", "value"), &service.Message{
		Msg: "hello world",
	}))
}
