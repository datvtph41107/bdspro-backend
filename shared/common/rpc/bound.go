package rpc

import "google.golang.org/grpc"

// NewBoundClient constructs one process-owned channel and binds a typed RPC
// capability to it. The caller owns the returned cleanup function.
func NewBoundClient[T any](
	cfg ClientConfig,
	bind func(grpc.ClientConnInterface) T,
) (T, func(), error) {
	var zero T

	conn, err := NewClient(cfg)
	if err != nil {
		return zero, nil, err
	}

	cleanup := func() {
		_ = conn.Close()
	}

	return bind(conn), cleanup, nil
}
