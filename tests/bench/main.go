package main

import (
	bench "../proto"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"sync"
	"sync/atomic"
	"time"
)

var in chan *bench.Message
var stop chan interface{}
var wg sync.WaitGroup

var done, fail int64

var n = flag.Int("n", 0, "Number of messages")
var c = flag.Int("c", 1, "Parallel connections")

func init() {
	in = make(chan *bench.Message)
	stop = make(chan interface{})
}

func main() {
	flag.Parse()

	for i := 0; i < *c; i++ {
		go worker()
	}

	go func() {
		for {
			select {
			case <-stop:
				return
			case <-time.NewTicker(time.Second).C:
				fmt.Printf("Done %v, failed %v \n", atomic.LoadInt64(&done), atomic.LoadInt64(&fail))
			}
		}
	}()

	start := time.Now()
	for i := 0; i < *n; i++ {
		wg.Add(1)
		in <- &bench.Message{Msg: "PING"}
	}

	wg.Wait()
	elapsed := time.Now().Sub(start).Seconds()
	fmt.Printf("Elapsed %v s, %v rps\n", elapsed, float64(*n)/elapsed)
	fmt.Printf("Done %v, failed %v\n", done, fail)
}

func worker() {
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		select {
		case <-in:
			in, out := new(bench.Message), new(bench.Message)
			in.Msg = "message forwarding?"
			err := conn.Invoke(context.Background(), "/test.PHP/Ping", in, out)

			if err == nil {
				atomic.AddInt64(&done, 1)
			} else {
				atomic.AddInt64(&fail, 1)
			}

			wg.Done()
		case <-stop:
			return
		}
	}
}
