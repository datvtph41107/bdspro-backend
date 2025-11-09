package client

import (
	_enum "common/domain/enum"
	provider_grpc "common/provider/grpc"
	"context"
	"fmt"
	sharepb "pb/types/shared"
	"social/env"
)

// @bind: social/internal/interface.NotiClient
type NotificationClient struct {
	*provider_grpc.NotificationGrpcClient
}

func NewNotificationClient() *NotificationClient {
	return &NotificationClient{
		NotificationGrpcClient: provider_grpc.MakeNotificationGrpcClient(env.ENV_RUNTIME),
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

func (s *NotificationClient) CreateToOwner(
	ctx context.Context,
	avatar string,
	title string,
	message []string,
	notificationType _enum.ENotificationType,
	targetId *uint64,
	ownerId uint64,
	ownerOf _enum.EOwnerOf,
	attachData []string,
) error {
	_, err := s.Client.CreateToOwner(ctx, &sharepb.NotificationRequest{
		Avatar:           avatar,
		Title:            title,
		Message:          message,
		NotificationType: uint32(notificationType),
		OwnerId:          ownerId,
		TargetId:         *targetId,
		OwnerOf:          uint32(ownerOf),
		AttachData:       attachData,
		SendToDevice:     true,
	})
	return err
}
