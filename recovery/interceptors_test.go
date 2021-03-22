package recovery_test

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/crumbandbase/expect"
	"github.com/crumbandbase/grpctest"
	"github.com/crumbandbase/service-core-go/recovery"
	echopb "github.com/crumbandbase/service-core-go/testdata/proto"
)

const nilError = "<nil>"

var (
	goodGRPCRequest     = &echopb.EchoRequest{Message: "good"}
	panicGRPCRequest    = &echopb.EchoRequest{Message: "panic"}
	nilPanicGRPCRequest = &echopb.EchoRequest{Message: "nilPanic"}

	nilResponse = (*echopb.EchoResponse)(nil)
)

type echoServer struct {
	echopb.UnimplementedEchoServer
}

func (s *echoServer) Echo(ctx context.Context, in *echopb.EchoRequest) (*echopb.EchoResponse, error) {
	s.returnPanics(in.Message)
	return &echopb.EchoResponse{Message: in.Message}, nil
}

func (s *echoServer) ServerStreamingEcho(in *echopb.EchoRequest, ss echopb.Echo_ServerStreamingEchoServer) error {
	s.returnPanics(in.Message)

	res := &echopb.EchoResponse{Message: in.Message}
	if err := ss.Send(res); err != nil {
		return err
	}

	return nil
}

func (s *echoServer) returnPanics(m string) {
	switch m {
	case "panic":
		panic("very bad thing happened")
	case "nilPanic":
		panic(nil)
	}
}

func setupRecoveryServer(t *testing.T, opts ...recovery.Option) (grpctest.Closer, echopb.EchoClient) {
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

	echopb.RegisterEchoServer(s, &echoServer{})
	s.Serve()

	return s.Close, echopb.NewEchoClient(conn)
}

func TestUnaryServerInterceptor_Default(t *testing.T) {
	t.Parallel()

	t.Run("does not invoke the recovery function when the request is successful", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		resp, err := client.Echo(context.Background(), goodGRPCRequest)

		expect.Equal(t, err, nil, "error was not <nil>")
		expect.Equal(t, resp.Message, goodGRPCRequest.Message, "messages are not equal")
	})

	t.Run("recovers and returns an error when the request panics", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		resp, err := client.Echo(context.Background(), panicGRPCRequest)
		statusErr := status.Convert(err)

		expect.Equal(t, resp, nilResponse, "responses are not equal")
		expect.Equal(t, statusErr.Code(), codes.Internal, "error codes are not equal")
		expect.Equal(t, statusErr.Message(), "very bad thing happened", "error messages are not equal")
	})

	t.Run("recovers and returns an error when the request causes a nil panic", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		resp, err := client.Echo(context.Background(), nilPanicGRPCRequest)
		statusErr := status.Convert(err)

		expect.Equal(t, resp, nilResponse, "responses are not equal")
		expect.Equal(t, statusErr.Code(), codes.Internal, "error codes are no equal")
		expect.Equal(t, statusErr.Message(), nilError, "error messages are not equal")
	})
}

func TestStreamServerInterceptor_Default(t *testing.T) {
	t.Parallel()

	t.Run("does not invoke the recovery function when the request is successful", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		stream, err := client.ServerStreamingEcho(context.Background(), goodGRPCRequest)
		expect.Equal(t, err, nil, "error was not <nil>")

		resp, err := stream.Recv()

		expect.Equal(t, err, nil, "error was not <nil>")
		expect.Equal(t, resp.Message, goodGRPCRequest.Message, "messages are not equal")
	})

	t.Run("recovers and returns an error when the request panics", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		stream, err := client.ServerStreamingEcho(context.Background(), panicGRPCRequest)
		expect.Equal(t, err, nil, "error was not <nil>")

		resp, err := stream.Recv()
		statusErr := status.Convert(err)

		expect.Equal(t, resp, nilResponse, "responses are not equal")
		expect.Equal(t, statusErr.Code(), codes.Internal, "error codes are not equal")
		expect.Equal(t, statusErr.Message(), "very bad thing happened", "error messages are not equal")
	})

	t.Run("recovers and returns an error when the request causes a nil panic", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		stream, err := client.ServerStreamingEcho(context.Background(), nilPanicGRPCRequest)
		expect.Equal(t, err, nil, "error was not <nil>")

		resp, err := stream.Recv()
		statusErr := status.Convert(err)

		expect.Equal(t, resp, nilResponse, "responses are not equal")
		expect.Equal(t, statusErr.Code(), codes.Internal, "error codes are not equal")
		expect.Equal(t, statusErr.Message(), nilError, "error messages are not equal")
	})
}

func setupOverrideRecoveryServer(t *testing.T) (grpctest.Closer, echopb.EchoClient) {
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

		resp, err := client.Echo(context.Background(), goodGRPCRequest)

		expect.Equal(t, err, nil, "error was not <nil>")
		expect.Equal(t, resp.Message, goodGRPCRequest.Message, "messages are not equal")
	})

	t.Run("recovers and returns an error when the request panics", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		resp, err := client.Echo(context.Background(), panicGRPCRequest)
		statusErr := status.Convert(err)

		expect.Equal(t, resp, nilResponse, "responses are not equal")
		expect.Equal(t, statusErr.Code(), codes.Internal, "error codes are not equal")
		expect.Equal(t, statusErr.Message(), "very bad thing happened", "error messages are not equal")
	})

	t.Run("recovers and returns an error when the request causes a nil panic", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		resp, err := client.Echo(context.Background(), nilPanicGRPCRequest)
		statusErr := status.Convert(err)

		expect.Equal(t, resp, nilResponse, "responses are not equal")
		expect.Equal(t, statusErr.Code(), codes.Internal, "error codes are not equal")
		expect.Equal(t, statusErr.Message(), nilError, "error messages are not equal")
	})
}

func TestStreamServerInterceptor_Override(t *testing.T) {
	t.Parallel()

	t.Run("does not invoke the recovery function when the request is successful", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		stream, err := client.ServerStreamingEcho(context.Background(), goodGRPCRequest)
		expect.Equal(t, err, nil, "error was not <nil>")

		resp, err := stream.Recv()

		expect.Equal(t, err, nil, "error was not <nil>")
		expect.Equal(t, resp.Message, goodGRPCRequest.Message, "messages are not equal")
	})

	t.Run("recovers and returns an error when the request panics", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		stream, err := client.ServerStreamingEcho(context.Background(), panicGRPCRequest)
		expect.Equal(t, err, nil, "error was not <nil>")

		resp, err := stream.Recv()
		statusErr := status.Convert(err)

		expect.Equal(t, resp, nilResponse, "responses are not equal")
		expect.Equal(t, statusErr.Code(), codes.Internal, "error codes are no equal")
		expect.Equal(t, statusErr.Message(), "very bad thing happened", "error messages are not equal")
	})

	t.Run("recovers and returns an error when the request causes a nil panic", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		stream, err := client.ServerStreamingEcho(context.Background(), nilPanicGRPCRequest)
		expect.Equal(t, err, nil, "error was not <nil>")

		resp, err := stream.Recv()
		statusErr := status.Convert(err)

		expect.Equal(t, resp, nilResponse, "responses are not equal")
		expect.Equal(t, statusErr.Code(), codes.Internal, "error codes are no equal")
		expect.Equal(t, statusErr.Message(), nilError, "error messages are not equal")
	})
}
