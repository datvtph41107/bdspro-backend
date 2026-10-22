package services

import (
	"bytes"
	"common/validation"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	_dto "common/domain/dto"
	domain "mdn4/modules/tqd/entity"
	qh "mdn4/modules/tqd/entity/qh"
	"mdn4/modules/tqd/enums"
	"mdn4/modules/tqd/geometry"
	repo "mdn4/modules/tqd/repository"
)

// ErrRegionNotFound và ErrRegionSyncReferenceNotFound do usecase region sở hữu.
var ErrRegionNotFound = errors.New("region not found")
var ErrRegionSyncReferenceNotFound = errors.New("legend with landUse not found")

// RegionUpdatePatch — các field có thể sửa qua UpdateRegion (không gồm geometry).
// Chỉ gán pointer cho field thực sự muốn đổi; LabelID = pointer tới 0 nghĩa là gỡ label.
type RegionUpdatePatch struct {
	RegionID    uint64
	Name        *string
	DisplayName *string
	Description *string
	Status      *enums.RegionStatus
	LabelID     *uint64
}

type RegionSyncResult struct {
	LayerID     uint64
	LabelID     uint64
	LandUseID   uint64
	LegendID    uint64
	SyncedCount int64
}

type RegionUsecase interface {
	Create(ctx context.Context, layerID uint64, name string, geometryBytes []byte, labelID uint64) (*qh.QHRegion, error)
	Update(ctx context.Context, patch RegionUpdatePatch) (*qh.QHRegion, error)
	GetByID(ctx context.Context, id uint64) (*qh.QHRegion, error)
	List(ctx context.Context, filter *repo.RegionFilter) ([]qh.QHRegion, int64, error)
	// ListForAdminByLayer kiểm tra layer tồn tại rồi list regions (admin)
	ListForAdminByLayer(ctx context.Context, layerID uint64, filter *repo.RegionFilter) ([]qh.QHRegion, int64, error)
	Delete(ctx context.Context, id uint64) error
	SyncRegion(ctx context.Context, layerID, labelID uint64) (*RegionSyncResult, error)
	UpdateProcessingStatus(ctx context.Context, regionID uint64, status enums.ProcessingStatus, userID uint64, notes string) error

	// Label Operations (1-1)
	AssignLabel(ctx context.Context, regionID uint64, labelID uint64, userID uint64, sourceFile string) error
	GetLabelByRegion(ctx context.Context, regionID uint64) (*qh.QHLabel, error)

	// Client Operations
	ListClient(ctx context.Context, layerID uint64, bBox *BBox, page, size uint32) ([]qh.QHRegion, int64, error)
	FindByPoint(ctx context.Context, lat, lng float64, layerID *uint64) (*qh.QHRegion, error)
	FindByBBox(ctx context.Context, minLng, minLat, maxLng, maxLat float64, layerID *uint64, labelID *uint64, page, size int) ([]qh.QHRegion, int64, error)
}

type BBox struct {
	MinLng, MinLat, MaxLng, MaxLat float64
}

type regionUsecaseImpl struct {
	regionRepo repo.RegionRepository
	layerRepo  repo.LayerRepository
	labelRepo  repo.QHLabelRepository
}

func NewRegionUsecase(
	regionRepo repo.RegionRepository,
	layerRepo repo.LayerRepository,
	labelRepo repo.QHLabelRepository,
) RegionUsecase {
	return &regionUsecaseImpl{
		regionRepo: regionRepo,
		layerRepo:  layerRepo,
		labelRepo:  labelRepo,
	}
}

func extractGeometryObjectFromBytes(raw []byte) (map[string]interface{}, error) {
	trim := bytes.TrimSpace(raw)
	if len(trim) == 0 {
		return nil, validation.New("geometry", "REGION_GEOMETRY_INVALID", "geometry is empty")
	}
	var root map[string]interface{}
	if err := json.Unmarshal(trim, &root); err != nil {
		return nil, validation.New("geometry", "REGION_GEOMETRY_INVALID", "geometry is not valid JSON")
	}
	t, _ := root["type"].(string)
	t = strings.TrimSpace(t)
	switch t {
	case "Feature":
		g, ok := root["geometry"].(map[string]interface{})
		if !ok || g == nil {
			return nil, validation.New("geometry", "REGION_GEOMETRY_INVALID", "GeoJSON Feature missing geometry object")
		}
		return g, nil
	case "Polygon", "MultiPolygon":
		return root, nil
	default:
		return nil, validation.New("geometry", "REGION_GEOMETRY_INVALID", fmt.Sprintf("unsupported GeoJSON type %q (expected Feature, Polygon, or MultiPolygon)", t))
	}
}

func (u *regionUsecaseImpl) Create(ctx context.Context, layerID uint64, name string, geometryBytes []byte, labelID uint64) (*qh.QHRegion, error) {
	if layerID == 0 {
		return nil, validation.New("layer_id", "REGION_LAYER_ID_INVALID", "layerId is required")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, validation.New("name", "REGION_NAME_INVALID", "name is required")
	}
	layer, err := u.layerRepo.GetByID(ctx, layerID)
	if err != nil {
		return nil, fmt.Errorf("get layer failed: %w", err)
	}
	if layer == nil {
		return nil, fmt.Errorf("%w: id=%d", ErrLayerNotFound, layerID)
	}

	geomObj, err := extractGeometryObjectFromBytes(geometryBytes)
	if err != nil {
		return nil, err
	}
	geomBytes, err := geometry.ToMultiPolygonGeoJSONBytes(geomObj)
	if err != nil {
		return nil, err
	}
	geoInfo, err := geometry.ParseGeometry(string(geomBytes))
	if err != nil {
		return nil, validation.New("geometry", "REGION_GEOMETRY_INVALID", "invalid geometry")
	}

	areaSqm := 0.0
	if geoInfo != nil {
		areaSqm = geoInfo.Area
	}

	propsBytes, _ := json.Marshal(map[string]interface{}{"source": "admin_create"})
	region := &qh.QHRegion{
		LayerID:            layerID,
		Name:               name,
		DisplayName:        name,
		LabelID:            &labelID,
		Geometry:           domain.RawGeometry{Raw: geomBytes},
		OriginalProperties: propsBytes,
		AreaSqm:            areaSqm,
		AreaHa:             areaSqm / 10000,
		Status:             enums.RegionStatusActive,
		ProcessingStatus:   enums.ProcessingStatusPending,
		Version:            1,
		IsLatest:           true,
	}
	if err := u.regionRepo.Create(ctx, region); err != nil {
		return nil, fmt.Errorf("create region failed: %w", err)
	}
	if region.ID == 0 {
		return nil, errors.New("create region failed: missing id")
	}
	created, err := u.regionRepo.GetByID(ctx, region.ID)
	if err != nil {
		return nil, fmt.Errorf("reload region failed: %w", err)
	}
	return created, nil
}

func (u *regionUsecaseImpl) Update(ctx context.Context, patch RegionUpdatePatch) (*qh.QHRegion, error) {
	if patch.RegionID == 0 {
		return nil, validation.New("id", "REGION_ID_INVALID", "id is required")
	}
	has := patch.Name != nil || patch.DisplayName != nil || patch.Description != nil ||
		patch.Status != nil || patch.LabelID != nil
	if !has {
		return nil, validation.New("fields", "REGION_FIELDS_INVALID", "at least one field to update is required")
	}

	region, err := u.regionRepo.GetByID(ctx, patch.RegionID)
	if err != nil {
		return nil, fmt.Errorf("get region failed: %w", err)
	}
	if region == nil {
		return nil, fmt.Errorf("%w: id=%d", ErrRegionNotFound, patch.RegionID)
	}

	if patch.Name != nil {
		v := strings.TrimSpace(*patch.Name)
		if v == "" {
			return nil, validation.New("name", "REGION_NAME_INVALID", "name cannot be empty")
		}
		region.Name = v
	}
	if patch.DisplayName != nil {
		region.DisplayName = strings.TrimSpace(*patch.DisplayName)
	}
	if patch.Description != nil {
		region.Description = *patch.Description
	}
	if patch.Status != nil {
		if !patch.Status.IsValid() {
			return nil, validation.New("status", "REGION_STATUS_INVALID", fmt.Sprintf("invalid status: %d", *patch.Status))
		}
		region.Status = *patch.Status
	}
	if patch.LabelID != nil {
		lid := *patch.LabelID
		if lid == 0 {
			region.LabelID = nil
		} else {
			label, lerr := u.labelRepo.GetByID(ctx, lid)
			if lerr != nil {
				return nil, fmt.Errorf("get label failed: %w", lerr)
			}
			if label == nil {
				return nil, fmt.Errorf("%w: id=%d", ErrLabelNotFound, lid)
			}
			// if label.LayerID != region.LayerID {
			// 	return nil, errors.New("label does not belong to region layer")
			// }
			region.LabelID = &lid
		}
	}

	if err := u.regionRepo.Update(ctx, region); err != nil {
		return nil, fmt.Errorf("update region failed: %w", err)
	}
	if region.LabelID != nil {
		u.labelRepo.EnsureLayerLink(ctx, *region.LabelID, region.LayerID)
	}
	return u.regionRepo.GetByID(ctx, patch.RegionID)
}

func (u *regionUsecaseImpl) GetByID(ctx context.Context, id uint64) (*qh.QHRegion, error) {
	region, err := u.regionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get region failed: %w", err)
	}
	if region == nil {
		return nil, fmt.Errorf("%w: id=%d", ErrRegionNotFound, id)
	}
	return region, nil
}

func (u *regionUsecaseImpl) List(ctx context.Context, filter *repo.RegionFilter) ([]qh.QHRegion, int64, error) {
	return u.regionRepo.List(ctx, filter)
}

func (u *regionUsecaseImpl) ListForAdminByLayer(ctx context.Context, layerID uint64, filter *repo.RegionFilter) ([]qh.QHRegion, int64, error) {
	if layerID == 0 {
		return nil, 0, validation.New("layer_id", "REGION_LAYER_ID_INVALID", "layerId is required")
	}
	layer, err := u.layerRepo.GetByID(ctx, layerID)
	if err != nil {
		return nil, 0, fmt.Errorf("get layer failed: %w", err)
	}
	if layer == nil {
		return nil, 0, fmt.Errorf("%w: id=%d", ErrLayerNotFound, layerID)
	}
	lid := layerID
	filter.LayerID = &lid
	return u.regionRepo.List(ctx, filter)
}

func (u *regionUsecaseImpl) Delete(ctx context.Context, id uint64) error {
	region, err := u.regionRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get region failed: %w", err)
	}
	if region == nil {
		return fmt.Errorf("%w: id=%d", ErrRegionNotFound, id)
	}
	return u.regionRepo.Delete(ctx, id)
}

func (u *regionUsecaseImpl) SyncRegion(ctx context.Context, layerID, labelID uint64) (*RegionSyncResult, error) {
	if layerID == 0 {
		return nil, validation.New("layer_id", "REGION_LAYER_ID_INVALID", "layerId is required")
	}
	if labelID == 0 {
		return nil, errors.New("labelId is required")
	}

	layer, err := u.layerRepo.GetByID(ctx, layerID)
	if err != nil {
		return nil, fmt.Errorf("get layer failed: %w", err)
	}
	if layer == nil {
		return nil, fmt.Errorf("%w: id=%d", ErrLayerNotFound, layerID)
	}

	label, err := u.labelRepo.GetByID(ctx, labelID)
	if err != nil {
		return nil, fmt.Errorf("get label failed: %w", err)
	}
	if label == nil {
		return nil, fmt.Errorf("%w: id=%d", ErrLabelNotFound, labelID)
	}

	ref, err := u.regionRepo.GetSyncReferenceByLayerLabel(ctx, layerID, labelID)
	if err != nil {
		return nil, fmt.Errorf("get region sync reference failed: %w", err)
	}
	if ref == nil {
		return nil, fmt.Errorf("%w: layer_id=%d label_id=%d", ErrRegionSyncReferenceNotFound, layerID, labelID)
	}

	total, err := u.regionRepo.SyncLandUseAndLegendByLayerLabel(ctx, layerID, labelID, ref.LandUseID, ref.LegendID)
	if err != nil {
		return nil, fmt.Errorf("sync region failed: %w", err)
	}

	return &RegionSyncResult{
		LayerID:     layerID,
		LabelID:     labelID,
		LandUseID:   ref.LandUseID,
		LegendID:    ref.LegendID,
		SyncedCount: total,
	}, nil
}

func (u *regionUsecaseImpl) UpdateProcessingStatus(ctx context.Context, regionID uint64, status enums.ProcessingStatus, userID uint64, notes string) error {
	region, err := u.regionRepo.GetByID(ctx, regionID)
	if err != nil {
		return fmt.Errorf("get region failed: %w", err)
	}
	if region == nil {
		return fmt.Errorf("%w: id=%d", ErrRegionNotFound, regionID)
	}
	region.ProcessingStatus = status
	return u.regionRepo.Update(ctx, region)
}

func (u *regionUsecaseImpl) AssignLabel(ctx context.Context, regionID uint64, labelID uint64, userID uint64, sourceFile string) error {
	// Kiểm tra region tồn tại
	region, err := u.regionRepo.GetByID(ctx, regionID)
	if err != nil {
		return fmt.Errorf("get region failed: %w", err)
	}
	if region == nil {
		return fmt.Errorf("%w: id=%d", ErrRegionNotFound, regionID)
	}

	// Kiểm tra label tồn tại
	label, err := u.labelRepo.GetByID(ctx, labelID)
	if err != nil {
		return fmt.Errorf("get label failed: %w", err)
	}
	if label == nil {
		return fmt.Errorf("%w: id=%d", ErrLabelNotFound, labelID)
	}

	// Gán label trực tiếp vào region
	region.LabelID = &labelID
	return u.regionRepo.Update(ctx, region)
}

func (u *regionUsecaseImpl) GetLabelByRegion(ctx context.Context, regionID uint64) (*qh.QHLabel, error) {
	region, err := u.regionRepo.GetByID(ctx, regionID)
	if err != nil {
		return nil, fmt.Errorf("get region failed: %w", err)
	}
	if region == nil {
		return nil, fmt.Errorf("%w: id=%d", ErrRegionNotFound, regionID)
	}
	if region.LabelID == nil {
		return nil, nil // Chưa có label
	}
	return u.labelRepo.GetByID(ctx, *region.LabelID)
}

// Client
func (u *regionUsecaseImpl) ListClient(ctx context.Context, layerID uint64, bBox *BBox, page, size uint32) ([]qh.QHRegion, int64, error) {
	offset := int((page - 1) * size)
	limit := int(size)

	if bBox != nil {
		return u.regionRepo.FindByBBox(ctx, bBox.MinLng, bBox.MinLat, bBox.MaxLng, bBox.MaxLat, &layerID, nil, limit, offset)
	}

	filter := &repo.ClientRegionFilter{
		LayerID: &layerID,
		Pagable: &_dto.Pagable{Page: page, Size: size},
	}
	return u.regionRepo.ListClient(ctx, filter)
}

func (u *regionUsecaseImpl) FindByPoint(ctx context.Context, lat, lng float64, layerID *uint64) (*qh.QHRegion, error) {
	if lat < -90 || lat > 90 {
		return nil, validation.New("latitude", "REGION_LATITUDE_INVALID", "invalid latitude")
	}
	if lng < -180 || lng > 180 {
		return nil, validation.New("longitude", "REGION_LONGITUDE_INVALID", "invalid longitude")
	}
	region, err := u.regionRepo.FindByPoint(ctx, lat, lng, layerID)
	if err != nil {
		return nil, fmt.Errorf("find by point failed: %w", err)
	}
	if region == nil {
		return nil, fmt.Errorf("no region found at location (%.6f, %.6f)", lat, lng)
	}
	return region, nil
}

func (u *regionUsecaseImpl) FindByBBox(ctx context.Context, minLng, minLat, maxLng, maxLat float64, layerID *uint64, labelID *uint64, page, size int) ([]qh.QHRegion, int64, error) {
	if minLng > maxLng || minLat > maxLat {
		return nil, 0, errors.New("invalid bbox coordinates")
	}

	offset := (page - 1) * size
	if offset < 0 {
		offset = 0
	}
	limit := size
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return u.regionRepo.FindByBBox(ctx, minLng, minLat, maxLng, maxLat, layerID, labelID, limit, offset)
}
