package rpc

import (
	"time"

	"common/configloader"
	_provider "common/domain/provider"
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

func NewAuthRPCClient() (*clients.AuthGrpcClient, func(), error) {
	return bindDependency(
		"rpc.auth.address",
		func(conn grpc.ClientConnInterface) *clients.AuthGrpcClient {
			return clients.BindAuthGrpcClient(conn, nil)
		},
	)
}

func NewPaymentRPCClient() (*clients.PaymentClient, func(), error) {
	return bindDependency(
		"rpc.payment.address",
		clients.BindPaymentClient,
	)
}

func NewTQDRPCClient() (*clients.TQDGrpcClient, func(), error) {
	return bindDependency(
		"rpc.tqd.address",
		clients.BindTQDClient,
	)
}

func NewBdsproRPCClient() (*clients.BdsproGrpcClient, func(), error) {
	return bindDependency(
		"rpc.bdspro.address",
		clients.BindBdsproGrpcClient,
	)
}

func NewNotificationRPCClient() (*clients.NotificationClient, func(), error) {
	return bindDependency(
		"rpc.notification.address",
		clients.BindNotificationClient,
	)
}

func NewOrganizationRPCClient() (*clients.OrganizationClient, func(), error) {
	return bindDependency(
		"rpc.organization.address",
		clients.BindOrganizationClient,
	)
}

func NewChatRPCClient() (_provider.ChatProvider, func(), error) {
	client, cleanup, err := bindDependency(
		"rpc.chat.address",
		clients.BindChatGrpcClient,
	)
	if err != nil {
		return nil, nil, err
	}
	return client, cleanup, nil
}
