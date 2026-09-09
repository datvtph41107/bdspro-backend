package rpc

import (
	"fmt"
	"strings"
	"time"

	qhprorpc "common/rpc"
	"pb/clients"

	"google.golang.org/grpc/credentials/insecure"
)

func NewUserRPCClient(
	target string,
	transport qhprorpc.TransportConfig,
) (*clients.UserGrpcClient, func(), error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return nil, nil, fmt.Errorf("auth user RPC target is required")
	}

	return qhprorpc.NewBoundClient(qhprorpc.ClientConfig{
		Target:          target,
		BackoffMaxDelay: 5 * time.Second,
		Credentials:     insecure.NewCredentials(),
		Transport:       transport,
	}, clients.BindUserGrpcClient)
}
