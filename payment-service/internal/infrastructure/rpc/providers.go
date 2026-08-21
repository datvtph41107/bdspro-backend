package rpc

import (
	"fmt"
	"time"

	"common/rpc"
	"common/transporttrust"
	"payment/appconfig"
	"pb/clients"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ProvideNotificationRPCClient is explicitly wired by the payment process.
// It is intentionally not named New* because payment has a custom nested Wire graph.
func ProvideNotificationRPCClient() (*clients.NotificationClient, func(), error) {
	target := fmt.Sprintf("%s:%s", appconfig.NOTIFICATION_SERVICE_HOST, appconfig.NOTIFICATION_SERVICE_PORT)
	return rpc.NewBoundClient(rpc.ClientConfig{
		Target:          target,
		BackoffMaxDelay: 5 * time.Second,
		Credentials:     insecure.NewCredentials(),
		Transport:       rpc.TransportConfig{Trust: transporttrust.ConfigFromEnv()},
	}, func(conn grpc.ClientConnInterface) *clients.NotificationClient {
		return clients.BindNotificationClient(conn)
	})
}
