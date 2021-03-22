package recovery

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a unary server interceptor for panic
// recovery.
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := applyOptions(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		panicked := true

		defer func() {
			if p := recover(); p != nil || panicked {
				err = recoverRPC(ctx, p, o.recovery)
			}
		}()

		// If the call to the handler does not result in a panic then the code
		// following it will be executed. Therefore if the handler panics then
		// `panicked` will not be set to false. This is essential to be able to
		// detect nil panics from the handler.
		resp, err := handler(ctx, req)
		panicked = false
		return resp, err
	}
}

// StreamServerInterceptor returns a streaming server interceptor for panic
// recovery.
func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	o := applyOptions(opts)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		panicked := true

		defer func() {
			if p := recover(); p != nil || panicked {
				err = recoverRPC(stream.Context(), p, o.recovery)
			}
		}()

		// If the call to the handler does not result in a panic then the code
		// following it will be executed. Therefore if the handler panics then
		// `panicked` will not be set to false. This is essential to be able to
		// detect nil panics from the handler.
		err = handler(srv, stream)
		panicked = false
		return err
	}
}

func recoverRPC(ctx context.Context, p interface{}, fn RecoveryContextFunc) error {
	err := fn(ctx, p)

	if err == nil {
		return status.Errorf(codes.Internal, "%v", p)
	}

	return err
}
