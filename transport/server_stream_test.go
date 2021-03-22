package transport_test

import (
	"context"
	"testing"

	"google.golang.org/grpc"

	"github.com/kapetndev/connect/transport"
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
	t.Parallel()

	t.Run("values from the server stream context are passed when wrapped", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), registrationKey{}, "NCC-1701")
		var ss grpc.ServerStream = &serverStream{ctx: ctx}

		ss, err := transport.NewServerStream(ss)
		if err != nil {
			t.Fatalf("error was not <nil>: %s", err)
		}

		if ss.Context().Value(registrationKey{}) != "NCC-1701" {
			t.Errorf("context has incorrect value: %s != %s", ss.Context().Value(registrationKey{}), "NCC-1701")
		}
	})

	t.Run("returns an error if the context is nil", func(t *testing.T) {
		ss := &serverStream{ctx: nil}

		_, err := transport.NewServerStream(ss)
		if err == nil {
			t.Error("error was <nil>")
		}
	})
}

func TestNewServerStreamWithContext(t *testing.T) {
	t.Parallel()

	t.Run("values provided by a new context are passed when wrapped", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), registrationKey{}, "NCC-1701")
		var ss grpc.ServerStream = &serverStream{ctx: context.Background()}

		ss, err := transport.NewServerStreamWithContext(ctx, ss)
		if err != nil {
			t.Fatalf("error was not <nil>: %s", err)
		}

		if ss.Context().Value(registrationKey{}) != "NCC-1701" {
			t.Errorf("context has incorrect value: %s != %s", ss.Context().Value(registrationKey{}), "NCC-1701")
		}
	})

	t.Run("returns an error if the context is nil", func(t *testing.T) {
		ss := &serverStream{ctx: context.Background()}

		_, err := transport.NewServerStreamWithContext(nil, ss)
		if err == nil {
			t.Error("error was <nil>")
		}
	})
}
