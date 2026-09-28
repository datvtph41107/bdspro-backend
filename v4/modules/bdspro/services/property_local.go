package services

import (
	"context"
	"strconv"
	"time"

	_errors "common/errors"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"

	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/entity"
	"userv4/modules/bdspro/repository"
)

// LocalPropertyService owns property list/detail against the v4 DSN + SyncUtil.
type LocalPropertyService struct {
	sync *_utils.SyncUtil
	repo repository.PropertyRepository
}

func NewLocalPropertyService(sync *_utils.SyncUtil, repo repository.PropertyRepository) *LocalPropertyService {
	return &LocalPropertyService{sync: sync, repo: repo}
}

func (s *LocalPropertyService) ListMe(ctx context.Context, req dto.PropertyListRequest) (*bdspropb.PropertyListResponse, error) {
	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyPropertyMe, 0)
		if !s.sync.HasUpdated(ctx, key, req.Timestamp) && req.Page == 0 && req.Timestamp > 0 {
			return &bdspropb.PropertyListResponse{
				Data: []*bdspropb.PropertyDetailResponse{},
				Pagination: &bdspropb.Pagination{
					Page: int32(req.Page),
					Size: int32(req.Size),
				},
			}, nil
		}
	}

	profileID := _utils.GetProfileIdWithContext(ctx)
	size := req.Size
	if size == 0 {
		size = 20
	}
	list, total, err := s.repo.ListMine(ctx, profileID, req.Page, size)
	if err != nil {
		return nil, err
	}

	data := make([]*bdspropb.PropertyDetailResponse, 0, len(list))
	for _, p := range list {
		data = append(data, mapPropertyToPB(p))
	}

	pages := int32(0)
	if size > 0 {
		pages = int32((total + int64(size) - 1) / int64(size))
	}
	hasNext := int64(req.Page+1)*int64(size) < total
	hasPrev := req.Page > 0

	if s.sync != nil && req.Page == 0 {
		_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyPropertyMe, 0), time.Now().UnixMilli())
	}

	return &bdspropb.PropertyListResponse{
		Data: data,
		Pagination: &bdspropb.Pagination{
			Page:        int32(req.Page),
			Size:        int32(size),
			Total:       total,
			Pages:       pages,
			HasNext:     hasNext,
			HasPrevious: hasPrev,
		},
	}, nil
}

func (s *LocalPropertyService) GetDetail(ctx context.Context, req dto.PropertyListRequest) (*bdspropb.PropertyDetailResponse, error) {
	if req.ID == 0 {
		return nil, _errors.BadRequestException("id is required")
	}
	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyPropertyDetail, req.ID)
		if !s.sync.HasUpdated(ctx, key, req.Timestamp) && req.Timestamp > 0 {
			return &bdspropb.PropertyDetailResponse{Id: req.ID}, nil
		}
	}
	p, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, _errors.NotFoundException("property not found")
	}
	return mapPropertyToPB(p), nil
}

func mapPropertyToPB(p *entity.PropertyLineage) *bdspropb.PropertyDetailResponse {
	if p == nil {
		return nil
	}
	out := &bdspropb.PropertyDetailResponse{Id: p.ID}
	if p.CreatedAt != nil {
		out.CreatedAt = strconv.FormatInt(p.CreatedAt.UnixMilli(), 10)
	}
	if p.UpdatedAt != nil {
		out.UpdatedAt = strconv.FormatInt(p.UpdatedAt.UnixMilli(), 10)
	}
	if p.Info != nil {
		title := p.Info.Title
		ident := p.Info.Identifier
		info := &bdspropb.PropertyInfo{
			Title:      &title,
			Identifier: &ident,
		}
		if p.Info.PropertyTypeID != nil {
			info.PropertyTypeId = p.Info.PropertyTypeID
		}
		if p.Info.ProjectID != nil {
			info.ProjectId = p.Info.ProjectID
		}
		if p.Info.AvatarID != nil {
			info.AvatarId = p.Info.AvatarID
		}
		out.Info = info
	}
	if p.Location != nil {
		loc := &bdspropb.PropertyLocationInfo{
			Id:     p.Location.ID,
			Detail: p.Location.AddressDetail,
			MapUrl: p.Location.MapURL,
		}
		if p.Location.ProvinceID != nil {
			loc.ProvinceId = p.Location.ProvinceID
		}
		if p.Location.DistrictID != nil {
			loc.DistrictId = p.Location.DistrictID
		}
		if p.Location.WardID != nil {
			loc.WardId = p.Location.WardID
		}
		if p.Location.Latitude != nil {
			loc.Latitude = p.Location.Latitude
		}
		if p.Location.Longitude != nil {
			loc.Longitude = p.Location.Longitude
		}
		out.Location = loc
	}
	return out
}
