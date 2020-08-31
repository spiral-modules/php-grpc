package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// EventUnaryCall thrown after every unary call.
const EventUnaryCall = 8001

// UnaryCallEvent contains information about invoked method.
type UnaryCallEvent struct {
	// Info contains unary call info.
	Info *grpc.UnaryServerInfo

	// Context associated with the call.
	Context context.Context

	// Error associated with event.
	Error error

	// event timings
	start   time.Time
	elapsed time.Duration
}

// Elapsed returns duration of the invocation.
func (e *UnaryCallEvent) Elapsed() time.Duration {
	return e.elapsed
}
