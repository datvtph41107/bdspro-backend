package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/geometry"
	"tqd/internal/interface/repo"
)

type ImportUsecase interface {
	Preview(ctx context.Context, geoJSON string, layerID uint64) (*dto.PreviewResult, error)
	Import(ctx context.Context, req *dto.ImportRequest) (*dto.ImportResult, error)
	GetBatchStatus(ctx context.Context, batchID string) (*dto.BatchStatus, error)
	Rollback(ctx context.Context, batchID string) (int, error)
}

type importUsecaseImpl struct {
	regionRepo      repo.RegionRepository
	labelRepo       repo.QHLabelRepository
	regionLabelRepo repo.RegionLabelRepository
}

func NewImportUsecase(
	regionRepo repo.RegionRepository,
	labelRepo repo.QHLabelRepository,
	regionLabelRepo repo.RegionLabelRepository,
) ImportUsecase {
	return &importUsecaseImpl{
		regionRepo:      regionRepo,
		labelRepo:       labelRepo,
		regionLabelRepo: regionLabelRepo,
	}
}

func (u *importUsecaseImpl) Preview(ctx context.Context, geoJSON string, layerID uint64) (*dto.PreviewResult, error) {
	var fc struct {
		Features []struct {
			Geometry   map[string]interface{} `json:"geometry"`
			Properties map[string]interface{} `json:"properties"`
		} `json:"features"`
	}

	if err := json.Unmarshal([]byte(geoJSON), &fc); err != nil {
		return nil, fmt.Errorf("invalid GeoJSON: %w", err)
	}

	fieldCount := make(map[string]int)
	fieldValues := make(map[string]map[string]bool)

	for _, feat := range fc.Features {
		for k, v := range feat.Properties {
			fieldCount[k]++
			if fieldValues[k] == nil {
				fieldValues[k] = make(map[string]bool)
			}
			fieldValues[k][fmt.Sprintf("%v", v)] = true
		}
	}

	var fieldStats []dto.FieldStat
	availableFields := make([]string, 0, len(fieldCount))
	for field, count := range fieldCount {
		availableFields = append(availableFields, field)
		distinct := make([]string, 0, len(fieldValues[field]))
		for v := range fieldValues[field] {
			distinct = append(distinct, v)
		}
		fieldStats = append(fieldStats, dto.FieldStat{
			FieldName:      field,
			Coverage:       count,
			DistinctValues: distinct,
		})
	}

	var samples []dto.FeatureSample
	for i, feat := range fc.Features {
		if i >= 5 {
			break
		}
		geomType := "Unknown"
		if geom, ok := feat.Geometry["type"]; ok {
			geomType = fmt.Sprintf("%v", geom)
		}
		samples = append(samples, dto.FeatureSample{
			Index:        i,
			GeometryType: geomType,
			Properties:   feat.Properties,
		})
	}

	return &dto.PreviewResult{
		TotalFeatures:   len(fc.Features),
		AvailableFields: availableFields,
		FieldStats:      fieldStats,
		Samples:         samples,
	}, nil
}

func (u *importUsecaseImpl) Import(ctx context.Context, req *dto.ImportRequest) (*dto.ImportResult, error) {
	var fc struct {
		Features []struct {
			Geometry   map[string]interface{} `json:"geometry"`
			Properties map[string]interface{} `json:"properties"`
		} `json:"features"`
	}

	if err := json.Unmarshal([]byte(req.GeoJSON), &fc); err != nil {
		return nil, fmt.Errorf("invalid GeoJSON: %w", err)
	}

	batchID := generateBatchID()
	result := &dto.ImportResult{
		BatchID:         batchID,
		TotalFeatures:   len(fc.Features),
		LabelStatistics: make(map[uint64]int32),
		Failures:        make([]dto.ImportFailure, 0),
	}

	for i, feat := range fc.Features {
		// Validate geometry
		if feat.Geometry == nil {
			if !req.SkipInvalid {
				result.FailedCount++
				result.Failures = append(result.Failures, dto.ImportFailure{
					FeatureIndex: i,
					Reason:       "empty geometry",
				})
			}
			continue
		}

		// Tính diện tích
		geomBytes, _ := json.Marshal(feat.Geometry)
		geoInfo, err := geometry.ParseGeometry(string(geomBytes))
		if err != nil && req.ValidateGeometry {
			result.FailedCount++
			result.Failures = append(result.Failures, dto.ImportFailure{
				FeatureIndex: i,
				Reason:       fmt.Sprintf("invalid geometry: %v", err),
			})
			continue
		}
		area := 0.0
		if geoInfo != nil {
			area = geoInfo.Area
		}
		result.TotalAreaHa += area / 10000

		// Xác định label
		labelID, err := u.resolveLabelID(ctx, feat.Properties, req)
		if err != nil {
			if req.SkipInvalid {
				continue
			}
			result.FailedCount++
			result.Failures = append(result.Failures, dto.ImportFailure{
				FeatureIndex: i,
				Reason:       fmt.Sprintf("label resolution: %v", err),
			})
			continue
		}

		// Tạo region
		propsBytes, _ := json.Marshal(feat.Properties)
		region := &qh_domain.QHRegion{
			LayerID:            req.LayerID,
			Geometry:           qh_domain.RawGeometry{Raw: geomBytes},
			OriginalProperties: propsBytes,
			AreaSqm:            area,
			AreaHa:             area / 10000,
			SourceFile:         req.SourceFileName,
			ImportBatchID:      batchID,
			Status:             10, // active
			ProcessingStatus:   10, // THÊM: pending
			Version:            1,
			IsLatest:           true,
		}

		// Lấy name từ properties nếu có
		if name, ok := feat.Properties["name"]; ok {
			region.Name = fmt.Sprintf("%v", name)
		}
		if displayName, ok := feat.Properties["display_name"]; ok {
			region.DisplayName = fmt.Sprintf("%v", displayName)
		}

		if err := u.regionRepo.Create(ctx, region); err != nil {
			if req.SkipInvalid {
				continue
			}
			result.FailedCount++
			result.Failures = append(result.Failures, dto.ImportFailure{
				FeatureIndex: i,
				Reason:       fmt.Sprintf("db insert: %v", err),
			})
			continue
		}

		// Gán label
		if err := u.regionLabelRepo.AddLabels(ctx, region.ID, []uint64{labelID}, req.UserID, req.SourceFileName, ""); err != nil {
			// Log error but continue
		}

		result.LabelStatistics[labelID]++
		result.SuccessCount++
		result.TotalRegionsCreated++
	}

	return result, nil
}

func (u *importUsecaseImpl) resolveLabelID(ctx context.Context, properties map[string]interface{}, req *dto.ImportRequest) (uint64, error) {
	if req.LabelField == "" {
		if req.DefaultLabelID != nil {
			return *req.DefaultLabelID, nil
		}
		return 0, fmt.Errorf("no label_field specified")
	}

	rawValue, ok := properties[req.LabelField]
	if !ok {
		if req.DefaultLabelID != nil {
			return *req.DefaultLabelID, nil
		}
		return 0, fmt.Errorf("field '%s' not found", req.LabelField)
	}

	value := strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", rawValue)))

	if labelID, ok := req.LabelMappings[value]; ok {
		return labelID, nil
	}

	label, err := u.labelRepo.GetByNameAndLayer(ctx, value, req.LayerID)
	if err == nil && label != nil {
		return label.ID, nil
	}

	// Auto create label
	newLabel := &qh_domain.QHLabel{
		LayerID:     req.LayerID,
		Name:        value,
		DisplayName: value,
		Color:       generateColorFromString(value),
		FillOpacity: 0.6,
		StrokeColor: "#000000",
		StrokeWidth: 1,
		IsVisible:   true,
		Status:      10,
	}
	if err := u.labelRepo.Create(ctx, newLabel); err != nil {
		return 0, fmt.Errorf("failed to create label: %w", err)
	}
	return newLabel.ID, nil
}

func (u *importUsecaseImpl) GetBatchStatus(ctx context.Context, batchID string) (*dto.BatchStatus, error) {
	regions, _, err := u.regionRepo.List(ctx, &repo.RegionFilter{
		Search: &batchID,
	})
	if err != nil {
		return nil, err
	}

	if len(regions) == 0 {
		return nil, fmt.Errorf("batch %s not found", batchID)
	}

	return &dto.BatchStatus{
		BatchID:       batchID,
		Status:        "completed",
		TotalFeatures: len(regions),
		SuccessCount:  len(regions),
		FailedCount:   0,
		CreatedAt:     regions[0].CreatedAt.Format(time.RFC3339),
		CompletedAt:   time.Now().Format(time.RFC3339),
		ErrorMessages: []string{},
	}, nil
}

func (u *importUsecaseImpl) Rollback(ctx context.Context, batchID string) (int, error) {
	regions, _, err := u.regionRepo.List(ctx, &repo.RegionFilter{
		Search: &batchID,
	})
	if err != nil {
		return 0, err
	}

	deletedCount := 0
	for _, r := range regions {
		if err := u.regionRepo.Delete(ctx, r.ID); err != nil {
			continue
		}
		deletedCount++
	}

	return deletedCount, nil
}

func generateBatchID() string {
	hash := md5.Sum([]byte(time.Now().String()))
	return hex.EncodeToString(hash[:])[:16]
}

func generateColorFromString(s string) string {
	hash := md5.Sum([]byte(s))
	return fmt.Sprintf("#%02X%02X%02X", hash[0], hash[1], hash[2])
}
