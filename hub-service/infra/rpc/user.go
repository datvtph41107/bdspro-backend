package rpc

import (
	"fmt"
	"hub/config"
	"strings"
	"time"

	qhprorpc "common/rpc"
	userpb "pb/types/user"

	"google.golang.org/grpc/credentials/insecure"
)

// NewAdminUserProfileClient constructs the Hub process-owned User RPC
// capability. Wire propagates cleanup to the process composition root.
func NewAdminUserProfileClient(
	runtime config.Runtime,
) (
	userpb.AdminUserProfileServiceClient,
	func(),
	error,
) {
	target := strings.TrimSpace(runtime.UserRPCTarget)
	if target == "" {
		return nil, nil, fmt.Errorf("hub user RPC target is required")
	}

	conn, err := qhprorpc.NewClient(qhprorpc.ClientConfig{
		Target:          target,
		BackoffMaxDelay: 5 * time.Second,
		Credentials:     insecure.NewCredentials(),
		Transport:       runtime.Transport,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create hub user gRPC connection: %w", err)
	}

	cleanup := func() {
		_ = conn.Close()
	}

	return userpb.NewAdminUserProfileServiceClient(conn), cleanup, nil
}
