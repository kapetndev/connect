package recovery_test

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kapetndev/connect/recovery"
	echopb "github.com/kapetndev/connect/testdata/echo/v1"
	"github.com/kapetndev/grpctest"
)

var (
	goodEchoRequest     = &echopb.EchoRequest{Message: "good"}
	panicEchoRequest    = &echopb.EchoRequest{Message: "panic"}
	nilPanicEchoRequest = &echopb.EchoRequest{Message: "nilPanic"}

	goodServerStreamingEchoRequest     = &echopb.ServerStreamingEchoRequest{Message: "good"}
	panicServerStreamingEchoRequest    = &echopb.ServerStreamingEchoRequest{Message: "panic"}
	nilPanicServerStreamingEchoRequest = &echopb.ServerStreamingEchoRequest{Message: "nilPanic"}

	nilEchoResponse                = (*echopb.EchoResponse)(nil)
	nilServerStreamingEchoResponse = (*echopb.ServerStreamingEchoResponse)(nil)
)

type echoServer struct {
	echopb.UnimplementedEchoServiceServer
}

func (s *echoServer) Echo(ctx context.Context, in *echopb.EchoRequest) (*echopb.EchoResponse, error) {
	returnPanics(in.Message)
	return &echopb.EchoResponse{Message: in.Message}, nil
}

func (s *echoServer) ServerStreamingEcho(in *echopb.ServerStreamingEchoRequest, ss echopb.EchoService_ServerStreamingEchoServer) error {
	returnPanics(in.Message)

	res := &echopb.ServerStreamingEchoResponse{Message: in.Message}
	if err := ss.Send(res); err != nil {
		return err
	}

	return nil
}

func setupRecoveryServer(t *testing.T, opts ...recovery.Option) (grpctest.Closer, echopb.EchoServiceClient) {
	s := grpctest.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(opts...),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(opts...),
		),
	)

	conn, err := s.ClientConn()
	if err != nil {
		t.Fatal(err)
	}

	echopb.RegisterEchoServiceServer(s, &echoServer{})
	s.Serve()

	return s.Close, echopb.NewEchoServiceClient(conn)
}

func TestUnaryServerInterceptor_Default(t *testing.T) {
	t.Parallel()

	t.Run("does not invoke the recovery function when the request is successful", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		resp, err := client.Echo(context.Background(), goodEchoRequest)
		if err != nil {
			t.Errorf("error was not <nil>: %s", err)
		}

		if resp.Message != goodEchoRequest.Message {
			t.Errorf("messages are not equal: %s != %s", resp.Message, goodEchoRequest.Message)
		}
	})

	t.Run("recovers and returns an error when the request panics", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		resp, err := client.Echo(context.Background(), panicEchoRequest)
		if err == nil {
			t.Errorf("error was <nil>")
		}

		if resp != nilEchoResponse {
			t.Errorf("responses are not equal: %v != %v", resp, nilEchoResponse)
		}

		statusErr := status.Convert(err)
		if statusErr.Code() != codes.Internal {
			t.Errorf("error codes are not equal: %d != %d", statusErr.Code(), codes.Internal)
		}

		if statusErr.Message() != panicMessage {
			t.Errorf("error messages are not equal: %s != %s", statusErr.Message(), panicMessage)
		}
	})

	t.Run("recovers and returns an error when the request causes a nil panic", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		resp, err := client.Echo(context.Background(), nilPanicEchoRequest)
		if err == nil {
			t.Errorf("error was <nil>")
		}

		if resp != nilEchoResponse {
			t.Errorf("responses are not equal: %v != %v", resp, nilEchoResponse)
		}

		statusErr := status.Convert(err)
		if statusErr.Code() != codes.Internal {
			t.Errorf("error codes are not equal: %d != %d", statusErr.Code(), codes.Internal)
		}

		if statusErr.Message() != "<nil>" {
			t.Errorf("error messages are not equal: %s != %s", statusErr.Message(), "<nil>")
		}
	})
}

func TestStreamServerInterceptor_Default(t *testing.T) {
	t.Parallel()

	t.Run("does not invoke the recovery function when the request is successful", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		stream, err := client.ServerStreamingEcho(context.Background(), goodServerStreamingEchoRequest)
		if err != nil {
			t.Errorf("error was not <nil>: %s", err)
		}

		resp, err := stream.Recv()
		if err != nil {
			t.Errorf("error was not <nil>: %s", err)
		}

		if resp.Message != goodServerStreamingEchoRequest.Message {
			t.Errorf("messages are not equal: %s != %s", resp.Message, goodServerStreamingEchoRequest.Message)
		}
	})

	t.Run("recovers and returns an error when the request panics", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		stream, err := client.ServerStreamingEcho(context.Background(), panicServerStreamingEchoRequest)
		if err != nil {
			t.Errorf("error was not <nil>: %s", err)
		}

		resp, err := stream.Recv()
		if err == nil {
			t.Errorf("error was <nil>")
		}

		if resp != nilServerStreamingEchoResponse {
			t.Errorf("responses are not equal: %v != %v", resp, nilServerStreamingEchoResponse)
		}

		statusErr := status.Convert(err)
		if statusErr.Code() != codes.Internal {
			t.Errorf("error codes are not equal: %d != %d", statusErr.Code(), codes.Internal)
		}

		if statusErr.Message() != panicMessage {
			t.Errorf("error messages are not equal: %s != %s", statusErr.Message(), panicMessage)
		}
	})

	t.Run("recovers and returns an error when the request causes a nil panic", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		stream, err := client.ServerStreamingEcho(context.Background(), nilPanicServerStreamingEchoRequest)
		if err != nil {
			t.Errorf("error was not <nil>: %s", err)
		}

		resp, err := stream.Recv()
		if err == nil {
			t.Errorf("error was <nil>")
		}

		if resp != nilServerStreamingEchoResponse {
			t.Errorf("responses are not equal: %v != %v", resp, nilServerStreamingEchoResponse)
		}

		statusErr := status.Convert(err)
		if statusErr.Code() != codes.Internal {
			t.Errorf("error codes are not equal: %d != %d", statusErr.Code(), codes.Internal)
		}

		if statusErr.Message() != "<nil>" {
			t.Errorf("error messages are not equal: %s != %s", statusErr.Message(), "<nil>")
		}
	})
}

func setupOverrideRecoveryServer(t *testing.T) (grpctest.Closer, echopb.EchoServiceClient) {
	recoveryFunc := func(ctx context.Context, p interface{}) error {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}

	return setupRecoveryServer(t, recovery.WithRecoveryContext(recoveryFunc))
}

func TestUnaryServerInterceptor_Override(t *testing.T) {
	t.Parallel()

	t.Run("does not invoke the recovery function when the request is successful", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		resp, err := client.Echo(context.Background(), goodEchoRequest)
		if err != nil {
			t.Errorf("error was not <nil>: %s", err)
		}

		if resp.Message != goodEchoRequest.Message {
			t.Errorf("messages are not equal: %s != %s", resp.Message, goodEchoRequest.Message)
		}
	})

	t.Run("recovers and returns an error when the request panics", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		resp, err := client.Echo(context.Background(), panicEchoRequest)
		if err == nil {
			t.Errorf("error was <nil>")
		}

		if resp != nilEchoResponse {
			t.Errorf("responses are not equal: %v != %v", resp, nilEchoResponse)
		}

		statusErr := status.Convert(err)
		if statusErr.Code() != codes.Internal {
			t.Errorf("error codes are not equal: %d != %d", statusErr.Code(), codes.Internal)
		}

		if statusErr.Message() != panicMessage {
			t.Errorf("error messages are not equal: %s != %s", statusErr.Message(), panicMessage)
		}
	})

	t.Run("recovers and returns an error when the request causes a nil panic", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		resp, err := client.Echo(context.Background(), nilPanicEchoRequest)
		if err == nil {
			t.Errorf("error was <nil>")
		}

		if resp != nilEchoResponse {
			t.Errorf("responses are not equal: %v != %v", resp, nilEchoResponse)
		}

		statusErr := status.Convert(err)
		if statusErr.Code() != codes.Internal {
			t.Errorf("error codes are not equal: %d != %d", statusErr.Code(), codes.Internal)
		}

		if statusErr.Message() != "<nil>" {
			t.Errorf("error messages are not equal: %s != %s", statusErr.Message(), "<nil>")
		}
	})
}

func TestStreamServerInterceptor_Override(t *testing.T) {
	t.Parallel()

	t.Run("does not invoke the recovery function when the request is successful", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		stream, err := client.ServerStreamingEcho(context.Background(), goodServerStreamingEchoRequest)
		if err != nil {
			t.Errorf("error was not <nil>: %s", err)
		}

		resp, err := stream.Recv()
		if err != nil {
			t.Errorf("error was not <nil>: %s", err)
		}

		if resp.Message != goodServerStreamingEchoRequest.Message {
			t.Errorf("messages are not equal: %s != %s", resp.Message, goodServerStreamingEchoRequest.Message)
		}
	})

	t.Run("recovers and returns an error when the request panics", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		stream, err := client.ServerStreamingEcho(context.Background(), panicServerStreamingEchoRequest)
		if err != nil {
			t.Errorf("error was not <nil>: %s", err)
		}

		resp, err := stream.Recv()
		if err == nil {
			t.Errorf("error was <nil>")
		}

		if resp != nilServerStreamingEchoResponse {
			t.Errorf("responses are not equal: %v != %v", resp, nilServerStreamingEchoResponse)
		}

		statusErr := status.Convert(err)
		if statusErr.Code() != codes.Internal {
			t.Errorf("error codes are not equal: %d != %d", statusErr.Code(), codes.Internal)
		}

		if statusErr.Message() != panicMessage {
			t.Errorf("error messages are not equal: %s != %s", statusErr.Message(), panicMessage)
		}
	})

	t.Run("recovers and returns an error when the request causes a nil panic", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		stream, err := client.ServerStreamingEcho(context.Background(), nilPanicServerStreamingEchoRequest)
		if err != nil {
			t.Errorf("error was not <nil>: %s", err)
		}

		resp, err := stream.Recv()
		if err == nil {
			t.Errorf("error was <nil>")
		}

		if resp != nilServerStreamingEchoResponse {
			t.Errorf("responses are not equal: %v != %v", resp, nilServerStreamingEchoResponse)
		}

		statusErr := status.Convert(err)
		if statusErr.Code() != codes.Internal {
			t.Errorf("error codes are not equal: %d != %d", statusErr.Code(), codes.Internal)
		}

		if statusErr.Message() != "<nil>" {
			t.Errorf("error messages are not equal: %s != %s", statusErr.Message(), "<nil>")
		}
	})
}
