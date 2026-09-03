package rpc

import (
	"time"

	"common/configloader"
	qhprorpc "common/rpc"
	"common/rpcenv"
	"pb/clients"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func bindDependency[T any](
	key string,
	bind func(grpc.ClientConnInterface) T,
) (T, func(), error) {
	target, err := configloader.RequiredString(key)
	if err != nil {
		var zero T
		return zero, nil, err
	}

	return qhprorpc.NewBoundClient(qhprorpc.ClientConfig{
		Target:          target,
		BackoffMaxDelay: 5 * time.Second,
		Credentials:     insecure.NewCredentials(),
		Transport:       rpcenv.LoadTransportConfig(),
	}, bind)
}

func NewUserRPCClient() (*clients.UserGrpcClient, func(), error) {
	return bindDependency(
		"rpc.user.address",
		clients.BindUserGrpcClient,
	)
}

func NewBdsproRPCClient() (*clients.BdsproGrpcClient, func(), error) {
	return bindDependency(
		"rpc.bdspro.address",
		clients.BindBdsproGrpcClient,
	)
}
