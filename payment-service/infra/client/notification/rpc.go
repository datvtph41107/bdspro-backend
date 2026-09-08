package notificationclient

import (
	"fmt"
	"strings"
	"time"

	qhprorpc "common/rpc"
	"common/rpcenv"
	"pb/clients"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ProvideNotificationRPCClientAt binds the legacy Wallet notification adapter
// to a process-owned target. Config/env ownership remains in cmd composition.
func ProvideNotificationRPCClientAt(target string) (*clients.NotificationClient, func(), error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return nil, nil, fmt.Errorf("notification gRPC address is required")
	}
	return qhprorpc.NewBoundClient(qhprorpc.ClientConfig{
		Target:          target,
		BackoffMaxDelay: 5 * time.Second,
		Credentials:     insecure.NewCredentials(),
		Transport:       rpcenv.LoadTransportConfig(),
	}, func(conn grpc.ClientConnInterface) *clients.NotificationClient {
		return clients.BindNotificationClient(conn)
	})
}
