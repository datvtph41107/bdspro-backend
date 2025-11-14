package client

import (
	"context"
	"pb/clients"

	"auth/env"
	"auth/internal/dto"
	_enum "common/domain/enum"
	provider_grpc "common/provider/grpc"
	notificationpb "pb/types/notification"

	"google.golang.org/protobuf/types/known/structpb"
)

// @bind: auth/internal/interface/provider.NotificationProvider
type NotificationClient struct {
	*provider_grpc.NotificationGrpcClient
	*clients.NotificationClient
}

func NewNotificationClient() *NotificationClient {
	notificationClient := clients.MakeNotificationClient(env.ENV_RUNTIME)
	return &NotificationClient{
		NotificationClient:     notificationClient,
		NotificationGrpcClient: provider_grpc.MakeNotificationGrpcClient(env.ENV_RUNTIME),
	}
}

// LogAdminLogin ghi log khi admin đăng nhập thành công
func (c *NotificationClient) LogAdminLogin(ctx context.Context, adminID uint64, adminRole string, ipAddress string, userAgent string) error {
	return c.NotificationClient.LogAdminLogin(ctx, adminID, adminRole, ipAddress, userAgent)
}

// CreateAdminHistory tạo lịch sử admin
func (c *NotificationClient) CreateAdminHistory(ctx context.Context,
	adminID uint64,
	targetID uint64,
	targetType int32,
	actionType int32,
	title string,
	notes []string,
	preStage string,
	afterStage string,
	ownerID *uint64,
	ownerType int32,
	adminRole string,
	ipAddress string,
	userAgent string,
) error {
	return c.NotificationClient.CreateAdminHistory(ctx, adminID, targetID, targetType, actionType, title, notes, preStage, afterStage, ownerID, ownerType, adminRole, ipAddress, userAgent)
}

func (c *NotificationClient) CreateHistoryAuth(ctx context.Context, payload *dto.HistoryAuthCreateDTO) error {
	if payload == nil {
		return nil
	}

	var metadataStruct *structpb.Struct
	if payload.Metadata != nil {
		if converted, err := structpb.NewStruct(payload.Metadata); err == nil {
			metadataStruct = converted
		}
	}

	req := &notificationpb.HistoryAuthCreateRequest{
		UserId:         payload.UserID,
		OrganizationId: payload.OrganizationID,
		ActionType:     payload.ActionType,
		ActionName:     payload.ActionName,
		Description:    payload.Description,
		Reason:         payload.Reason,
		IpAddress:      payload.IPAddress,
		UserAgent:      payload.UserAgent,
		PerformedBy:    payload.PerformedBy,
		SessionId:      payload.SessionID,
		Channel:        payload.Channel,
		DeviceId:       payload.DeviceID,
		Location:       payload.Location,
		AdditionalNote: payload.AdditionalNote,
		SourceService:  payload.SourceService,
		Metadata:       metadataStruct,
	}

	if payload.Success != nil {
		req.Success = *payload.Success
	} else {
		req.Success = true
	}

	_, err := c.NotificationGrpcClient.CreateHistoryAuth(ctx, req)
	return err
}

// CreateNotification tạo thông báo cho user
// func (c *NotificationClient) CreateNotification(ctx context.Context,
// 	title string,
// 	message []string,
// 	notificationType int32,
// 	link string,
// 	ownerID uint64,
// 	attachData []string,
// ) error {
// 	return c.client.CreateNotification(ctx, title, message, notificationType, link, ownerID, attachData)
// }

func (c *NotificationClient) CreateNotification(
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
	return c.NotificationGrpcClient.CreateNotification(ctx, avatar, title, message, notificationType, targetId, ownerId, ownerOf, attachData)
}
