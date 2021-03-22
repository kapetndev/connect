package transport

import (
	"context"
	"errors"

	"google.golang.org/grpc"
)

// serverStream implements a server side Stream.
type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

// NewServerStream returns a new ServerStream.
func NewServerStream(ss grpc.ServerStream) (*serverStream, error) {
	return NewServerStreamWithContext(ss.Context(), ss)
}

// NewServerStreamWithContext returns a new ServerStream with a context.
func NewServerStreamWithContext(ctx context.Context, ss grpc.ServerStream) (*serverStream, error) {
	if ctx == nil {
		return nil, errors.New("transport/grpc: nil Context")
	}

	return &serverStream{
		ServerStream: ss,
		ctx:          ctx,
	}, nil
}

// Context returns the underlying context.
func (ss *serverStream) Context() context.Context {
	return ss.ctx
}
