package client

import (
	"context"
)

type NotificationClient struct {
}

func (s *NotificationClient) NotifyToUser(ctx context.Context, userId uint64, message string, newsFeedID uint64) error {
	return nil
}
