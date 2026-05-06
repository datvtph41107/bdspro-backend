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
	regionRepo repo.RegionRepository
	labelRepo  repo.QHLabelRepository
}

func NewImportUsecase(
	regionRepo repo.RegionRepository,
	labelRepo repo.QHLabelRepository,
) ImportUsecase {
	return &importUsecaseImpl{
		regionRepo: regionRepo,
		labelRepo:  labelRepo,
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

		// Tạo region - THÊM label_id TRỰC TIẾP
		propsBytes, _ := json.Marshal(feat.Properties)

		// Lấy màu từ properties nếu có
		colorRed := uint8(200)
		colorGreen := uint8(200)
		colorBlue := uint8(200)
		if props, ok := feat.Properties["red"]; ok {
			if val, ok := props.(float64); ok {
				colorRed = uint8(val)
			}
		}
		if props, ok := feat.Properties["green"]; ok {
			if val, ok := props.(float64); ok {
				colorGreen = uint8(val)
			}
		}
		if props, ok := feat.Properties["blue"]; ok {
			if val, ok := props.(float64); ok {
				colorBlue = uint8(val)
			}
		}

		region := &qh_domain.QHRegion{
			LayerID:            req.LayerID,
			LabelID:            labelID, // ← THÊM LABEL_ID TRỰC TIẾP
			Geometry:           qh_domain.RawGeometry{Raw: geomBytes},
			OriginalProperties: propsBytes,
			AreaSqm:            area,
			AreaHa:             area / 10000,
			SourceFile:         req.SourceFileName,
			ImportBatchID:      batchID,
			Status:             10, // active
			ProcessingStatus:   10, // pending
			Version:            1,
			IsLatest:           true,
			ColorRed:           colorRed,
			ColorGreen:         colorGreen,
			ColorBlue:          colorBlue,
		}

		// Lấy name từ properties nếu có
		if name, ok := feat.Properties["name"]; ok {
			region.Name = fmt.Sprintf("%v", name)
		}
		if displayName, ok := feat.Properties["display_name"]; ok {
			region.DisplayName = fmt.Sprintf("%v", displayName)
		}
		if name, ok := feat.Properties["ten_khu"]; ok && region.Name == "" {
			region.Name = fmt.Sprintf("%v", name)
		}
		if displayName, ok := feat.Properties["ten_khu"]; ok && region.DisplayName == "" {
			region.DisplayName = fmt.Sprintf("%v", displayName)
		}

		// Insert region (đã có label_id)
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

		// KHÔNG CẦN gán label riêng nữa (vì đã có label_id trong region)
		// ĐÃ XÓA: u.regionLabelRepo.AddLabels(...)

		result.LabelStatistics[labelID]++
		result.SuccessCount++
		result.TotalRegionsCreated++
	}

	return result, nil
}

func (u *importUsecaseImpl) resolveLabelID(ctx context.Context, properties map[string]interface{}, req *dto.ImportRequest) (uint64, error) {
	// Ưu tiên dùng LabelField từ request
	if req.LabelField != "" {
		rawValue, ok := properties[req.LabelField]
		if ok {
			value := normalizeValue(rawValue)
			// Kiểm tra mapping
			if labelID, ok := req.LabelMappings[value]; ok {
				return labelID, nil
			}
			// Tìm label theo tên
			label, err := u.labelRepo.GetByNameAndLayer(ctx, value, req.LayerID)
			if err == nil && label != nil {
				return label.ID, nil
			}
			// Tạo label mới
			return u.createNewLabel(ctx, value, req.LayerID)
		}
	}

	// Fallback: tìm bất kỳ field nào có giá trị có nghĩa
	for _, val := range properties {
		// Bỏ qua các field không phải string hoặc quá dài
		strVal, ok := val.(string)
		if !ok || len(strVal) > 100 {
			continue
		}
		// Bỏ qua các field tọa độ, id, số
		if isCoordinateOrNumeric(strVal) {
			continue
		}
		value := normalizeValue(strVal)
		if value == "" || value == "null" || value == "undefined" {
			continue
		}
		// Tìm hoặc tạo label
		label, err := u.labelRepo.GetByNameAndLayer(ctx, value, req.LayerID)
		if err == nil && label != nil {
			return label.ID, nil
		}
		return u.createNewLabel(ctx, value, req.LayerID)
	}

	// Default label
	if req.DefaultLabelID != nil {
		return *req.DefaultLabelID, nil
	}
	return u.createNewLabel(ctx, "Chưa phân loại", req.LayerID)
}

func (u *importUsecaseImpl) createNewLabel(ctx context.Context, name string, layerID uint64) (uint64, error) {
	normalizedName := normalizeLabelName(name)
	newLabel := &qh_domain.QHLabel{
		LayerID:     layerID,
		Name:        normalizedName,
		DisplayName: normalizedName,
		Color:       generateColorFromString(normalizedName),
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

func normalizeValue(v interface{}) string {
	if v == nil {
		return ""
	}
	str := strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", v)))
	return str
}

func normalizeLabelName(rawName string) string {
	if rawName == "" {
		return "Chưa phân loại"
	}
	// Capitalize first letter of each word
	words := strings.Split(rawName, " ")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

func isCoordinateOrNumeric(value string) bool {
	// Kiểm tra nếu là số
	if _, err := fmt.Sscanf(value, "%f", new(float64)); err == nil && len(value) < 20 {
		return true
	}
	// Kiểm tra nếu là tọa độ (chứa dấu phẩy hoặc khoảng trắng và số)
	if strings.ContainsAny(value, ", ") {
		parts := strings.FieldsFunc(value, func(r rune) bool {
			return r == ',' || r == ' '
		})
		numberCount := 0
		for _, p := range parts {
			if _, err := fmt.Sscanf(p, "%f", new(float64)); err == nil {
				numberCount++
			}
		}
		if numberCount >= 2 && numberCount <= 4 {
			return true
		}
	}
	return false
}
