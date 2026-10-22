package services

import (
	"context"
	"strconv"

	_utils "common/utils"
	"mdn4/internal/interface/providers"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/entity"
	"mdn4/modules/bdspro/repository"
	notificationpb "pb/types/notification"
)

type PropertyHistoryService interface {
	ListBySubject(ctx context.Context, subjectID uint64, req dto.SyncByIDRequest) (*notificationpb.GetPropertyHistoryResponse, error)
}

// @bind: mdn4/modules/bdspro/services.PropertyHistoryService
type LocalPropertyHistoryService struct {
	sync    *_utils.SyncUtil
	repo    repository.PropertyHistoryRepository
	profile providers.ProfileProvider
}

func NewLocalPropertyHistoryService(
	sync *_utils.SyncUtil,
	repo repository.PropertyHistoryRepository,
	profile providers.ProfileProvider,
) *LocalPropertyHistoryService {
	return &LocalPropertyHistoryService{sync: sync, repo: repo, profile: profile}
}

func (s *LocalPropertyHistoryService) ListBySubject(ctx context.Context, subjectID uint64, req dto.SyncByIDRequest) (*notificationpb.GetPropertyHistoryResponse, error) {
	if subjectID == 0 {
		return &notificationpb.GetPropertyHistoryResponse{}, nil
	}
	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyPropertyHistories, subjectID)
		if !s.sync.HasUpdated(ctx, key, req.Timestamp) && req.Page == 0 && req.Timestamp > 0 {
			return &notificationpb.GetPropertyHistoryResponse{}, nil
		}
	}

	filter := repository.PropertyHistoryFilter{
		SubjectID: subjectID,
		Keyword:   req.Keyword,
		Page:      req.Page,
		Size:      req.Size,
	}
	if req.Action > 0 {
		a := req.Action
		filter.Action = &a
	}
	if req.ActorID > 0 {
		id := req.ActorID
		filter.ActorID = &id
	}

	list, total, err := s.repo.ListBySubject(ctx, filter)
	if err != nil {
		return nil, err
	}

	actorIDs := make([]uint64, 0, len(list))
	seen := map[uint64]struct{}{}
	for _, h := range list {
		if h.ActorID == 0 {
			continue
		}
		if _, ok := seen[h.ActorID]; ok {
			continue
		}
		seen[h.ActorID] = struct{}{}
		actorIDs = append(actorIDs, h.ActorID)
	}
	actors := map[uint64]*notificationpb.ActorInfo{}
	if s.profile != nil && len(actorIDs) > 0 {
		if m, err := s.profile.GetMapProfileByIds(ctx, actorIDs); err == nil {
			for id, p := range m {
				if p == nil {
					continue
				}
				actors[id] = &notificationpb.ActorInfo{Id: id, FullName: p.FullName, Avatar: p.Avatar}
			}
		}
	}

	data := make([]*notificationpb.PropertyHistoryItem, 0, len(list))
	for _, h := range list {
		data = append(data, mapPropertyHistory(h, actors[h.ActorID]))
	}
	return &notificationpb.GetPropertyHistoryResponse{Data: data, Total: total}, nil
}

func mapPropertyHistory(h *entity.PropertyHistory, actor *notificationpb.ActorInfo) *notificationpb.PropertyHistoryItem {
	if h == nil {
		return nil
	}
	item := &notificationpb.PropertyHistoryItem{
		Id: h.ID, SubjectId: h.SubjectID, SubjectType: h.SubjectType,
		Action: h.Action.String(), Description: h.Description, Actor: actor,
	}
	if h.OccurredAt != nil {
		item.OccurredAt = strconv.FormatInt(h.OccurredAt.UnixMilli(), 10)
	}
	if h.CreatedAt != nil {
		item.CreatedAt = strconv.FormatInt(h.CreatedAt.UnixMilli(), 10)
	}
	if h.UpdatedAt != nil {
		item.UpdatedAt = strconv.FormatInt(h.UpdatedAt.UnixMilli(), 10)
	}
	return item
}
