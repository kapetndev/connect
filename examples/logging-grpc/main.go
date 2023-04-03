package main

import (
	"context"
	"net"
	"os"
	"time"

	"golang.org/x/exp/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/kapetndev/connect/logging"
	echopb "github.com/kapetndev/connect/testdata/echo/v1"
)

const timeout = 10 * time.Second

func main() {
	h := logging.NewGoogleCloudHandler(os.Stdout, slog.LevelDebug)

	s := grpc.NewServer(
		grpc.ConnectionTimeout(timeout),
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(logging.WithHandler(h)),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(logging.WithHandler(h)),
		),
	)

	echopb.RegisterEchoServiceServer(s, &server{})
	healthpb.RegisterHealthServer(s, &health.Server{})

	// Register reflection service on gRPC server for debugging.
	reflection.Register(s)

	logger := slog.New(h)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("failed to listen: " + err.Error())
		os.Exit(1)
	}

	logger.Info("server started on [::]:50051")
	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve: " + err.Error())
	}
}

type server struct {
	echopb.UnimplementedEchoServiceServer
}

func (*server) Echo(ctx context.Context, in *echopb.EchoRequest) (*echopb.EchoResponse, error) {
	l := logging.FromContext(ctx)
	l.Info(ctx, "logging from request handler")

	return &echopb.EchoResponse{
		Message: in.Message,
	}, nil
}
