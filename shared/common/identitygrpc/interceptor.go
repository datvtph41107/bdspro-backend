package identitygrpc

import (
	"common/transporttrust"
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

/**
 * UnaryServerInterceptor verify metadata signature rồi mới restore Actor/Caller.
 */
func UnaryServerInterceptor(config transporttrust.Config) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		method := ""
		if info != nil {
			method = info.FullMethod
		}

		transportCaller, err := transporttrust.Verify(method, md, config)
		if err != nil {
			if errors.Is(err, transporttrust.ErrMissingAssertion) ||
				errors.Is(err, transporttrust.ErrInvalidAssertion) ||
				errors.Is(err, transporttrust.ErrExpiredAssertion) ||
				errors.Is(err, transporttrust.ErrNotConfigured) {
				return nil, status.Error(codes.Unauthenticated, "trusted gRPC metadata is required")
			}
			return nil, status.Error(codes.Unauthenticated, "trusted gRPC metadata is invalid")
		}

		bound, err := transporttrust.BindCaller(ctx, transportCaller)
		if err != nil {
			return nil, status.Error(codes.Internal, "transport caller context conflict")
		}

		bound, err = RestoreFromMetadata(bound, md)
		if err != nil {
			return nil, err
		}

		return handler(bound, req)
	}
}
