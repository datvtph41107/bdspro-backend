package provider

import "context"

// IFirebaseProvider interface cho Firebase provider
type IFirebaseProvider interface {
	// SendPushNotification gửi push notification đến nhiều tokens
	SendPushNotification(ctx context.Context, tokens []string, title string, message string, data map[string]string) error
}
