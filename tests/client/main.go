package main

import (
	test "../proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	start := time.Now()
	var n = 10000
	for i := 0; i < n; i++ {
		in, out := new(test.Message), new(test.Message)
		in.Msg = "message forwarding?"
		err := conn.Invoke(context.Background(), "/test.PHP/Ping", in, out)

		if err != nil {
			panic(err)
		}
	}
	elapsed := time.Now().Sub(start).Seconds()
	fmt.Printf("elapsed %v s, %v rps\n", elapsed, float64(n)/elapsed)
}
