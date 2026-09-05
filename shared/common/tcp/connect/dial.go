package connect

import (
	_middleware "common/middleware"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// DefaultDialOptions trả về DialOption chuẩn: insecure + ConnectParams (auto-reconnect)
// + keepalive + middleware interceptor của common.
func DefaultDialOptions(extra ...grpc.DialOption) []grpc.DialOption {
	cfg := DefaultDialConfig()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  cfg.BaseDelay,
				Multiplier: 1.6,
				Jitter:     0.2,
				MaxDelay:   cfg.MaxDelay,
			},
			MinConnectTimeout: cfg.MinConnectTimeout,
		}),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                cfg.KeepaliveTime,
			Timeout:             cfg.KeepaliveTimeout,
			PermitWithoutStream: true,
		}),
		grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(_middleware.StreamClientInterceptor),
	}
	return append(opts, extra...)
}

// Dial tạo ClientConn tới target. Không block; gRPC tự reconnect khi disconnect
// theo ConnectParams.
func Dial(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if target == "" {
		return nil, fmt.Errorf("tcp/connect: empty target")
	}
	if len(opts) == 0 {
		opts = DefaultDialOptions()
	}
	return grpc.NewClient(target, opts...)
}

// DialNamed dial theo tên service và cache trong Manager mặc định (reuse conn).
func DialNamed(name, target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return DefaultManager().GetOrDial(name, target, opts...)
}

// AddressOf ghép host từ env_runtime + port.
func AddressOf(envRuntime, name, port string) string {
	host := HostOf(envRuntime, name)
	if port == "" {
		port = DefaultPort(name)
	}
	if port == "" {
		return host
	}
	if len(port) > 0 && port[0] == ':' {
		return host + port
	}
	return fmt.Sprintf("%s:%s", host, port)
}

// HostOf trả hostname service theo ENV_RUNTIME (docker name vs localhost).
func HostOf(envRuntime, name string) string {
	if envRuntime != "" {
		return name
	}
	return "localhost"
}

// WaitReady chờ conn chuyển sang READY hoặc context hết hạn.
func WaitReady(ctx context.Context, conn *grpc.ClientConn) error {
	if conn == nil {
		return fmt.Errorf("tcp/connect: nil connection")
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	for {
		state := conn.GetState()
		if state == connectivity.Ready {
			return nil
		}
		if state == connectivity.Shutdown {
			return fmt.Errorf("tcp/connect: connection shutdown")
		}
		if state == connectivity.Idle {
			conn.Connect()
		}
		if !conn.WaitForStateChange(ctx, state) {
			return fmt.Errorf("tcp/connect: wait ready timeout (state=%v): %w", conn.GetState(), ctx.Err())
		}
	}
}
