package client

import (
	"context"
	"fmt"
	clients "pb/clients"
	sharepb "pb/types/shared"
	"social/env"
)

// @bind: social/internal/interface.NotiClient
type NotificationClient struct {
	*clients.NotificationClient
}

func NewNotificationClient() *NotificationClient {
	return &NotificationClient{
		NotificationClient: clients.MakeNotificationClient(env.ENV_RUNTIME),
	}
}

func (s *NotificationClient) NotifyToUser(
	ctx context.Context,
	typeNoti int32,
	userId *uint64,
	attachData []string,
	targetID uint64,
	isMerge bool,
) error {
	_, err := s.Client.Create(ctx, &sharepb.NotificationRequest{
		AttachData: attachData,
		Link:       fmt.Sprintf("/news-feed/%d", targetID),
		Type:       typeNoti,
		OwnerId:    *userId,
		TargetId:   targetID,
		TargetType: typeNoti,
		IsMerge:    isMerge,
	})
	return err
}
