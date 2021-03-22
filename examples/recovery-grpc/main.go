package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/crumbandbase/service-core-go/recovery"
	echopb "github.com/crumbandbase/service-core-go/testdata/proto"
)

const timeout = 10 * time.Second

func main() {
	recoveryFn := func(ctx context.Context, p interface{}) error {
		return status.Errorf(codes.Internal, "recovered from panic: %v", p)
	}

	s := grpc.NewServer(
		grpc.ConnectionTimeout(timeout),
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(
				recovery.WithRecoveryContext(recoveryFn),
			),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(
				recovery.WithRecoveryContext(recoveryFn),
			),
		),
	)

	echopb.RegisterEchoServer(s, &server{})
	healthpb.RegisterHealthServer(s, &health.Server{})

	// Register reflection service on grpc server for debugging.
	reflection.Register(s)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("server started on [::]:50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct {
	echopb.UnimplementedEchoServer
}

// grpcurl -plaintext -d '{"message": "Hello, world"}' :50051 echo.v1.Echo/Echo | jq -r .message
func (*server) Echo(ctx context.Context, in *echopb.EchoRequest) (*echopb.EchoResponse, error) {
	panic("shut up, Wesley!")

	return &echopb.EchoResponse{
		Message: in.Message,
	}, nil
}
