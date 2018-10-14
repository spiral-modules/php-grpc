package main

import (
	"../proto"
	"context"
	"google.golang.org/grpc"
	"log"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := test.NewPingClient(conn)
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(
		ctx,
		metadata.Pairs("key1", "val1", "key2", "val2"),
	)

	log.Println(client.Ping(ctx, &test.Message{
		Msg: "hi",
	}))
}
