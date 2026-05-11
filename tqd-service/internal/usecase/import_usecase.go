package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/enums"
	"tqd/internal/geometry"
	"tqd/internal/interface/repo"
)

type ImportUsecase interface {
	Preview(ctx context.Context, geoJSON string, layerID uint64) (*dto.PreviewResult, error)
	EnqueueImportFromFile(ctx context.Context, req *dto.ImportFileRequest) (*dto.ImportEnqueueResult, error)
	// EnqueueImportFromReader — stream body vào file tạm (multipart HTTP); cùng logic xóa temp sau job
	EnqueueImportFromReader(ctx context.Context, src io.Reader, sourceFileName, fileFormat string, job *dto.ImportRequest) (*dto.ImportEnqueueResult, error)
	GetBatchStatus(ctx context.Context, batchID string) (*dto.BatchStatus, error)
	Rollback(ctx context.Context, batchID string) (int, error)
}

type importUsecaseImpl struct {
	regionRepo repo.RegionRepository
	labelRepo  repo.QHLabelRepository
	layerRepo  repo.LayerRepository
}

func NewImportUsecase(
	regionRepo repo.RegionRepository,
	labelRepo repo.QHLabelRepository,
	layerRepo repo.LayerRepository,
) ImportUsecase {
	return &importUsecaseImpl{
		regionRepo: regionRepo,
		labelRepo:  labelRepo,
		layerRepo:  layerRepo,
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

func (u *importUsecaseImpl) validateLayerForImport(ctx context.Context, layerID uint64) (*qh_domain.QHLayer, error) {
	if layerID == 0 {
		return nil, fmt.Errorf("layerId is required")
	}
	layer, err := u.layerRepo.GetByID(ctx, layerID)
	if err != nil {
		return nil, err
	}
	if layer == nil {
		return nil, fmt.Errorf("layer %d not found", layerID)
	}
	if layer.ImportStatus == enums.LayerImportStatusProcessing {
		return nil, fmt.Errorf("layer %d already has an import in progress", layerID)
	}
	return layer, nil
}

func (u *importUsecaseImpl) finalizeImportEnqueue(ctx context.Context, tmpPath, batchID, fileFormat string, jobReq *dto.ImportRequest) (*dto.ImportEnqueueResult, error) {
	if err := u.layerRepo.UpdateImportStatus(ctx, jobReq.LayerID, enums.LayerImportStatusProcessing, &batchID); err != nil {
		_ = os.Remove(tmpPath)
		return nil, fmt.Errorf("update layer import status: %w", err)
	}
	go u.runImportJob(tmpPath, batchID, strings.TrimSpace(fileFormat), jobReq)
	return &dto.ImportEnqueueResult{
		Success: true,
		BatchID: batchID,
		Status:  "processing",
		Message: "File accepted; import is running in the background",
	}, nil
}

func (u *importUsecaseImpl) EnqueueImportFromFile(ctx context.Context, req *dto.ImportFileRequest) (*dto.ImportEnqueueResult, error) {
	if len(req.FileContent) == 0 {
		return nil, fmt.Errorf("file_content is required")
	}
	if _, err := u.validateLayerForImport(ctx, req.LayerID); err != nil {
		return nil, err
	}

	batchID := generateBatchID()
	tmpRoot := filepath.Join(os.TempDir(), "tqd-import")
	tmpPath := filepath.Join(tmpRoot, batchID+".upload")
	if err := os.MkdirAll(tmpRoot, 0o755); err != nil {
		return nil, err
	}
	if err := writeUploadFileChunked(tmpPath, req.FileContent); err != nil {
		return nil, fmt.Errorf("save upload: %w", err)
	}

	jobReq := &dto.ImportRequest{
		LayerID:          req.LayerID,
		LabelField:       req.LabelField,
		LabelMappings:    req.LabelMappings,
		DefaultLabelID:   req.DefaultLabelID,
		SourceFileName:   req.SourceFileName,
		UserID:           req.UserID,
		ValidateGeometry: req.ValidateGeometry,
		AutoFixGeometry:  req.AutoFixGeometry,
		SkipInvalid:      req.SkipInvalid,
	}
	return u.finalizeImportEnqueue(ctx, tmpPath, batchID, req.FileFormat, jobReq)
}

func (u *importUsecaseImpl) EnqueueImportFromReader(ctx context.Context, src io.Reader, sourceFileName, fileFormat string, job *dto.ImportRequest) (*dto.ImportEnqueueResult, error) {
	if _, err := u.validateLayerForImport(ctx, job.LayerID); err != nil {
		return nil, err
	}

	batchID := generateBatchID()
	tmpRoot := filepath.Join(os.TempDir(), "tqd-import")
	tmpPath := filepath.Join(tmpRoot, batchID+".upload")
	if err := os.MkdirAll(tmpRoot, 0o755); err != nil {
		return nil, err
	}

	out, err := os.Create(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("create temp file: %w", err)
	}
	buf := make([]byte, 1024*1024)
	_, copyErr := io.CopyBuffer(out, src, buf)
	closeErr := out.Close()
	if copyErr != nil {
		_ = os.Remove(tmpPath)
		return nil, fmt.Errorf("write temp file: %w", copyErr)
	}
	if closeErr != nil {
		_ = os.Remove(tmpPath)
		return nil, fmt.Errorf("close temp file: %w", closeErr)
	}

	if job.SourceFileName == "" {
		job.SourceFileName = sourceFileName
	}
	return u.finalizeImportEnqueue(ctx, tmpPath, batchID, fileFormat, job)
}

func (u *importUsecaseImpl) runImportJob(tmpPath, batchID, fileFormat string, jobReq *dto.ImportRequest) {
	defer func() {
		if err := os.Remove(tmpPath); err != nil && !os.IsNotExist(err) {
			log.Printf("[import job] remove temp file %s: %v", tmpPath, err)
		}
	}()

	jobCtx := context.Background()
	result := &dto.ImportResult{
		BatchID:         batchID,
		LabelStatistics: make(map[uint64]int32),
		Failures:        make([]dto.ImportFailure, 0),
	}

	err := processImportFileByFormat(tmpPath, fileFormat, func(feat geoJSONFeature, index int) error {
		u.ingestFeature(jobCtx, batchID, index, feat, jobReq, result)
		return nil
	})

	if err != nil {
		log.Printf("[import job] batch=%s layer=%d error: %v", batchID, jobReq.LayerID, err)
		_ = u.layerRepo.UpdateImportStatus(jobCtx, jobReq.LayerID, enums.LayerImportStatusFailed, &batchID)
		return
	}

	if err := u.layerRepo.UpdateImportStatus(jobCtx, jobReq.LayerID, enums.LayerImportStatusDone, &batchID); err != nil {
		log.Printf("[import job] batch=%s update done: %v", batchID, err)
	}
	log.Printf("[import job] batch=%s layer=%d success=%d failed=%d", batchID, jobReq.LayerID, result.SuccessCount, result.FailedCount)
}

func (u *importUsecaseImpl) ingestFeature(ctx context.Context, batchID string, index int, feat geoJSONFeature, req *dto.ImportRequest, result *dto.ImportResult) {
	result.TotalFeatures++

	if feat.Geometry == nil {
		if !req.SkipInvalid {
			result.FailedCount++
			result.Failures = append(result.Failures, dto.ImportFailure{
				FeatureIndex: index,
				Reason:       "empty geometry",
			})
		}
		return
	}

	geomBytes, normErr := geometry.ToMultiPolygonGeoJSONBytes(feat.Geometry)
	if normErr != nil {
		if req.SkipInvalid {
			return
		}
		result.FailedCount++
		result.Failures = append(result.Failures, dto.ImportFailure{
			FeatureIndex: index,
			Reason:       fmt.Sprintf("geometry normalize: %v", normErr),
		})
		return
	}

	geoInfo, err := geometry.ParseGeometry(string(geomBytes))
	if err != nil && req.ValidateGeometry {
		result.FailedCount++
		result.Failures = append(result.Failures, dto.ImportFailure{
			FeatureIndex: index,
			Reason:       fmt.Sprintf("invalid geometry: %v", err),
		})
		return
	}
	area := 0.0
	if geoInfo != nil {
		area = geoInfo.Area
	}
	result.TotalAreaHa += area / 10000

	props := feat.Properties
	if props == nil {
		props = map[string]interface{}{}
	}

	labelID, err := u.resolveLabelID(ctx, props, req)
	if err != nil {
		if req.SkipInvalid {
			return
		}
		result.FailedCount++
		result.Failures = append(result.Failures, dto.ImportFailure{
			FeatureIndex: index,
			Reason:       fmt.Sprintf("label resolution: %v", err),
		})
		return
	}

	propsBytes, _ := json.Marshal(props)

	colorRed := uint8(200)
	colorGreen := uint8(200)
	colorBlue := uint8(200)
	if props, ok := props["red"]; ok {
		if val, ok := props.(float64); ok {
			colorRed = uint8(val)
		}
	}
	if props, ok := props["green"]; ok {
		if val, ok := props.(float64); ok {
			colorGreen = uint8(val)
		}
	}
	if props, ok := props["blue"]; ok {
		if val, ok := props.(float64); ok {
			colorBlue = uint8(val)
		}
	}

	region := &qh_domain.QHRegion{
		LayerID:            req.LayerID,
		LabelID:            labelID,
		Geometry:           qh_domain.RawGeometry{Raw: geomBytes},
		OriginalProperties: propsBytes,
		AreaSqm:            area,
		AreaHa:             area / 10000,
		SourceFile:         req.SourceFileName,
		ImportBatchID:      batchID,
		Status:             10,
		ProcessingStatus:   10,
		Version:            1,
		IsLatest:           true,
		ColorRed:           colorRed,
		ColorGreen:         colorGreen,
		ColorBlue:          colorBlue,
	}

	if name, ok := props["name"]; ok {
		region.Name = fmt.Sprintf("%v", name)
	}
	if displayName, ok := props["display_name"]; ok {
		region.DisplayName = fmt.Sprintf("%v", displayName)
	}
	if name, ok := props["ten_khu"]; ok && region.Name == "" {
		region.Name = fmt.Sprintf("%v", name)
	}
	if displayName, ok := props["ten_khu"]; ok && region.DisplayName == "" {
		region.DisplayName = fmt.Sprintf("%v", displayName)
	}

	if err := u.regionRepo.Create(ctx, region); err != nil {
		if req.SkipInvalid {
			return
		}
		result.FailedCount++
		result.Failures = append(result.Failures, dto.ImportFailure{
			FeatureIndex: index,
			Reason:       fmt.Sprintf("db insert: %v", err),
		})
		return
	}

	result.LabelStatistics[labelID]++
	result.SuccessCount++
	result.TotalRegionsCreated++
}

func (u *importUsecaseImpl) resolveLabelID(ctx context.Context, properties map[string]interface{}, req *dto.ImportRequest) (uint64, error) {
	if req.LabelField != "" {
		rawValue, ok := properties[req.LabelField]
		if ok {
			value := normalizeValue(rawValue)
			if labelID, ok := req.LabelMappings[value]; ok {
				return labelID, nil
			}
			label, err := u.labelRepo.GetByNameAndLayer(ctx, value, req.LayerID)
			if err == nil && label != nil {
				return label.ID, nil
			}
			return u.createNewLabel(ctx, value, req.LayerID)
		}
	}

	for _, val := range properties {
		strVal, ok := val.(string)
		if !ok || len(strVal) > 100 {
			continue
		}
		if isCoordinateOrNumeric(strVal) {
			continue
		}
		value := normalizeValue(strVal)
		if value == "" || value == "null" || value == "undefined" {
			continue
		}
		label, err := u.labelRepo.GetByNameAndLayer(ctx, value, req.LayerID)
		if err == nil && label != nil {
			return label.ID, nil
		}
		return u.createNewLabel(ctx, value, req.LayerID)
	}

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
	if batchID == "" {
		return nil, fmt.Errorf("batchId is required")
	}

	layer, err := u.layerRepo.GetByImportBatchID(ctx, batchID)
	if err != nil {
		return nil, err
	}

	regions, err := u.regionRepo.ListByImportBatchID(ctx, batchID)
	if err != nil {
		return nil, err
	}

	if layer == nil && len(regions) == 0 {
		return nil, fmt.Errorf("batch %s not found", batchID)
	}

	statusStr := "completed"
	errMsgs := []string(nil)
	if layer != nil {
		switch layer.ImportStatus {
		case enums.LayerImportStatusProcessing:
			statusStr = "processing"
		case enums.LayerImportStatusFailed:
			statusStr = "failed"
			errMsgs = []string{"import job failed"}
		case enums.LayerImportStatusDone:
			statusStr = "completed"
		}
	} else if len(regions) > 0 {
		// Batch cũ không gắn import_batch_id trên layer
		statusStr = "completed"
	}

	createdAt := time.Now().Format(time.RFC3339)
	completedAt := ""
	if len(regions) > 0 {
		createdAt = regions[0].CreatedAt.Format(time.RFC3339)
	}
	if statusStr == "completed" || statusStr == "failed" {
		completedAt = time.Now().Format(time.RFC3339)
	}

	return &dto.BatchStatus{
		BatchID:       batchID,
		Status:        statusStr,
		TotalFeatures: len(regions),
		SuccessCount:  len(regions),
		FailedCount:   0,
		CreatedAt:     createdAt,
		CompletedAt:   completedAt,
		ErrorMessages: errMsgs,
	}, nil
}

func (u *importUsecaseImpl) Rollback(ctx context.Context, batchID string) (int, error) {
	regions, err := u.regionRepo.ListByImportBatchID(ctx, batchID)
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

	layer, err := u.layerRepo.GetByImportBatchID(ctx, batchID)
	if err != nil {
		return deletedCount, err
	}
	if layer != nil {
		_ = u.layerRepo.ResetLayerImportTracking(ctx, layer.ID)
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
	words := strings.Split(rawName, " ")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

func isCoordinateOrNumeric(value string) bool {
	if _, err := fmt.Sscanf(value, "%f", new(float64)); err == nil && len(value) < 20 {
		return true
	}
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
