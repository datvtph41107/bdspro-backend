package rpc

import (
	"time"

	qhprorpc "common/rpc"
	"common/rpcenv"
	organizationconfig "organization/config"
	"pb/clients"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func bindDependency[T any](service string, bind func(grpc.ClientConnInterface) T) (T, func(), error) {
	target, err := organizationconfig.RPCAddress(service)
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
	return bindDependency("user", clients.BindUserGrpcClient)
}
func NewBdsproRPCClient() (*clients.BdsproGrpcClient, func(), error) {
	return bindDependency("bdspro", clients.BindBdsproGrpcClient)
}
func NewNotificationRPCClient() (*clients.NotificationClient, func(), error) {
	return bindDependency("notification", clients.BindNotificationClient)
}
func NewPaymentRPCClient() (*clients.PaymentClient, func(), error) {
	return bindDependency("payment", clients.BindPaymentClient)
}
func NewTransactionRPCClient() (*clients.TransactionGrpcClient, func(), error) {
	return bindDependency("transaction", clients.BindTransactionClient)
}
func NewChatRPCClient() (*clients.ChatGrpcClient, func(), error) {
	return bindDependency("chat", clients.BindChatGrpcClient)
}
func NewAuthRPCClient() (*clients.AuthGrpcClient, func(), error) {
	return bindDependency("auth", func(conn grpc.ClientConnInterface) *clients.AuthGrpcClient {
		return clients.BindAuthGrpcClient(conn, nil)
	})
}
