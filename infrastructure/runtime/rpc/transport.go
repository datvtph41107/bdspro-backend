package rpc

import (
	"common/identity"
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// TransportConfig contains explicit cross-service assertion policy. It owns no
// service discovery, business authorization, quota, logging, or process exit.
type TransportConfig struct {
	ServiceAssertion        ServiceAssertionConfig
	RequireServiceAssertion bool
}

// ClientInterceptors is the canonical unary/stream outbound transport pair.
type ClientInterceptors struct {
	Unary  grpc.UnaryClientInterceptor
	Stream grpc.StreamClientInterceptor
}

// ServerInterceptors is the canonical unary/stream inbound transport pair.
type ServerInterceptors struct {
	Unary  grpc.UnaryServerInterceptor
	Stream grpc.StreamServerInterceptor
}

func ClientTransport(cfg TransportConfig) ClientInterceptors {
	return ClientInterceptors{
		Unary:  UnaryClientInterceptor(cfg),
		Stream: StreamClientInterceptor(cfg),
	}
}

func ServerTransport(cfg TransportConfig) ServerInterceptors {
	return ServerInterceptors{
		Unary:  UnaryServerInterceptor(cfg),
		Stream: StreamServerInterceptor(cfg),
	}
}

func UnaryClientInterceptor(cfg TransportConfig) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		outgoing, err := WithOutgoingContext(ctx, method, cfg)
		if err != nil {
			return err
		}
		return invoker(outgoing, method, req, reply, cc, opts...)
	}
}

func StreamClientInterceptor(cfg TransportConfig) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		outgoing, err := WithOutgoingContext(ctx, method, cfg)
		if err != nil {
			return nil, err
		}
		return streamer(outgoing, desc, cc, method, opts...)
	}
}

// WithOutgoingContext preserves existing outgoing metadata, replaces canonical
// request/identity fields from context, and signs the complete protected QHPRO
// transport contract. Audit mode permits an unsigned envelope only while the
// migration configuration is explicitly not configured; enforce mode rejects.
func WithOutgoingContext(
	ctx context.Context,
	method string,
	cfg TransportConfig,
) (context.Context, error) {
	if ctx == nil {
		return nil, status.Error(codes.Internal, "context is nil")
	}

	md := metadata.MD{}
	if existing, ok := metadata.FromOutgoingContext(ctx); ok {
		md = existing.Copy()
	}
	md = AppendRequestFromContext(ctx, md)
	md = AppendIdentityFromContext(ctx, md)

	if err := SignServiceAssertion(method, md, cfg.ServiceAssertion); err != nil {
		if errors.Is(err, ErrServiceAssertionNotConfigured) && !cfg.RequireServiceAssertion {
			return metadata.NewOutgoingContext(ctx, md), nil
		}
		if errors.Is(err, ErrServiceAssertionNotConfigured) {
			return nil, status.Error(codes.FailedPrecondition, "gRPC metadata signing is not configured")
		}
		return nil, status.Error(codes.Internal, "cannot sign gRPC metadata")
	}
	return metadata.NewOutgoingContext(ctx, md), nil
}

func UnaryServerInterceptor(cfg TransportConfig) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		method := ""
		if info != nil {
			method = info.FullMethod
		}
		bound, err := PrepareIncomingContext(ctx, method, cfg)
		if err != nil {
			return nil, err
		}
		return handler(bound, req)
	}
}

func StreamServerInterceptor(cfg TransportConfig) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if stream == nil {
			return status.Error(codes.Internal, "server stream is nil")
		}
		method := ""
		if info != nil {
			method = info.FullMethod
		}
		bound, err := PrepareIncomingContext(stream.Context(), method, cfg)
		if err != nil {
			return err
		}
		return handler(srv, &serverStreamWithContext{ServerStream: stream, ctx: bound})
	}
}

// PrepareIncomingContext verifies service-hop provenance before promoting
// Actor/Caller. Request tracing identities are restored independently so audit
// and unsigned operational RPCs retain observability without gaining privilege.
func PrepareIncomingContext(
	ctx context.Context,
	method string,
	cfg TransportConfig,
) (context.Context, error) {
	if ctx == nil {
		return nil, status.Error(codes.Internal, "context is nil")
	}
	md, _ := metadata.FromIncomingContext(ctx)
	if md == nil {
		md = metadata.MD{}
	}

	transportCaller, verifyErr := VerifyServiceAssertion(method, md, cfg.ServiceAssertion)

	verified := verifyErr == nil
	hasPrivilegedMetadata := HasPrivilegedMetadata(md)
	hasServiceAssertion := HasServiceAssertion(md)
	unsignedAllowed := AllowsUnsignedServiceCall(method) &&
		!hasPrivilegedMetadata &&
		!hasServiceAssertion

	switch {
	case verified:
		bound, err := identity.BindServiceCaller(ctx, transportCaller)
		if err != nil {
			return ctx, status.Error(codes.Internal, "transport caller context conflict")
		}
		ctx = bound

	case unsignedAllowed:
		// Explicit infrastructure exception. No verified service caller is promoted.

	case cfg.RequireServiceAssertion:
		return ctx, status.Error(codes.Unauthenticated, "service assertion required")

	case hasPrivilegedMetadata || hasServiceAssertion:
		// Verification failed. No verified service caller or identity is promoted.
	}

	var err error
	ctx, err = RestoreRequestFromMetadata(ctx, md)
	if err != nil {
		return ctx, err
	}
	ctx, err = restoreIdentityFromMetadata(ctx, md, verified, verified)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

// StrictIdentityUnaryServerInterceptor preserves the historical identitygrpc
// contract: the call must carry a valid signed transport assertion before
// Actor/Caller are restored. It exists only for compatibility facades.
func StrictIdentityUnaryServerInterceptor(cfg ServiceAssertionConfig) grpc.UnaryServerInterceptor {
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
		transportCaller, err := VerifyServiceAssertion(method, md, cfg)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "service assertion is required")
		}
		bound, err := identity.BindServiceCaller(ctx, transportCaller)
		if err != nil {
			return nil, status.Error(codes.Internal, "transport caller context conflict")
		}
		bound, err = RestoreIdentityFromMetadata(bound, md, true)
		if err != nil {
			return nil, err
		}
		return handler(bound, req)
	}
}

type serverStreamWithContext struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *serverStreamWithContext) Context() context.Context { return s.ctx }
