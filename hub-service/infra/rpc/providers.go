package rpc

import (
	"fmt"
	"hub/config"
	"strings"
	"time"

	qhprorpc "common/rpc"
	"pb/clients"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func bindDependency[T any](
	target string,
	transport qhprorpc.TransportConfig,
	bind func(grpc.ClientConnInterface) T,
) (T, func(), error) {
	target = strings.TrimSpace(target)
	if target == "" {
		var zero T
		return zero, nil, fmt.Errorf("hub RPC target is required")
	}
	return qhprorpc.NewBoundClient(qhprorpc.ClientConfig{
		Target:          target,
		BackoffMaxDelay: 5 * time.Second,
		Credentials:     insecure.NewCredentials(),
		Transport:       transport,
	}, bind)
}

func NewBdsproRPCClient(runtime config.Runtime) (*clients.BdsproGrpcClient, func(), error) {
	return bindDependency(
		runtime.BDSProRPCTarget,
		runtime.Transport,
		clients.BindBdsproGrpcClient,
	)
}

func NewNotificationRPCClient(runtime config.Runtime) (*clients.NotificationClient, func(), error) {
	return bindDependency(
		runtime.NotificationRPCTarget,
		runtime.Transport,
		clients.BindNotificationClient,
	)
}
