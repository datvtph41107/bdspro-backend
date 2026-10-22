package services

import (
	"context"
	"time"

	_errors "common/errors"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"

	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/mapper"
	"mdn4/modules/bdspro/repository"
)

// ListingListOutcome is listing/me response + sync metadata.
type ListingListOutcome struct {
	List      *bdspropb.PostListResponse
	Timestamp int64
	Unchanged bool
}

// ListingService lists v4 listing rows (product-attached publish metadata).
type ListingService interface {
	ListMine(ctx context.Context, req dto.ListingListRequest) (*ListingListOutcome, error)
}

// @bind: mdn4/modules/bdspro/services.ListingService
type LocalListingService struct {
	repo repository.ListingRepository
	sync *_utils.SyncUtil
}

func NewLocalListingService(repo repository.ListingRepository, sync *_utils.SyncUtil) *LocalListingService {
	return &LocalListingService{repo: repo, sync: sync}
}

func (s *LocalListingService) ListMine(ctx context.Context, req dto.ListingListRequest) (*ListingListOutcome, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, _errors.UnauthorizedException()
	}

	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyListingMe, 0)
		if req.Timestamp > 0 && req.Page == 0 && !s.sync.HasUpdated(ctx, key, req.Timestamp) {
			ts := s.sync.GetTimestamp(ctx, key)
			if ts == 0 {
				ts = req.Timestamp
			}
			return &ListingListOutcome{
				Timestamp: ts,
				Unchanged: true,
			}, nil
		}
	}

	rows, total, err := s.repo.List(ctx, req, profileID)
	if err != nil {
		return nil, err
	}

	ts := s.listingMeTimestamp(ctx, req.Page, true)

	return &ListingListOutcome{
		List: &bdspropb.PostListResponse{
			Data:  mapper.ListingsToPosts(rows),
			Total: int32(total),
		},
		Timestamp: ts,
	}, nil
}

// listingMeTimestamp returns the Redis sync ts for listing/me (same key as GET /v2/bdspro/ts listing_me).
// Does not overwrite an existing ts on list reads; seeds only when missing (e.g. after Del on create).
func (s *LocalListingService) listingMeTimestamp(ctx context.Context, page uint32, seedIfMissing bool) int64 {
	if s.sync == nil {
		return time.Now().UnixMilli()
	}
	key := s.sync.GetKey(ctx, _utils.SyncKeyListingMe, 0)
	if ts := s.sync.GetTimestamp(ctx, key); ts > 0 {
		return ts
	}
	if !seedIfMissing || page != 0 {
		return 0
	}
	ts := time.Now().UnixMilli()
	_ = s.sync.PutTimestamp(ctx, key, ts)
	return ts
}
