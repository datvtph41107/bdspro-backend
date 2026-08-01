package injection

import (
	_provider "common/domain/provider"
	"os"
	"pb/clients"
)

var runtimeEnv = os.Getenv("ENV_RUNTIME")

func NewAuthClient() *clients.AuthGrpcClient {
	return clients.MakeAuthGrpcClient(runtimeEnv)
}

func NewBDSProClient() *clients.BdsproGrpcClient {
	return clients.MakeBdsproGrpcClient(runtimeEnv)
}

func NewCRMClient() _provider.CrmProvider {
	return clients.MakeCrmGrpcClient(runtimeEnv)
}

func NewChatClient() _provider.ChatProvider {
	return clients.NewChatGrpcClient(runtimeEnv)
}

func NewHubClient() *clients.HubGrpcClient {
	return clients.MakeHubGrpcClient(runtimeEnv)
}

func NewNotification() *clients.NotificationClient {
	return clients.MakeNotificationClient(runtimeEnv)
}

func NewOrganizationClient() *clients.OrganizationClient {
	return clients.MakeOrganizationClient(runtimeEnv)
}

func NewPaymentClient() *clients.PaymentClient {
	return clients.MakePaymentClient(runtimeEnv)
}

func NewTQClient() *clients.TQDGrpcClient {
	return clients.MakeTQDClient(runtimeEnv)
}

func NewTransactionClient() *clients.TransactionGrpcClient {
	return clients.NewTransactionGrpcClient(runtimeEnv)
}

func NewUserClient() *clients.UserGrpcClient {
	return clients.MakeUserGrpcClient(runtimeEnv)
}
