package services

import (
	"context"
	"strings"
	"time"

	_utils "common/utils"
	tqdpb "pb/types/tqd"
	"userv4/modules/tqd/dto"
	"userv4/modules/tqd/entity"
	"userv4/modules/tqd/enums"
	"userv4/modules/tqd/repository"
)

const pageSizeMax = 500

// LocalLocationService serves TQD province/ward catalogs from the TQD DSN + SyncUtil.
type LocalLocationService struct {
	repo repository.LocationRepository
	sync *_utils.SyncUtil
}

func NewLocalLocationService(repo repository.LocationRepository, sync *_utils.SyncUtil) *LocalLocationService {
	return &LocalLocationService{repo: repo, sync: sync}
}

func (s *LocalLocationService) ListProvinces(ctx context.Context, req dto.ListProvincesRequest) (*tqdpb.ListProvincesResponse, error) {
	catalogTS, updated := s.checkProvincesWardsCatalog(ctx, req.Timestamp)
	if !updated {
		return &tqdpb.ListProvincesResponse{Timestamp: catalogTS}, nil
	}

	provinces, err := s.repo.ListAllProvinces(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*tqdpb.Province, len(provinces))
	for i, p := range provinces {
		out[i] = mapProvince(p)
	}
	return &tqdpb.ListProvincesResponse{
		Provinces: out,
		Timestamp: catalogTS,
	}, nil
}

func (s *LocalLocationService) ListWards(ctx context.Context, req dto.ListWardsRequest) (*tqdpb.ClientListWardsResponse, error) {
	catalogTS, updated := s.checkProvincesWardsCatalog(ctx, req.Timestamp)
	if !updated {
		return &tqdpb.ClientListWardsResponse{Timestamp: catalogTS}, nil
	}

	page := 1
	pageSize := pageSizeMax
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.PageSize != nil {
		pageSize = int(*req.PageSize)
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = pageSizeMax
	}
	if pageSize > pageSizeMax {
		pageSize = pageSizeMax
	}

	wards, total, err := s.repo.ListAllWards(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}
	data := make([]*tqdpb.Ward, len(wards))
	for i, w := range wards {
		data[i] = mapWardListItem(w)
	}
	return &tqdpb.ClientListWardsResponse{
		Data: data,
		Pagination: &tqdpb.LocationPagination{
			Page:  int32(page),
			Size:  int32(pageSize),
			Total: total,
		},
		Timestamp: catalogTS,
	}, nil
}

func (s *LocalLocationService) checkProvincesWardsCatalog(ctx context.Context, clientTimestamp int64) (int64, bool) {
	if s.sync == nil {
		now := time.Now().UnixMilli()
		return now, true
	}
	key := _utils.SyncKeyTQDProvincesWards
	catalogTS := s.sync.EnsureTimestamp(ctx, key)
	return catalogTS, s.sync.HasUpdated(ctx, key, clientTimestamp)
}

func mapProvince(p *entity.Province) *tqdpb.Province {
	if p == nil {
		return nil
	}
	createdAt := p.CreatedAt
	updatedAt := p.UpdatedAt
	resp := &tqdpb.Province{
		Id:                    p.ID,
		Code:                  p.Code,
		FullName:              p.FullName,
		ShortName:             p.ShortName,
		Latitude:              p.Lat,
		Longitude:             p.Lng,
		DivisionType:          p.DivisionType,
		WardCount:             p.WardCount,
		OfficialName:          p.OfficialName,
		OfficialOperationDate: _utils.FormatTimeToString(p.OfficialOperationDate),
		NaturalAreaKm2:        p.NaturalAreaKm2,
		Population:            p.Population,
		PopulationDensity:     p.PopulationDensity,
		Description:           p.Description,
		UnitType:              uint32(p.UnitType),
		UnitTypeName:          p.UnitType.String(),
		DataSource:            uint32(p.DataSource),
		DataSourceName:        p.DataSource.String(),
		OfficialRole:          p.OfficialRole,
		OfficialAddress:       p.OfficialAddress,
		CompositeId:           encodeLocationCompositeID(p.ID, _utils.LocationIDKindProvince),
		CreatedAt:             _utils.FormatTimeToString(&createdAt),
		UpdatedAt:             _utils.FormatTimeToString(&updatedAt),
	}
	if p.SeoID != nil && *p.SeoID > 0 {
		seoID := *p.SeoID
		resp.SeoId = &seoID
	}
	if p.PlanDocID != nil && *p.PlanDocID > 0 {
		planDocID := *p.PlanDocID
		resp.PlanDocId = &planDocID
	}
	return resp
}

func mapWardListItem(w *entity.Ward) *tqdpb.Ward {
	if w == nil {
		return nil
	}
	unitType := wardUnitType(w)
	if w.UnitType != 0 {
		unitType = w.UnitType
	}
	resp := &tqdpb.Ward{
		Id:           w.ID,
		ShortName:    w.ShortName,
		Latitude:     w.Lat,
		Longitude:    w.Lng,
		ProvinceId:   w.ProvinceID,
		CompositeId:  encodeLocationCompositeID(w.ID, _utils.LocationIDKindWard),
		UnitType:     uint32(unitType),
		UnitTypeName: unitType.String(),
	}
	if w.SeoID != nil && *w.SeoID > 0 {
		seoID := *w.SeoID
		resp.SeoId = &seoID
	}
	return resp
}

func encodeLocationCompositeID(sourceID uint64, kind uint32) uint64 {
	id, _ := _utils.EncodeLocationID(sourceID, kind)
	return id
}

func wardUnitType(w *entity.Ward) enums.LocationUnitType {
	if w == nil {
		return enums.LocationUnitTypeCity
	}
	name := strings.ToLower(strings.TrimSpace(w.DivisionType + " " + w.FullName + " " + w.OfficialName))
	if strings.Contains(name, "đặc khu") || strings.Contains(name, "dac khu") {
		return enums.LocationUnitTypeSpecialZone
	}
	return enums.LocationUnitTypeCity
}
