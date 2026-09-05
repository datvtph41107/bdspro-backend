package rpc

import (
	"time"

	"common/configloader"
	commonprovider "common/domain/provider"
	qhprorpc "common/rpc"
	"common/rpcenv"
	"pb/clients"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func bindDependency[T any](key string, bind func(grpc.ClientConnInterface) T) (T, func(), error) {
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
	return bindDependency("rpc.user.address", clients.BindUserGrpcClient)
}

func NewAuthRPCClient() (*clients.AuthGrpcClient, func(), error) {
	return bindDependency(
		"rpc.auth.address",
		func(conn grpc.ClientConnInterface) *clients.AuthGrpcClient {
			return clients.BindAuthGrpcClient(conn, nil)
		},
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

func NewPaymentRPCClient() (*clients.PaymentClient, func(), error) {
	return bindDependency(
		"rpc.payment.address",
		clients.BindPaymentClient,
	)
}

func NewHubRPCClient() (*clients.HubGrpcClient, func(), error) {
	return bindDependency(
		"rpc.hub.address",
		clients.BindHubGrpcClient,
	)
}

func NewChatRPCClient() (commonprovider.ChatProvider, func(), error) {
	return bindDependency(
		"rpc.chat.address",
		func(conn grpc.ClientConnInterface) commonprovider.ChatProvider {
			return clients.BindChatGrpcClient(conn)
		},
	)
}

func NewTQDRPCClient() (*clients.TQDGrpcClient, func(), error) {
	return bindDependency(
		"rpc.tqd.address",
		clients.BindTQDClient,
	)
}

func NewCRMRPCClient() (commonprovider.CrmProvider, func(), error) {
	return bindDependency(
		"rpc.crm.address",
		func(conn grpc.ClientConnInterface) commonprovider.CrmProvider {
			return clients.BindCrmGrpcClient(conn)
		},
	)
}
