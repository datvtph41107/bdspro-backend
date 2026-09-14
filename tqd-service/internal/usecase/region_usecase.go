package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	_dto "common/domain/dto"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
	"tqd/internal/geometry"
	"tqd/internal/interface/repo"
)

// ErrLayerNotFound khi admin list region theo layer nhưng layer không tồn tại
var ErrLayerNotFound = errors.New("layer not found")

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
	Create(ctx context.Context, layerID uint64, name string, geometryBytes []byte, labelID uint64) (*qh_domain.QHRegion, error)
	Update(ctx context.Context, patch RegionUpdatePatch) (*qh_domain.QHRegion, error)
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHRegion, error)
	List(ctx context.Context, filter *repo.RegionFilter) ([]qh_domain.QHRegion, int64, error)
	// ListForAdminByLayer kiểm tra layer tồn tại rồi list regions (admin)
	ListForAdminByLayer(ctx context.Context, layerID uint64, filter *repo.RegionFilter) ([]qh_domain.QHRegion, int64, error)
	Delete(ctx context.Context, id uint64) error
	SyncRegion(ctx context.Context, layerID, labelID uint64) (*RegionSyncResult, error)
	UpdateProcessingStatus(ctx context.Context, regionID uint64, status enums.ProcessingStatus, userID uint64, notes string) error

	// Label Operations (1-1)
	AssignLabel(ctx context.Context, regionID uint64, labelID uint64, userID uint64, sourceFile string) error
	GetLabelByRegion(ctx context.Context, regionID uint64) (*qh_domain.QHLabel, error)

	// Client Operations
	ListClient(ctx context.Context, layerID uint64, bBox *BBox, page, size uint32) ([]qh_domain.QHRegion, int64, error)
	FindByPoint(ctx context.Context, lat, lng float64, layerID *uint64) (*qh_domain.QHRegion, error)
	FindByBBox(ctx context.Context, minLng, minLat, maxLng, maxLat float64, layerID *uint64, labelID *uint64, page, size int) ([]qh_domain.QHRegion, int64, error)
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
		return nil, errors.New("geometry is empty")
	}
	var root map[string]interface{}
	if err := json.Unmarshal(trim, &root); err != nil {
		return nil, fmt.Errorf("geometry is not valid JSON: %w", err)
	}
	t, _ := root["type"].(string)
	t = strings.TrimSpace(t)
	switch t {
	case "Feature":
		g, ok := root["geometry"].(map[string]interface{})
		if !ok || g == nil {
			return nil, errors.New("GeoJSON Feature missing geometry object")
		}
		return g, nil
	case "Polygon", "MultiPolygon":
		return root, nil
	default:
		return nil, fmt.Errorf("unsupported GeoJSON type %q (expected Feature, Polygon, or MultiPolygon)", t)
	}
}

func (u *regionUsecaseImpl) Create(
	ctx context.Context,
	layerID uint64,
	name string,
	geometryBytes []byte,
	labelID uint64,
) (*qh_domain.QHRegion, error) {
	if layerID == 0 {
		return nil, regionLayerIDRequired()
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, regionNameRequired()
	}

	layer, err := u.layerRepo.GetByID(ctx, layerID)
	if err != nil {
		return nil, regionInternal(
			err,
			"tqd.region.layer_lookup_failed",
		)
	}
	if layer == nil {
		return nil, regionLayerNotFound(layerID)
	}

	geomObj, err := extractGeometryObjectFromBytes(geometryBytes)
	if err != nil {
		return nil, regionGeometryInvalid(err)
	}

	geomBytes, err := geometry.ToMultiPolygonGeoJSONBytes(geomObj)
	if err != nil {
		return nil, regionGeometryInvalid(err)
	}

	geoInfo, err := geometry.ParseGeometry(string(geomBytes))
	if err != nil {
		return nil, regionGeometryInvalid(err)
	}

	areaSqm := 0.0
	if geoInfo != nil {
		areaSqm = geoInfo.Area
	}

	propsBytes, _ := json.Marshal(
		map[string]interface{}{
			"source": "admin_create",
		},
	)

	region := &qh_domain.QHRegion{
		LayerID:            layerID,
		Name:               name,
		DisplayName:        name,
		LabelID:            &labelID,
		Geometry:           qh_domain.RawGeometry{Raw: geomBytes},
		OriginalProperties: propsBytes,
		AreaSqm:            areaSqm,
		AreaHa:             areaSqm / 10000,
		Status:             enums.RegionStatusActive,
		ProcessingStatus:   enums.ProcessingStatusPending,
		Version:            1,
		IsLatest:           true,
	}

	if err := u.regionRepo.Create(ctx, region); err != nil {
		return nil, regionInternal(
			err,
			"tqd.region.create_failed",
		)
	}

	if region.ID == 0 {
		return nil, regionInternal(
			nil,
			"tqd.region.create_failed",
		)
	}

	created, err := u.regionRepo.GetByID(ctx, region.ID)
	if err != nil {
		return nil, regionInternal(
			err,
			"tqd.region.reload_failed",
		)
	}
	if created == nil {
		return nil, regionInternal(
			nil,
			"tqd.region.reload_failed",
		)
	}

	return created, nil
}

func (u *regionUsecaseImpl) Update(
	ctx context.Context,
	patch RegionUpdatePatch,
) (*qh_domain.QHRegion, error) {
	if patch.RegionID == 0 {
		return nil, regionIDRequired()
	}

	has := patch.Name != nil ||
		patch.DisplayName != nil ||
		patch.Description != nil ||
		patch.Status != nil ||
		patch.LabelID != nil

	if !has {
		return nil, regionUpdateRequired()
	}

	region, err := u.regionRepo.GetByID(
		ctx,
		patch.RegionID,
	)
	if err != nil {
		return nil, regionInternal(
			err,
			"tqd.region.lookup_failed",
		)
	}
	if region == nil {
		return nil, regionNotFound(patch.RegionID)
	}

	if patch.Name != nil {
		value := strings.TrimSpace(*patch.Name)
		if value == "" {
			return nil, regionNameEmpty()
		}
		region.Name = value
	}

	if patch.DisplayName != nil {
		region.DisplayName = strings.TrimSpace(
			*patch.DisplayName,
		)
	}

	if patch.Description != nil {
		region.Description = *patch.Description
	}

	if patch.Status != nil {
		if !patch.Status.IsValid() {
			return nil, regionStatusInvalid(
				uint32(*patch.Status),
			)
		}
		region.Status = *patch.Status
	}

	if patch.LabelID != nil {
		labelID := *patch.LabelID

		if labelID == 0 {
			region.LabelID = nil
		} else {
			label, lookupErr := u.labelRepo.GetByID(
				ctx,
				labelID,
			)
			if lookupErr != nil {
				return nil, regionInternal(
					lookupErr,
					"tqd.region.label_lookup_failed",
				)
			}
			if label == nil {
				return nil, regionLabelNotFound(
					labelID,
				)
			}

			region.LabelID = &labelID
		}
	}

	if err := u.regionRepo.Update(ctx, region); err != nil {
		return nil, regionInternal(
			err,
			"tqd.region.update_failed",
		)
	}

	if region.LabelID != nil {
		u.labelRepo.EnsureLayerLink(
			ctx,
			*region.LabelID,
			region.LayerID,
		)
	}

	updated, err := u.regionRepo.GetByID(
		ctx,
		patch.RegionID,
	)
	if err != nil {
		return nil, regionInternal(
			err,
			"tqd.region.reload_failed",
		)
	}
	if updated == nil {
		return nil, regionInternal(
			nil,
			"tqd.region.reload_failed",
		)
	}

	return updated, nil
}

func (u *regionUsecaseImpl) GetByID(ctx context.Context, id uint64) (*qh_domain.QHRegion, error) {
	region, err := u.regionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get region failed: %w", err)
	}
	if region == nil {
		return nil, fmt.Errorf("region %d not found", id)
	}
	return region, nil
}

func (u *regionUsecaseImpl) List(ctx context.Context, filter *repo.RegionFilter) ([]qh_domain.QHRegion, int64, error) {
	return u.regionRepo.List(ctx, filter)
}

func (u *regionUsecaseImpl) ListForAdminByLayer(ctx context.Context, layerID uint64, filter *repo.RegionFilter) ([]qh_domain.QHRegion, int64, error) {
	if layerID == 0 {
		return nil, 0, errors.New("layerId is required")
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
		return fmt.Errorf("region %d not found", id)
	}
	return u.regionRepo.Delete(ctx, id)
}

func (u *regionUsecaseImpl) SyncRegion(
	ctx context.Context,
	layerID uint64,
	labelID uint64,
) (*RegionSyncResult, error) {
	if layerID == 0 {
		return nil, regionLayerIDRequired()
	}
	if labelID == 0 {
		return nil, regionLabelIDRequired()
	}

	layer, err := u.layerRepo.GetByID(ctx, layerID)
	if err != nil {
		return nil, regionInternal(
			err,
			"tqd.region.layer_lookup_failed",
		)
	}
	if layer == nil {
		return nil, regionLayerNotFound(layerID)
	}

	label, err := u.labelRepo.GetByID(ctx, labelID)
	if err != nil {
		return nil, regionInternal(
			err,
			"tqd.region.label_lookup_failed",
		)
	}
	if label == nil {
		return nil, regionLabelNotFound(labelID)
	}

	ref, err := u.regionRepo.GetSyncReferenceByLayerLabel(
		ctx,
		layerID,
		labelID,
	)
	if err != nil {
		return nil, regionInternal(
			err,
			"tqd.region.sync_reference_lookup_failed",
		)
	}
	if ref == nil {
		return nil, regionSyncReferenceNotFound(
			layerID,
			labelID,
		)
	}

	total, err :=
		u.regionRepo.SyncLandUseAndLegendByLayerLabel(
			ctx,
			layerID,
			labelID,
			ref.LandUseID,
			ref.LegendID,
		)
	if err != nil {
		return nil, regionInternal(
			err,
			"tqd.region.sync_failed",
		)
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
		return fmt.Errorf("region %d not found", regionID)
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
		return fmt.Errorf("region %d not found", regionID)
	}

	// Kiểm tra label tồn tại
	label, err := u.labelRepo.GetByID(ctx, labelID)
	if err != nil {
		return fmt.Errorf("get label failed: %w", err)
	}
	if label == nil {
		return fmt.Errorf("label %d not found", labelID)
	}

	// Gán label trực tiếp vào region
	region.LabelID = &labelID
	return u.regionRepo.Update(ctx, region)
}

func (u *regionUsecaseImpl) GetLabelByRegion(ctx context.Context, regionID uint64) (*qh_domain.QHLabel, error) {
	region, err := u.regionRepo.GetByID(ctx, regionID)
	if err != nil {
		return nil, fmt.Errorf("get region failed: %w", err)
	}
	if region == nil {
		return nil, fmt.Errorf("region %d not found", regionID)
	}
	if region.LabelID == nil {
		return nil, nil // Chưa có label
	}
	return u.labelRepo.GetByID(ctx, *region.LabelID)
}

// Client
func (u *regionUsecaseImpl) ListClient(ctx context.Context, layerID uint64, bBox *BBox, page, size uint32) ([]qh_domain.QHRegion, int64, error) {
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

func (u *regionUsecaseImpl) FindByPoint(ctx context.Context, lat, lng float64, layerID *uint64) (*qh_domain.QHRegion, error) {
	if lat < -90 || lat > 90 {
		return nil, errors.New("invalid latitude")
	}
	if lng < -180 || lng > 180 {
		return nil, errors.New("invalid longitude")
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

func (u *regionUsecaseImpl) FindByBBox(ctx context.Context, minLng, minLat, maxLng, maxLat float64, layerID *uint64, labelID *uint64, page, size int) ([]qh_domain.QHRegion, int64, error) {
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
