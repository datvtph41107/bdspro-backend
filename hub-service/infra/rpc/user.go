package rpc

import (
	"fmt"
	"time"

	qhprorpc "common/rpc"
	"common/rpcenv"
	hubconfig "hub/config"
	userpb "pb/types/user"

	"google.golang.org/grpc/credentials/insecure"
)

// NewAdminUserProfileClient constructs the Hub process-owned User RPC
// capability. Wire propagates cleanup to the process composition root.
func NewAdminUserProfileClient() (
	userpb.AdminUserProfileServiceClient,
	func(),
	error,
) {
	conn, err := qhprorpc.NewClient(qhprorpc.ClientConfig{
		Target:          hubconfig.AppProperties.RPC.User.Address,
		BackoffMaxDelay: 5 * time.Second,
		Credentials:     insecure.NewCredentials(),
		Transport:       rpcenv.LoadTransportConfig(),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create hub user gRPC connection: %w", err)
	}

	cleanup := func() {
		_ = conn.Close()
	}

	return userpb.NewAdminUserProfileServiceClient(conn), cleanup, nil
}
