// Package grpcclient is the compatibility facade for the historical outgoing
// canonical metadata interceptor. New code should use common/rpc.
package grpcclient

import (
	"common/rpc"
	"common/transporttrust"
	"context"

	"google.golang.org/grpc"
)

// Deprecated: use rpc.UnaryClientInterceptor with explicit TransportConfig.
func UnaryClientInterceptor(config transporttrust.Config) grpc.UnaryClientInterceptor {
	return rpc.UnaryClientInterceptor(rpc.TransportConfig{Trust: config})
}

// Deprecated: use rpc.StreamClientInterceptor with explicit TransportConfig.
func StreamClientInterceptor(config transporttrust.Config) grpc.StreamClientInterceptor {
	return rpc.StreamClientInterceptor(rpc.TransportConfig{Trust: config})
}

// Deprecated: use rpc.WithOutgoingContext with explicit TransportConfig.
func WithOutgoingContext(
	ctx context.Context,
	method string,
	config transporttrust.Config,
) (context.Context, error) {
	return rpc.WithOutgoingContext(ctx, method, rpc.TransportConfig{Trust: config})
}
