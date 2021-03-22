package logging

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/crumbandbase/service-core-go/transport"
)

var (
	// JsonpbMarshaller is the marshaller used for serializing protobuf messages.
	JsonpbMarshaller = &jsonpb.Marshaler{}
)

// UnaryServerInterceptor is a server side unary interceptor logging the
// payloads for a single request/response.
func UnaryServerInterceptor(log Logger, opts ...Option) grpc.UnaryServerInterceptor {
	o := applyOptions(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		// Configure the logger passed into the middleware. This function will
		// always return a new logger and is therefore thread safe.
		entry := newRequestLogger(ctx, log, o.fields, "POST", info.FullMethod, startTime, o.timestampFormat)

		// Invoke the hander and log the response.
		resp, err := handler(NewContext(ctx, entry), req)

		// Suppress request logs matching some pattern.
		if o.shouldDiscard(ctx, info.FullMethod, err) {
			return resp, err
		}

		durationKey, duration := o.durationFieldValue(time.Since(startTime))
		logRPC(entry, resp, err, durationKey, duration)
		return resp, err
	}
}

// StreamServerInterceptor is a server side stream interceptor logging the
// payloads for a single stream. Unlike the unary interceptor the payload of
// each message in the stream will be collected and logged together.
func StreamServerInterceptor(log Logger, opts ...Option) grpc.StreamServerInterceptor {
	o := applyOptions(opts)
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()

		// Configure the logger passed into the middleware. This function will
		// always return a new logger and is therefore thread safe.
		entry := newRequestLogger(ss.Context(), log, o.fields, "POST", info.FullMethod, startTime, o.timestampFormat)

		ss, err := transport.NewServerStreamWithContext(NewContext(ss.Context(), entry), ss)
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		// Invoke the hander and log the response.
		err = handler(srv, ss)

		// Suppress request logs matching some pattern.
		if o.shouldDiscard(ss.Context(), info.FullMethod, err) {
			return err
		}

		durationKey, duration := o.durationFieldValue(time.Since(startTime))
		logRPC(entry, nil, err, durationKey, duration)
		return err
	}
}

func logRPC(log Logger, pbMsg interface{}, err error, durationKey string, duration interface{}) {
	log = log.WithField(durationKey, duration)
	log = entryWithProtoFields(log, pbMsg, PayloadKey)

	if err != nil {
		log.WithField("error", err.Error()).Error()
		return
	}

	log.Info()
}

func entryWithProtoFields(log Logger, pbMsg interface{}, key string) Logger {
	if p, ok := pbMsg.(proto.Message); ok {
		return log.WithField(key, &jsonpbMarshalleble{p})
	}

	return log
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

	if err := JsonpbMarshaller.Marshal(b, j.Message); err != nil {
		return nil, fmt.Errorf("failed to marshal jsonpb: %v", err)
	}

	return b.Bytes(), nil
}
