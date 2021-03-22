package logging

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"golang.org/x/exp/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	"github.com/kapetndev/connect/transport"
)

// jsonpbMarshaller is the marshaller used for serializing protobuf messages.
var jsonpbMarshaller = &jsonpb.Marshaler{}

// UnaryServerInterceptor is a server side unary interceptor logging the
// payloads for a single request/response.
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := applyOptions(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		// Configure the logger passed into the middleware.
		logger := New(o.handler)

		// Invoke the handler and log the response.
		resp, err := handler(NewContext(ctx, logger), req)

		// Suppress request logs matching some pattern.
		if o.shouldDiscard(ctx, info.FullMethod, err) {
			return resp, err
		}

		if err != nil {
			o.handler.Handle(ctx, newRPCErrorRecord(ctx, startTime, info.FullMethod, err))
			return resp, err
		}

		// Log the request/response.
		o.handler.Handle(ctx, newRPCRecord(ctx, startTime, info.FullMethod, resp))
		return resp, err
	}
}

// StreamServerInterceptor is a server side stream interceptor logging the
// payloads for a single stream. Unlike the unary interceptor the payload of
// each message in the stream will be collected and logged together.
func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	o := applyOptions(opts)
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()
		ctx := ss.Context()

		// Configure the logger passed into the middleware.
		logger := New(o.handler)

		ss, err := transport.NewServerStreamWithContext(NewContext(ctx, logger), ss)
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		// Invoke the handler and log the response.
		err = handler(srv, ss)

		// Suppress request logs matching some pattern.
		if o.shouldDiscard(ctx, info.FullMethod, err) {
			return err
		}

		if err != nil {
			o.handler.Handle(ctx, newRPCErrorRecord(ctx, startTime, info.FullMethod, err))
			return err
		}

		// Log the request/response.
		o.handler.Handle(ctx, newRPCRecord(ctx, startTime, info.FullMethod, nil))
		return err
	}
}

func newRPCErrorRecord(ctx context.Context, t time.Time, path string, err error) slog.Record {
	record := newCommonRecord(ctx, slog.LevelError, t, "POST", path)
	record.AddAttrs(slog.String("error", err.Error()))
	return record
}

func newRPCRecord(ctx context.Context, t time.Time, path string, pbMsg interface{}) slog.Record {
	record := newCommonRecord(ctx, slog.LevelInfo, t, "POST", path)

	// If the response includes a payload then add it to the log entry. This
	// assumes that the payload is a JSON object.
	if p, ok := pbMsg.(proto.Message); ok {
		record.AddAttrs(slog.Any(ResponseKey, &jsonpbMarshalleble{p}))
	}

	return record
}

// jsonpbMarshalleble is a wrapper type allowing us to implement the Marshaler
// interface for protobuf message types.
type jsonpbMarshalleble struct {
	proto.Message
}

// MarshalJSON handles generating a slice of bytes representing the protobuf
// message payload as JSON.
func (j *jsonpbMarshalleble) MarshalJSON() ([]byte, error) {
	b := &bytes.Buffer{}

	if err := jsonpbMarshaller.Marshal(b, j.Message); err != nil {
		return nil, fmt.Errorf("failed to marshal jsonpb: %s", err)
	}

	return b.Bytes(), nil
}
