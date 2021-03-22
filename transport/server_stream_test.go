package transport_test

import (
	"context"
	"testing"

	"google.golang.org/grpc"

	"github.com/crumbandbase/expect"
	"github.com/crumbandbase/service-core-go/transport"
)

type registrationKey struct{}

type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (ss *serverStream) Context() context.Context {
	return ss.ctx
}

func TestNewServerStream(t *testing.T) {
	t.Run("values from the server stream context are passed when wrapped", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), registrationKey{}, "NCC-1701")
		var ss grpc.ServerStream = &serverStream{ctx: ctx}

		ss, err := transport.NewServerStream(ss)

		expect.Equal(t, err, nil, "error was not <nil>")
		expect.Equal(t, ss.Context().Value(registrationKey{}), "NCC-1701", "context has incorrect value")
	})

	t.Run("returns an error if the context is nil", func(t *testing.T) {
		ss := &serverStream{ctx: nil}

		_, err := transport.NewServerStream(ss)
		expect.NotEqual(t, err, nil, "error was <nil>")
	})
}

func TestNewServerStreamWithContext(t *testing.T) {
	t.Run("values provided by a new context are passed when wrapped", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), registrationKey{}, "NCC-1701")
		var ss grpc.ServerStream = &serverStream{ctx: context.Background()}

		ss, err := transport.NewServerStreamWithContext(ctx, ss)

		expect.Equal(t, err, nil, "error was not <nil>")
		expect.Equal(t, ss.Context().Value(registrationKey{}), "NCC-1701", "context has incorrect value")
	})

	t.Run("returns an error if the context is nil", func(t *testing.T) {
		ss := &serverStream{ctx: context.Background()}

		_, err := transport.NewServerStreamWithContext(nil, ss)
		expect.NotEqual(t, err, nil, "error was <nil>")
	})
}
