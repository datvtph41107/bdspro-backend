package rpc

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

func OpenUserRPCClient(target string) (*clients.UserGrpcClient, func(), error) {
	target = strings.TrimSpace(target)
	if target == "" { return nil, nil, fmt.Errorf("notification user gRPC target is required") }
	return qhprorpc.NewBoundClient(qhprorpc.ClientConfig{
		Target: target,
		BackoffMaxDelay: 5 * time.Second,
		Credentials: insecure.NewCredentials(),
		Transport: rpcenv.LoadTransportConfig(),
	}, func(conn grpc.ClientConnInterface) *clients.UserGrpcClient { return clients.BindUserGrpcClient(conn) })
}
