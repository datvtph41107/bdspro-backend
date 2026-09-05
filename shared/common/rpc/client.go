package rpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

// ClientConfig contains channel-construction mechanics only. The composition
// owner supplies topology, credentials, service-assertion policy, and Close ownership.
type ClientConfig struct {
	Target             string
	BackoffMaxDelay    time.Duration
	Credentials        credentials.TransportCredentials
	Transport          TransportConfig
	DefaultCallOptions []grpc.CallOption
}

// NewClient builds the canonical QHPRO gRPC channel. It does not wait for
// READY, start monitoring goroutines, or own process lifecycle. Caller owns
// Close on the returned connection.
func NewClient(cfg ClientConfig) (*grpc.ClientConn, error) {
	return newClient(cfg, ClientTransport(cfg.Transport))
}

// NewClientWithInterceptors is a migration-only compatibility surface for the
// historical transport bridge. New process composition must use NewClient.
func NewClientWithInterceptors(
	cfg ClientConfig,
	interceptors ClientInterceptors,
) (*grpc.ClientConn, error) {
	return newClient(cfg, interceptors)
}

func newClient(cfg ClientConfig, interceptors ClientInterceptors) (*grpc.ClientConn, error) {
	target := strings.TrimSpace(cfg.Target)
	if target == "" {
		return nil, fmt.Errorf("gRPC target is required")
	}
	if cfg.BackoffMaxDelay <= 0 {
		return nil, fmt.Errorf("gRPC backoff max delay must be positive")
	}
	if cfg.Credentials == nil {
		return nil, fmt.Errorf("gRPC transport credentials are required")
	}
	if interceptors.Unary == nil || interceptors.Stream == nil {
		return nil, fmt.Errorf("gRPC unary and stream interceptors are required")
	}

	backoffConfig := backoff.DefaultConfig
	backoffConfig.MaxDelay = cfg.BackoffMaxDelay

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(cfg.Credentials),
		grpc.WithConnectParams(grpc.ConnectParams{Backoff: backoffConfig}),
		grpc.WithUnaryInterceptor(interceptors.Unary),
		grpc.WithStreamInterceptor(interceptors.Stream),
	}
	if len(cfg.DefaultCallOptions) > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(cfg.DefaultCallOptions...))
	}

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("create gRPC client channel %q: %w", target, err)
	}
	return conn, nil
}

// Unavailable returns a connection interface that fails every RPC with
// codes.Unavailable. Compatibility provider graphs use it while constructor
// errors cannot yet be propagated through generated Wire code.
func Unavailable(cause error) grpc.ClientConnInterface {
	message := "gRPC client unavailable"
	if cause != nil {
		message = fmt.Sprintf("%s: %v", message, cause)
	}
	return unavailableConn{message: message}
}

type unavailableConn struct{ message string }

func (c unavailableConn) Invoke(
	context.Context, string, interface{}, interface{}, ...grpc.CallOption,
) error {
	return status.Error(codes.Unavailable, c.message)
}

func (c unavailableConn) NewStream(
	context.Context, *grpc.StreamDesc, string, ...grpc.CallOption,
) (grpc.ClientStream, error) {
	return nil, status.Error(codes.Unavailable, c.message)
}
