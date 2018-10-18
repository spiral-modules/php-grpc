package main

import (
	"./service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:9001", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := service.NewEchoClient(conn)
	log.Println(client.Ping(context.Background(), &service.Message{
		Msg: "hello world",
	}))
}
