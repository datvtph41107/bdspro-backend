package client

import (
	_enum "common/domain/enum"
	"context"
	"errors"
	"pb/clients"
	notificationpb "pb/types/notification"
	sharepb "pb/types/shared"
	"time"
	"user/internal/dto"
	"user/internal/interface/providers"

	"google.golang.org/protobuf/types/known/structpb"
)

type NotificationClient struct {
	*clients.NotificationClient
}

func NewNotificationClient(
	rpcClient *clients.NotificationClient,
) providers.NotificationProvider {
	return &NotificationClient{
		NotificationClient: rpcClient,
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

	_, err := c.NotificationClient.CreateHistoryAuth(ctx, req)
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
	if c == nil || c.NotificationClient == nil || c.NotificationClient.Client == nil {
		return errors.New("notification grpc client chưa được khởi tạo")
	}
	callCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return c.NotificationClient.CreateNotification(callCtx, avatar, title, message, notificationType, targetId, ownerId, ownerOf, attachData)
}

func (c *NotificationClient) GetPersonConfigs(ctx context.Context, userID uint64) ([]*dto.PersonConfigDTO, error) {
	if userID == 0 {
		return []*dto.PersonConfigDTO{}, nil
	}

	if c.NotificationClient == nil || c.NotificationClient.Client == nil {
		return nil, errors.New("notification grpc client chưa được khởi tạo")
	}

	var configs []*sharepb.PersonConfigV3Proto
	var err error

	if configs, err = c.NotificationClient.GetPersonConfigs(ctx, userID); err != nil {
		return nil, err
	}

	result := make([]*dto.PersonConfigDTO, 0, len(configs))
	for _, cfg := range configs {
		if cfg == nil {
			continue
		}
		result = append(result, &dto.PersonConfigDTO{
			ID:        cfg.Id,
			UserID:    cfg.UserId,
			Key:       cfg.Key,
			Checked:   cfg.Checked,
			IsDefault: cfg.IsDefault,
			Channel:   _enum.EChannelNotification(cfg.Channel),
		})
	}

	return result, nil
}
