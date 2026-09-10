package injection

import (
	"time"

	"common/configloader"
	_provider "common/domain/provider"
	qhprorpc "common/rpc"
	"common/rpcenv"
	"pb/clients"

	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RPCClientSet is the shared process-owned gRPC client graph.
// Services include this set in wire instead of redefining dial/bind per service.
var RPCClientSet = wire.NewSet(
	NewAuthClient,
	NewBDSProClient,
	NewCRMClient,
	NewChatClient,
	NewHubClient,
	NewNotification,
	NewOrganizationClient,
	NewPaymentClient,
	NewTQClient,
	NewTransactionClient,
	NewUserClient,
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

// NewAuthClient opens rpc.auth.address and attaches RoleClient from user when present.
// RoleService is hosted by user-service; AuthInternalService is hosted by auth-service.
func NewAuthClient(userClient *clients.UserGrpcClient) (*clients.AuthGrpcClient, func(), error) {
	return bindDependency(
		"rpc.auth.address",
		func(conn grpc.ClientConnInterface) *clients.AuthGrpcClient {
			client := clients.BindAuthGrpcClient(conn)
			if userClient != nil {
				client.RoleClient = userClient.RoleClient
			}
			return client
		},
	)
}

func NewBDSProClient() (*clients.BdsproGrpcClient, func(), error) {
	return bindDependency("rpc.bdspro.address", clients.BindBdsproGrpcClient)
}

func NewCRMClient() (_provider.CrmProvider, func(), error) {
	return bindDependency(
		"rpc.crm.address",
		func(conn grpc.ClientConnInterface) _provider.CrmProvider {
			return clients.BindCrmGrpcClient(conn)
		},
	)
}

func NewChatClient() (_provider.ChatProvider, func(), error) {
	return bindDependency(
		"rpc.chat.address",
		func(conn grpc.ClientConnInterface) _provider.ChatProvider {
			return clients.BindChatGrpcClient(conn)
		},
	)
}

func NewHubClient() (*clients.HubGrpcClient, func(), error) {
	return bindDependency("rpc.hub.address", clients.BindHubGrpcClient)
}

func NewNotification() (*clients.NotificationClient, func(), error) {
	return bindDependency("rpc.notification.address", clients.BindNotificationClient)
}

func NewOrganizationClient() (*clients.OrganizationClient, func(), error) {
	return bindDependency("rpc.organization.address", clients.BindOrganizationClient)
}

func NewPaymentClient() (*clients.PaymentClient, func(), error) {
	return bindDependency("rpc.payment.address", clients.BindPaymentClient)
}

func NewTQClient() (*clients.TQDGrpcClient, func(), error) {
	return bindDependency("rpc.tqd.address", clients.BindTQDClient)
}

func NewTransactionClient() (*clients.TransactionGrpcClient, func(), error) {
	return bindDependency("rpc.transaction.address", clients.BindTransactionClient)
}

func NewUserClient() (*clients.UserGrpcClient, func(), error) {
	return bindDependency("rpc.user.address", clients.BindUserGrpcClient)
}
