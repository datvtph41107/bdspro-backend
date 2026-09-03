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

func NewAuthRPCClient(
	userClient *clients.UserGrpcClient,
) (*clients.AuthGrpcClient, func(), error) {
	return bindDependency(
		"rpc.auth.address",
		func(conn grpc.ClientConnInterface) *clients.AuthGrpcClient {
			return clients.BindAuthGrpcClient(
				conn,
				userClient.RoleClient,
			)
		},
	)
}

func NewOrganizationRPCClient() (*clients.OrganizationClient, func(), error) {
	return bindDependency(
		"rpc.organization.address",
		clients.BindOrganizationClient,
	)
}

func NewHubRPCClient() (*clients.HubGrpcClient, func(), error) {
	return bindDependency(
		"rpc.hub.address",
		clients.BindHubGrpcClient,
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

func NewChatRPCClient() (commonprovider.ChatProvider, func(), error) {
	return bindDependency(
		"rpc.chat.address",
		func(conn grpc.ClientConnInterface) commonprovider.ChatProvider {
			return clients.BindChatGrpcClient(conn)
		},
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
