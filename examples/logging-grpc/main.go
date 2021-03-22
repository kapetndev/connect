package main

import (
	"context"
	"encoding/hex"
	"net"
	"time"

	"github.com/sirupsen/logrus"

	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/crumbandbase/service-core-go/logging"
	logging_logrus "github.com/crumbandbase/service-core-go/logging/logrus"
	echopb "github.com/crumbandbase/service-core-go/testdata/proto"
)

const timeout = 10 * time.Second

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	grpclogger := logging_logrus.New(logger)
	grpclog.SetLoggerV2(grpclogger)

	s := grpc.NewServer(
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
		grpc.ConnectionTimeout(timeout),
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(
				grpclogger,
				logging.WithField(traceFunc),
			),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(
				grpclogger,
				logging.WithField(traceFunc),
			),
		),
	)

	echopb.RegisterEchoServer(s, &server{})
	healthpb.RegisterHealthServer(s, &health.Server{})

	// Register reflection service on grpc server for debugging.
	reflection.Register(s)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	logger.Info("server started on [::]:50051")
	if err := s.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}

type server struct {
	echopb.UnimplementedEchoServer
}

// grpcurl -plaintext -d '{"message": "Hello, world"}' :50051 echo.v1.Echo/Echo | jq -r .message
func (*server) Echo(ctx context.Context, in *echopb.EchoRequest) (*echopb.EchoResponse, error) {
	l := logging.FromContext(ctx)
	l.Info("logging from request handler")

	return &echopb.EchoResponse{
		Message: in.Message,
	}, nil
}

func traceFunc(ctx context.Context) (string, interface{}) {
	span := trace.FromContext(ctx)
	id := span.SpanContext().TraceID

	return logging.TraceKey, hex.EncodeToString(id[:])
}
