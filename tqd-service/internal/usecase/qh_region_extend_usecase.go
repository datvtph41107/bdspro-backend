package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	_dto "common/domain/dto"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
	"tqd/internal/geometry"
	"tqd/internal/interface/repo"
)

type RegionExtendCreateInput struct {
	LayerID            uint64
	Name               string
	DisplayName        string
	Description        string
	LabelID            *uint64
	Geometry           []byte
	LegalDoc           string
	PlanningName       string
	OriginalProperties []byte
	SourceFile         string
	ImportBatchID      string
	Status             *uint32
	MergeSourceIDs     []uint64
}

type RegionExtendUpdatePatch struct {
	ID                 uint64
	LayerID            *uint64
	Name               *string
	DisplayName        *string
	Description        *string
	LabelID            *uint64
	Geometry           []byte
	LegalDoc           *string
	PlanningName       *string
	OriginalProperties []byte
	SourceFile         *string
	ImportBatchID      *string
	Status             *uint32
}

type RegionExtendUsecase interface {
	Create(ctx context.Context, input RegionExtendCreateInput) (*qh_domain.QHRegionExtend, error)
	Update(ctx context.Context, patch RegionExtendUpdatePatch) (*qh_domain.QHRegionExtend, error)
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, layerID uint64, geomType *string, pagable *_dto.Pagable) ([]qh_domain.QHRegionExtend, int64, error)
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHRegionExtend, error)
}

type regionExtendUsecaseImpl struct {
	extendRepo repo.RegionExtendRepository
}

func NewRegionExtendUsecase(extendRepo repo.RegionExtendRepository) RegionExtendUsecase {
	return &regionExtendUsecaseImpl{extendRepo: extendRepo}
}

func buildExtendGeometry(raw []byte) (qh_domain.RawExtendGeometry, string, error) {
	if len(raw) == 0 {
		return qh_domain.RawExtendGeometry{}, "", fmt.Errorf("geometry is required")
	}
	var geom map[string]interface{}
	if err := json.Unmarshal(raw, &geom); err != nil {
		return qh_domain.RawExtendGeometry{}, "", fmt.Errorf("geometry is not valid JSON: %w", err)
	}
	geomBytes, err := geometry.ToRawGeoJSONBytes(geom)
	if err != nil {
		return qh_domain.RawExtendGeometry{}, "", err
	}
	geomType := geometry.GetGeomType(geom)
	if strings.TrimSpace(geomType) == "" {
		return qh_domain.RawExtendGeometry{}, "", fmt.Errorf("geometry type is required")
	}
	return qh_domain.RawExtendGeometry{Raw: geomBytes}, geomType, nil
}

func uniqueUint64IDs(ids []uint64) []uint64 {
	if len(ids) == 0 {
		return nil
	}
	seen := make(map[uint64]struct{}, len(ids))
	out := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func (u *regionExtendUsecaseImpl) Create(ctx context.Context, input RegionExtendCreateInput) (*qh_domain.QHRegionExtend, error) {
	if input.LayerID == 0 {
		return nil, fmt.Errorf("layerId is required")
	}

	mergeSourceIDs := uniqueUint64IDs(input.MergeSourceIDs)
	if len(mergeSourceIDs) == 0 && len(input.Geometry) == 0 {
		return nil, fmt.Errorf("geometry is required when mergeSourceIds is empty")
	}

	status := enums.RegionStatusActive
	if input.Status != nil {
		status = enums.RegionStatus(*input.Status)
	}

	record := &qh_domain.QHRegionExtend{
		LayerID:            input.LayerID,
		Name:               strings.TrimSpace(input.Name),
		DisplayName:        strings.TrimSpace(input.DisplayName),
		Description:        input.Description,
		LabelID:            input.LabelID,
		LegalDoc:           input.LegalDoc,
		PlanningName:       input.PlanningName,
		OriginalProperties: input.OriginalProperties,
		SourceFile:         input.SourceFile,
		ImportBatchID:      input.ImportBatchID,
		Status:             status,
		Version:            1,
		IsLatest:           true,
	}
	if record.DisplayName == "" {
		record.DisplayName = record.Name
	}

	if len(mergeSourceIDs) > 0 {
		var extraGeom []byte
		if len(input.Geometry) > 0 {
			geom, _, err := buildExtendGeometry(input.Geometry)
			if err != nil {
				return nil, err
			}
			extraGeom = geom.Raw
		}
		if err := u.extendRepo.CreateMergedReplaceSources(ctx, record, mergeSourceIDs, extraGeom); err != nil {
			return nil, fmt.Errorf("create merged region extend failed: %w", err)
		}
		return u.GetByID(ctx, record.ID)
	}

	geom, geomType, err := buildExtendGeometry(input.Geometry)
	if err != nil {
		return nil, err
	}
	record.GeomType = geomType
	record.Geometry = geom
	if err := u.extendRepo.Create(ctx, record); err != nil {
		return nil, fmt.Errorf("create region extend failed: %w", err)
	}
	return u.GetByID(ctx, record.ID)
}

func (u *regionExtendUsecaseImpl) Update(ctx context.Context, patch RegionExtendUpdatePatch) (*qh_domain.QHRegionExtend, error) {
	if patch.ID == 0 {
		return nil, fmt.Errorf("id is required")
	}
	record, err := u.GetByID(ctx, patch.ID)
	if err != nil {
		return nil, err
	}
	if patch.LayerID != nil {
		if *patch.LayerID == 0 {
			return nil, fmt.Errorf("layerId is required")
		}
		record.LayerID = *patch.LayerID
	}
	if patch.Name != nil {
		record.Name = strings.TrimSpace(*patch.Name)
	}
	if patch.DisplayName != nil {
		record.DisplayName = strings.TrimSpace(*patch.DisplayName)
	}
	if patch.Description != nil {
		record.Description = *patch.Description
	}
	if patch.LabelID != nil {
		if *patch.LabelID == 0 {
			record.LabelID = nil
		} else {
			record.LabelID = patch.LabelID
		}
	}
	if len(patch.Geometry) > 0 {
		geom, geomType, gerr := buildExtendGeometry(patch.Geometry)
		if gerr != nil {
			return nil, gerr
		}
		record.Geometry = geom
		record.GeomType = geomType
	}
	if patch.LegalDoc != nil {
		record.LegalDoc = *patch.LegalDoc
	}
	if patch.PlanningName != nil {
		record.PlanningName = *patch.PlanningName
	}
	if patch.OriginalProperties != nil {
		record.OriginalProperties = patch.OriginalProperties
	}
	if patch.SourceFile != nil {
		record.SourceFile = *patch.SourceFile
	}
	if patch.ImportBatchID != nil {
		record.ImportBatchID = *patch.ImportBatchID
	}
	if patch.Status != nil {
		record.Status = enums.RegionStatus(*patch.Status)
	}

	if err := u.extendRepo.Update(ctx, record); err != nil {
		return nil, fmt.Errorf("update region extend failed: %w", err)
	}
	return u.GetByID(ctx, patch.ID)
}

func (u *regionExtendUsecaseImpl) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return fmt.Errorf("id is required")
	}
	if _, err := u.GetByID(ctx, id); err != nil {
		return err
	}
	return u.extendRepo.Delete(ctx, id)
}

func (u *regionExtendUsecaseImpl) List(ctx context.Context, layerID uint64, geomType *string, pagable *_dto.Pagable) ([]qh_domain.QHRegionExtend, int64, error) {
	if layerID == 0 {
		return nil, 0, fmt.Errorf("layerId is required")
	}
	filter := &repo.RegionExtendFilter{
		LayerID:  &layerID,
		GeomType: geomType,
		Pagable:  pagable,
	}
	return u.extendRepo.List(ctx, filter)
}

func (u *regionExtendUsecaseImpl) GetByID(ctx context.Context, id uint64) (*qh_domain.QHRegionExtend, error) {
	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}
	rec, err := u.extendRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, fmt.Errorf("region extend %d not found", id)
	}
	return rec, nil
}
