package rpc

import (
	"time"

	"common/configloader"
	qhprorpc "common/rpc"
	"common/rpcenv"
	"pb/clients"

	"google.golang.org/grpc/credentials/insecure"
)

func NewUserRPCClient() (*clients.UserGrpcClient, func(), error) {
	target, err := configloader.RequiredString("rpc.user.address")
	if err != nil {
		return nil, nil, err
	}

	return qhprorpc.NewBoundClient(qhprorpc.ClientConfig{
		Target:          target,
		BackoffMaxDelay: 5 * time.Second,
		Credentials:     insecure.NewCredentials(),
		Transport:       rpcenv.LoadTransportConfig(),
	}, clients.BindUserGrpcClient)
}
