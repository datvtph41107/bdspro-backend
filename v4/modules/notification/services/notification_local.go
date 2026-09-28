package services

import (
	"context"
	"fmt"
	"time"

	_utils "common/utils"
	notificationpb "pb/types/notification"
	"userv4/modules/notification/dto"
	"userv4/modules/notification/entity"
	"userv4/modules/notification/repository"
)

// LocalNotificationService owns inbox flows against the v4 Postgres DSN.
type LocalNotificationService struct {
	repo repository.NotificationRepository
	sync *_utils.SyncUtil
}

func NewLocalNotificationService(repo repository.NotificationRepository, sync *_utils.SyncUtil) *LocalNotificationService {
	return &LocalNotificationService{repo: repo, sync: sync}
}

func (s *LocalNotificationService) Count(ctx context.Context, _ dto.CountRequest) (*notificationpb.CountResponse, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	count, err := s.repo.CountUnread(ctx, profileID)
	if err != nil {
		return nil, err
	}
	now := time.Now().UnixMilli()
	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyNotificationMe, 0)
		_ = s.sync.PutTimestamp(ctx, key, now)
	}
	return &notificationpb.CountResponse{
		Count:     count,
		Unread:    count,
		Timestamp: now,
		Unchanged: false,
	}, nil
}

func (s *LocalNotificationService) List(ctx context.Context, req dto.ListRequest) (*notificationpb.ListResponse, error) {
	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyNotificationMe, 0)
		updated := s.sync.HasUpdated(ctx, key, req.Timestamp)
		s.sync.PutTimeRequest(ctx, key, req.Timestamp)
		if !updated && req.Page == 0 {
			ts := s.sync.EnsureTimestamp(ctx, key)
			return &notificationpb.ListResponse{
				Data:      nil,
				Total:     0,
				Timestamp: ts,
				Unchanged: true,
			}, nil
		}
	}

	profileID := _utils.GetProfileIdWithContext(ctx)
	rows, total, err := s.repo.List(ctx, profileID, req)
	if err != nil {
		return nil, err
	}
	data := make([]*notificationpb.NotificationDTO, 0, len(rows))
	for _, row := range rows {
		data = append(data, mapNotificationPB(row))
	}
	return &notificationpb.ListResponse{
		Data:      data,
		Total:     total,
		Unchanged: false,
	}, nil
}

func (s *LocalNotificationService) Read(ctx context.Context, id uint64) (*notificationpb.StatusResponse, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	if _, err := s.repo.MarkRead(ctx, profileID, id); err != nil {
		return nil, err
	}
	s.invalidateSync(ctx)
	return &notificationpb.StatusResponse{
		Success: true,
		Message: "Thông báo đã được đánh dấu là đã đọc",
	}, nil
}

func (s *LocalNotificationService) ReadMany(ctx context.Context, ids []uint64) (*notificationpb.StatusResponse, error) {
	if len(ids) == 0 {
		return &notificationpb.StatusResponse{
			Success: true,
			Message: "Không có thông báo nào để đánh dấu",
		}, nil
	}
	profileID := _utils.GetProfileIdWithContext(ctx)
	count, err := s.repo.MarkReadMany(ctx, profileID, ids)
	if err != nil {
		return nil, err
	}
	s.invalidateSync(ctx)
	return &notificationpb.StatusResponse{
		Success: true,
		Message: fmt.Sprintf("Đã đánh dấu %d thông báo là đã đọc", count),
	}, nil
}

func (s *LocalNotificationService) ReadAll(ctx context.Context) (*notificationpb.StatusResponse, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	if _, err := s.repo.MarkReadAll(ctx, profileID); err != nil {
		return nil, err
	}
	s.invalidateSync(ctx)
	return &notificationpb.StatusResponse{
		Success: true,
		Message: "Thông báo đã được đánh dấu là đã đọc",
	}, nil
}

func (s *LocalNotificationService) Remove(ctx context.Context, id uint64) (*notificationpb.StatusResponse, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	if err := s.repo.Remove(ctx, profileID, id); err != nil {
		return nil, err
	}
	s.invalidateSync(ctx)
	return &notificationpb.StatusResponse{
		Success: true,
		Message: "Thông báo đã được xóa",
	}, nil
}

func (s *LocalNotificationService) invalidateSync(ctx context.Context) {
	if s.sync == nil {
		return
	}
	key := s.sync.GetKey(ctx, _utils.SyncKeyNotificationMe, 0)
	_ = s.sync.PutTimestamp(ctx, key, time.Now().UnixMilli())
}

func mapNotificationPB(n *entity.Notification) *notificationpb.NotificationDTO {
	if n == nil {
		return nil
	}
	targetID := n.TargetID
	return &notificationpb.NotificationDTO{
		Id:               n.ID,
		Avatar:           n.Avatar,
		Title:            n.Title,
		Message:          n.Message,
		Timestamp:        _utils.FormatTimeToString(n.CreatedAt),
		NotificationType: uint32(n.NotificationType),
		TargetId:         &targetID,
		IsRead:           n.IsRead,
	}
}
