package health

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UpstreamChecker kiểm tra gRPC upstream dial được không
type UpstreamChecker struct {
	name     string
	addr     string
	critical bool
}

func (u *UpstreamChecker) Check(ctx context.Context) error {
	conn, err := grpc.DialContext(ctx, u.addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(2*time.Second),
	)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}
