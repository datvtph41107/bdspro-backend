package grpcclient

import (
	"common/identitygrpc"
	"common/requestgrpc"
	"common/transporttrust"
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

/**
 * UnaryClientInterceptor truyền canonical context sang service kế tiếp và ký metadata.
 */
func UnaryClientInterceptor(config transporttrust.Config) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		outgoing, err := WithOutgoingContext(ctx, method, config)
		if err != nil {
			return err
		}
		return invoker(outgoing, method, req, reply, cc, opts...)
	}
}

/**
 * StreamClientInterceptor áp dụng cùng metadata rules cho streaming gRPC call.
 */
func StreamClientInterceptor(config transporttrust.Config) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		outgoing, err := WithOutgoingContext(ctx, method, config)
		if err != nil {
			return nil, err
		}
		return streamer(outgoing, desc, cc, method, opts...)
	}
}

/**
 * WithOutgoingContext merge metadata đang có với request identity + Actor/Caller rồi ký lại.
 */
func WithOutgoingContext(
	ctx context.Context,
	method string,
	config transporttrust.Config,
) (context.Context, error) {
	md := metadata.MD{}
	if existing, ok := metadata.FromOutgoingContext(ctx); ok {
		md = existing.Copy()
	}

	md = requestgrpc.AppendFromContext(ctx, md)
	md = identitygrpc.AppendFromContext(ctx, md)

	if err := transporttrust.Sign(method, md, config); err != nil {
		if errors.Is(err, transporttrust.ErrNotConfigured) {
			return nil, status.Error(codes.FailedPrecondition, "gRPC metadata signing is not configured")
		}
		return nil, status.Error(codes.Internal, "cannot sign gRPC metadata")
	}

	return metadata.NewOutgoingContext(ctx, md), nil
}
