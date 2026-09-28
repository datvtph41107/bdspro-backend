package services

import (
	"context"
	"fmt"
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

func (s *LocalPropertyService) GetTimestamps(ctx context.Context, id uint64) (*bdspropb.PropertyTimestampsResponse, error) {
	if id == 0 {
		return nil, _errors.BadRequestException("property id is required")
	}
	if s.sync == nil {
		return &bdspropb.PropertyTimestampsResponse{Timestamps: map[string]int64{}}, nil
	}

	redisKeys := []string{
		s.sync.GetKey(ctx, _utils.SyncKeyPropertyDetail, id),
		s.sync.GetKey(ctx, _utils.SyncKeyPropertyHistories, id),
	}
	names := []string{"property_detail", "property_histories"}

	vals, err := s.sync.MGet(ctx, redisKeys)
	if err != nil {
		return nil, _errors.InternalServerException(err.Error())
	}

	timestamps := make(map[string]int64, len(names))
	for i, v := range vals {
		timestamps[names[i]] = parseSyncTimestamp(v)
	}
	return &bdspropb.PropertyTimestampsResponse{Timestamps: timestamps}, nil
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

// mapPropertyToPB mirrors bdspro PropertyDetailToProto for response parity.
func mapPropertyToPB(p *entity.PropertyLineage) *bdspropb.PropertyDetailResponse {
	if p == nil {
		return nil
	}
	out := &bdspropb.PropertyDetailResponse{Id: p.ID}
	if p.CreatedAt != nil {
		out.CreatedAt = _utils.FormatTimeToString(p.CreatedAt)
	}
	if p.UpdatedAt != nil {
		out.UpdatedAt = _utils.FormatTimeToString(p.UpdatedAt)
	}
	if p.NationalID != "" {
		out.Lineage = &bdspropb.PropertyLineageInfo{NationalId: p.NationalID}
	}

	if p.Identify != nil && p.Identify.ID != 0 {
		out.Identifier = &bdspropb.PropertyIdentify{
			Id:        p.Identify.ID,
			Pid:       p.Identify.PID,
			Version:   p.Identify.Version,
			Code:      fmt.Sprintf("BDS%d", p.Identify.PID),
			LineageId: p.ID,
		}
		if p.Identify.LineageID != nil {
			out.Identifier.LineageId = *p.Identify.LineageID
		}
	}

	if p.Statistic != nil {
		out.Statistic = &bdspropb.PropertyStatistic{
			Id:           p.Statistic.ID,
			ProductCount: &p.Statistic.ProductCount,
			AssetCount:   &p.Statistic.AssetCount,
			ListingCount: &p.Statistic.ListingCount,
		}
	}

	if p.Info != nil {
		pi := p.Info
		info := &bdspropb.PropertyInfo{
			PropertyTypeId: pi.PropertyTypeID,
			ProjectId:      pi.ProjectID,
			AvatarId:       pi.AvatarID,
		}
		if pi.Title != "" {
			t := pi.Title
			info.Title = &t
		}
		if pi.Identifier != "" {
			id := pi.Identifier
			info.Identifier = &id
		}
		if pi.UnitCode != "" {
			u := pi.UnitCode
			info.UnitCode = &u
		}
		if pi.Level != "" {
			l := pi.Level
			info.Level = &l
		}
		if pi.Note != "" {
			n := pi.Note
			info.Note = &n
		}
		if pi.LegalStatus != 0 {
			v := pi.LegalStatus
			info.LegalStatus = &v
		}
		out.Info = info

		if pi.Avatar != nil {
			out.Avatar = &bdspropb.MediaItem{
				Id:        pi.Avatar.ID,
				MediaType: pi.Avatar.MediaType,
				MediaUrl:  pi.Avatar.MediaURL,
			}
		}
		if pi.PropertyType != nil {
			out.PropertyType = &bdspropb.PropertyType{
				Id:   pi.PropertyType.ID,
				Name: pi.PropertyType.Name,
			}
		}
	}

	if p.Location != nil {
		loc := p.Location
		location := &bdspropb.PropertyLocationInfo{
			Id:         loc.ID,
			Detail:     loc.AddressDetail,
			MapUrl:     loc.MapURL,
			RegionId:   loc.RegionID,
			ProvinceId: loc.ProvinceID,
			DistrictId: loc.DistrictID,
			WardId:     loc.WardID,
			Latitude:   loc.Latitude,
			Longitude:  loc.Longitude,
		}
		if loc.Province != nil {
			location.ProvinceName = loc.Province.Name
		}
		if loc.Ward != nil {
			location.WardName = loc.Ward.Name
		}
		out.Location = location
	}

	if p.LandInfo != nil {
		li := p.LandInfo
		land := &bdspropb.PropertyLandInfoDetail{
			Id:           li.ID,
			DocumentNo:   li.DocumentNo,
			IssuringAuth: li.IssuringAuth,
			Note:         li.Note,
			LandNote:     li.LandNote,
			DocumentType: li.DocumentType,
			AreaTotal:    li.AreaTotal,
			AreaLand:     li.AreaLand,
			AreaPlant:    li.AreaPlant,
			FrontWidth:   li.FrontWidth,
			Depth:        li.Depth,
			StreetWidth:  li.StreetWidth,
		}
		if li.Plot > 0 {
			v := li.Plot
			land.Plot = &v
		}
		if li.Sheet > 0 {
			v := li.Sheet
			land.Sheet = &v
		}
		if li.PurposeUsed != 0 {
			v := li.PurposeUsed
			land.PurposeUsed = &v
		}
		out.LandInfo = land
	}

	if p.BuildingInfo != nil {
		out.BuildingInfo = mapBuildingInfoToPB(p.BuildingInfo)
	}

	if len(p.MediaList) > 0 {
		out.MediaList = make([]*bdspropb.MediaItem, 0, len(p.MediaList))
		for _, m := range p.MediaList {
			out.MediaList = append(out.MediaList, &bdspropb.MediaItem{
				Id:        m.ID,
				MediaUrl:  m.MediaURL,
				MediaType: m.MediaType,
			})
		}
		if out.Avatar == nil {
			first := p.MediaList[0]
			out.Avatar = &bdspropb.MediaItem{
				Id:        first.ID,
				MediaUrl:  first.MediaURL,
				MediaType: first.MediaType,
			}
		}
	}
	return out
}
