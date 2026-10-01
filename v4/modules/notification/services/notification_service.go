package services

import (
	"context"
	"fmt"
	"time"

	_utils "common/utils"
	shared_enum "pb/enums"
	notificationpb "pb/types/notification"
	"userv4/internal/interface/providers"
	"userv4/modules/notification/dto"
	"userv4/modules/notification/entity"
	"userv4/modules/notification/repository"
)

// NotificationService owns the user inbox flows migrated into modules/notification.
type NotificationService interface {
	Count(ctx context.Context, req dto.CountRequest) (*notificationpb.CountResponse, error)
	List(ctx context.Context, req dto.ListRequest) (*notificationpb.ListResponse, error)
	Read(ctx context.Context, id uint64) (*notificationpb.StatusResponse, error)
	ReadMany(ctx context.Context, ids []uint64) (*notificationpb.StatusResponse, error)
	ReadAll(ctx context.Context) (*notificationpb.StatusResponse, error)
	Remove(ctx context.Context, id uint64) (*notificationpb.StatusResponse, error)
	SearchAdminHistory(ctx context.Context, req dto.AdminHistorySearchRequest) ([]*notificationpb.AdminHistoryDTO, int64, error)
}

// LocalNotificationService owns inbox flows against the v4 Postgres DSN.
type LocalNotificationService struct {
	repo     repository.NotificationRepository
	sync     *_utils.SyncUtil
	profiles providers.ProfileProvider
}

func NewLocalNotificationService(
	repo repository.NotificationRepository,
	sync *_utils.SyncUtil,
	profiles providers.ProfileProvider,
) *LocalNotificationService {
	return &LocalNotificationService{repo: repo, sync: sync, profiles: profiles}
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

func (s *LocalNotificationService) SearchAdminHistory(ctx context.Context, req dto.AdminHistorySearchRequest) ([]*notificationpb.AdminHistoryDTO, int64, error) {
	filter := dto.AdminHistoryFilter{
		Page:       int(req.Page),
		Size:       int(req.Size),
		AdminID:    req.AdminID,
		OwnerID:    req.OwnerID,
		TargetID:   req.TargetID,
		AdminRole:  req.AdminRole,
		FromDate:   req.FromDate,
		ToDate:     req.ToDate,
		IsInternal: req.IsInternal,
	}
	if req.OwnerType != nil {
		v := shared_enum.EOwnerType(*req.OwnerType)
		filter.OwnerType = &v
	}
	if req.TargetType != nil {
		v := shared_enum.ETargetHistory(*req.TargetType)
		filter.TargetType = &v
	}
	if req.ActionType != nil {
		v := shared_enum.EHistory(*req.ActionType)
		filter.ActionType = &v
	}

	rows, total, err := s.repo.SearchAdminHistory(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	s.populateAdminInfo(ctx, rows)

	out := make([]*notificationpb.AdminHistoryDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapAdminHistoryPB(row))
	}
	return out, total, nil
}

func (s *LocalNotificationService) populateAdminInfo(ctx context.Context, rows []*entity.AdminHistory) {
	if s.profiles == nil || len(rows) == 0 {
		return
	}
	ids := make([]uint64, 0, len(rows))
	seen := make(map[uint64]struct{}, len(rows))
	for _, row := range rows {
		if row == nil || row.AdminID == 0 {
			continue
		}
		if _, ok := seen[row.AdminID]; ok {
			continue
		}
		seen[row.AdminID] = struct{}{}
		ids = append(ids, row.AdminID)
	}
	if len(ids) == 0 {
		return
	}
	profiles, err := s.profiles.GetMapProfileByIds(ctx, ids)
	if err != nil || len(profiles) == 0 {
		return
	}
	for _, row := range rows {
		if row == nil {
			continue
		}
		empty := ""
		row.AdminName = &empty
		row.AdminAvatar = &empty
		status := "success"
		row.Status = &status
		if p := profiles[row.AdminID]; p != nil {
			if p.FullName != "" {
				name := p.FullName
				row.AdminName = &name
			}
			if p.Avatar != "" {
				avatar := p.Avatar
				row.AdminAvatar = &avatar
			}
		}
	}
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

func mapAdminHistoryPB(h *entity.AdminHistory) *notificationpb.AdminHistoryDTO {
	if h == nil {
		return nil
	}
	out := &notificationpb.AdminHistoryDTO{
		Id:         h.ID,
		TargetId:   h.TargetID,
		TargetType: int32(h.TargetType),
		ActionType: int32(h.ActionType),
		Title:      h.Title,
		Note:       []string(h.Note),
		PreStage:   h.PreStage,
		AfterStage: h.AfterStage,
		AdminId:    h.AdminID,
		OwnerId:    h.OwnerID,
		OwnerType:  int32(h.OwnerOf),
		IsInternal: h.IsInternal,
		AdminRole:  h.AdminRole,
		IpAddress:  h.IPAddress,
		UserAgent:  h.UserAgent,
		CreatedAt:  _utils.FormatTimeToString(h.CreatedAt),
		UpdatedAt:  _utils.FormatTimeToString(h.UpdatedAt),
		AdminName:  h.AdminName,
		AdminAvatar: h.AdminAvatar,
	}
	return out
}
