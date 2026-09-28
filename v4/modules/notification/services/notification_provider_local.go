package services

import (
	"context"
	"time"

	_dto "common/domain/dto"
	_enum "common/domain/enum"
	_utils "common/utils"
	sharepb "pb/types/shared"
	shared_enum "pb/enums"
	"userv4/internal/dto"
	"userv4/modules/notification/entity"
	"userv4/modules/notification/repository"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// LocalNotificationProvider writes notification/history/config to the v4 DSN.
// It replaces remote notification-service calls for in-process producers.
type LocalNotificationProvider struct {
	repo repository.NotificationRepository
	sync *_utils.SyncUtil
}

func NewLocalNotificationProvider(repo repository.NotificationRepository, sync *_utils.SyncUtil) *LocalNotificationProvider {
	return &LocalNotificationProvider{repo: repo, sync: sync}
}

func (p *LocalNotificationProvider) CreateHistory(ctx context.Context, history *_dto.HistoryDTO) error {
	if history == nil {
		return nil
	}
	h := &entity.ActivityHistory{
		TargetID:   history.TargetId,
		TargetType: shared_enum.ETargetHistory(history.TargetType),
		ActionType: shared_enum.EHistory(history.ActionType),
		Title:      history.Title,
		Note:       pq.StringArray(history.Note),
		OwnerID:    history.OwnerID,
		OwnerType:  shared_enum.EOwnerType(history.OwnerOf),
	}
	return p.repo.CreateActivityHistory(ctx, h)
}

func (p *LocalNotificationProvider) CreateNotification(
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
	n := &entity.Notification{
		Avatar:           avatar,
		Title:            title,
		Message:          pq.StringArray(message),
		NotificationType: notificationType,
		OwnerID:          ownerId,
		OwnerOf:          ownerOf,
		AttachData:       pq.StringArray(attachData),
	}
	if targetId != nil {
		n.TargetID = *targetId
	}
	if err := p.repo.Create(ctx, n); err != nil {
		return err
	}
	if p.sync != nil {
		key := p.sync.GetKey(ctx, _utils.SyncKeyNotificationMe, 0)
		_ = p.sync.PutTimestamp(ctx, key, time.Now().UnixMilli())
	}
	return nil
}

func (p *LocalNotificationProvider) BatchPersonConfig(ctx context.Context, configs []*sharepb.PersonConfigV3Proto) error {
	rows := make([]*entity.PersonConfig, 0, len(configs))
	for _, cfg := range configs {
		if cfg == nil {
			continue
		}
		rows = append(rows, &entity.PersonConfig{
			ID:        cfg.Id,
			UserID:    cfg.UserId,
			Key:       cfg.Key,
			Checked:   cfg.Checked,
			IsDefault: cfg.IsDefault,
			Channel:   _enum.EChannelNotification(cfg.Channel),
		})
	}
	return p.repo.BatchUpsertPersonConfigs(ctx, rows)
}

func (p *LocalNotificationProvider) CreateAdminHistory(
	ctx context.Context,
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
	return p.repo.CreateAdminHistory(ctx, &entity.AdminHistory{
		AdminID:    adminID,
		TargetID:   targetID,
		TargetType: shared_enum.ETargetHistory(targetType),
		ActionType: shared_enum.EHistory(actionType),
		Title:      title,
		Note:       pq.StringArray(notes),
		PreStage:   preStage,
		AfterStage: afterStage,
		OwnerID:    ownerID,
		OwnerOf:    shared_enum.EOwnerType(ownerType),
		AdminRole:  adminRole,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	})
}

func (p *LocalNotificationProvider) CreateHistoryAuth(ctx context.Context, payload *dto.HistoryAuthCreateDTO) error {
	if payload == nil {
		return nil
	}
	success := true
	if payload.Success != nil {
		success = *payload.Success
	}
	source := payload.SourceService
	if source == "" {
		source = "auth-service"
	}
	meta := datatypes.JSONMap{}
	for k, v := range payload.Metadata {
		meta[k] = v
	}
	return p.repo.CreateHistoryAuth(ctx, &entity.HistoryAuth{
		UserID:         payload.UserID,
		OrganizationID: payload.OrganizationID,
		ActionType:     payload.ActionType,
		ActionName:     payload.ActionName,
		Description:    payload.Description,
		Success:        success,
		Reason:         payload.Reason,
		IPAddress:      payload.IPAddress,
		UserAgent:      payload.UserAgent,
		Metadata:       meta,
		PerformedBy:    payload.PerformedBy,
		SessionID:      payload.SessionID,
		Channel:        payload.Channel,
		DeviceID:       payload.DeviceID,
		Location:       payload.Location,
		AdditionalNote: payload.AdditionalNote,
		SourceService:  source,
	})
}

func (p *LocalNotificationProvider) GetPersonConfigs(ctx context.Context, userID uint64) ([]*dto.PersonConfigDTO, error) {
	rows, err := p.repo.ListPersonConfigs(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]*dto.PersonConfigDTO, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		out = append(out, &dto.PersonConfigDTO{
			ID:        row.ID,
			UserID:    row.UserID,
			Key:       row.Key,
			Checked:   row.Checked,
			IsDefault: row.IsDefault,
			Channel:   row.Channel,
		})
	}
	return out, nil
}
