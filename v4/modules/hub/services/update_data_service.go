package services

import (
	"context"
	"time"

	"mdn4/modules/hub/dto"
	"mdn4/modules/hub/entity"
	"mdn4/modules/hub/repository"
)

// UpdateDataService tracks per-owner resource change versions for client sync.
type UpdateDataService interface {
	CheckVersionSync(ctx context.Context, ownerID uint64, resource string, lastSync int64, limit int64) (*dto.CheckVersionSyncResponse, error)
	RefreshSyncIDs(ctx context.Context, ownerID uint64, resource string, lastSync int64, limit int64) (*dto.RefreshSyncIdsResponse, error)
	FlushSyncIDs(ctx context.Context, ownerID uint64, resource string, limit int64) error
}

// @bind: mdn4/modules/hub/services.UpdateDataService
type LocalUpdateDataService struct {
	repo repository.UpdateDataRepository
}

func NewLocalUpdateDataService(repo repository.UpdateDataRepository) *LocalUpdateDataService {
	return &LocalUpdateDataService{repo: repo}
}

func (s *LocalUpdateDataService) CheckVersionSync(ctx context.Context, ownerID uint64, resource string, lastSync int64, limit int64) (*dto.CheckVersionSyncResponse, error) {
	rt := entity.ParseSyncResource(resource)
	lim := int(limit)
	if lim <= 0 {
		lim = 100
	}
	ids, err := s.repo.GetChangedIDsSince(ctx, ownerID, rt, lastSync, lim)
	if err != nil {
		return nil, err
	}
	return &dto.CheckVersionSyncResponse{
		NeedSync:       len(ids) > 0,
		CurrentVersion: time.Now().Unix(),
		ChangeCount:    int32(len(ids)),
		ChangedIDs:     ids,
	}, nil
}

func (s *LocalUpdateDataService) RefreshSyncIDs(ctx context.Context, ownerID uint64, resource string, lastSync int64, limit int64) (*dto.RefreshSyncIdsResponse, error) {
	rt := entity.ParseSyncResource(resource)
	lim := int(limit)
	if lim <= 0 {
		lim = 100
	}
	ids, err := s.repo.GetChangedIDsSince(ctx, ownerID, rt, lastSync, lim)
	if err != nil {
		return nil, err
	}
	return &dto.RefreshSyncIdsResponse{
		NeedSync:       len(ids) > 0,
		CurrentVersion: time.Now().Unix(),
		ChangeCount:    int32(len(ids)),
		ChangedIDs:     ids,
	}, nil
}

func (s *LocalUpdateDataService) FlushSyncIDs(ctx context.Context, ownerID uint64, resource string, limit int64) error {
	rt := entity.ParseSyncResource(resource)
	return s.repo.FlushOld(ctx, ownerID, rt, limit)
}
