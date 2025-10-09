package client

import (
	"user/env"

	"pb/clients"
)

// @bind: user/internal/interface/providers.NotificationProvider
type NotificationClient struct {
	*clients.NotificationClient
}

func NewNotificationClient() *NotificationClient {
	return &NotificationClient{
		NotificationClient: clients.MakeNotificationClient(env.ENV_RUNTIME),
	}
}
