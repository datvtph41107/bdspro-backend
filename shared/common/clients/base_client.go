package clients

import (
	"common/pkg/config"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BaseGrpcClient struct {
	conn   *grpc.ClientConn
	config config.ServiceEndpoint
}

func NewBaseGrpcClient(ep config.ServiceEndpoint) (*BaseGrpcClient, error) {
	var lastErr error
	for i := 0; i <= ep.Retries; i++ {
		for _, addr := range ep.Addresses {
			ctx, cancel := context.WithTimeout(context.Background(), ep.Timeout)
			conn, err := grpc.DialContext(ctx, addr,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
			)
			cancel()
			if err == nil {
				return &BaseGrpcClient{conn: conn, config: ep}, nil
			}
			lastErr = err
		}
		if i < ep.Retries {
			time.Sleep(ep.Backoff)
		}
	}
	return nil, fmt.Errorf("failed to connect after %d retries: %w", ep.Retries, lastErr)
}

func (c *BaseGrpcClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *BaseGrpcClient) Conn() *grpc.ClientConn {
	return c.conn
}
