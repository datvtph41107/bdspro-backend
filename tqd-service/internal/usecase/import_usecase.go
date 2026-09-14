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
	"runtime"
	"strings"
	"sync"
	"time"

	_dto "common/domain/dto"
	_utils "common/utils"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/enums"
	"tqd/internal/geometry"
	"tqd/internal/interface/repo"
)

// ─── tunables ──────────────────────────────────────────────────────────────
const (
	importBatchSize  = 500 // records per INSERT batch
	importWorkers    = 4   // parallel normalise/resolve goroutines
	importChanBuffer = importWorkers * importBatchSize * 2
)

// ─── interfaces ────────────────────────────────────────────────────────────

type ImportUsecase interface {
	Preview(ctx context.Context, geoJSON string, layerID uint64) (*dto.PreviewResult, error)
	EnqueueImportFromFile(ctx context.Context, req *dto.ImportFileRequest) (*dto.ImportEnqueueResult, error)
	EnqueueImportFromReader(ctx context.Context, src io.Reader, sourceFileName, fileFormat string, job *dto.ImportRequest) (*dto.ImportEnqueueResult, error)
	GetBatchStatus(ctx context.Context, batchID string) (*dto.BatchStatus, error)
	Rollback(ctx context.Context, batchID string) (int, error)
	ListImportRegionErrors(ctx context.Context, layerID uint64, pagable *_dto.Pagable) ([]qh_domain.QHRegionImportErrorLog, int64, error)
	RetryImportError(ctx context.Context, errorID uint64) (*dto.RetryImportErrorResult, error)
}

// ─── label cache (per-job, not shared across requests) ────────────────────

type labelCache struct {
	mu    sync.Mutex
	cache map[string]uint64 // normalizedName → labelID (không phụ thuộc layer)
}

func newLabelCache() *labelCache { return &labelCache{cache: make(map[string]uint64)} }

func (lc *labelCache) get(normalized string) (uint64, bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	id, ok := lc.cache[normalized]
	return id, ok
}

func (lc *labelCache) set(normalized string, id uint64) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.cache[normalized] = id
}

// ─── impl ──────────────────────────────────────────────────────────────────

type importUsecaseImpl struct {
	regionRepo       repo.RegionRepository
	regionExtendRepo repo.RegionExtendRepository
	labelRepo        repo.QHLabelRepository
	layerRepo        repo.LayerRepository
}

func NewImportUsecase(
	regionRepo repo.RegionRepository,
	regionExtendRepo repo.RegionExtendRepository,
	labelRepo repo.QHLabelRepository,
	layerRepo repo.LayerRepository,
) ImportUsecase {
	return &importUsecaseImpl{
		regionRepo:       regionRepo,
		regionExtendRepo: regionExtendRepo,
		labelRepo:        labelRepo,
		layerRepo:        layerRepo,
	}
}

// ─── Preview (unchanged) ───────────────────────────────────────────────────

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
		if g, ok := feat.Geometry["type"]; ok {
			geomType = fmt.Sprintf("%v", g)
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

// ─── Enqueue helpers ───────────────────────────────────────────────────────

func (u *importUsecaseImpl) validateLayerForImport(ctx context.Context, layerID uint64) (*qh_domain.QHLayer, error) {
	if layerID == 0 {
		return nil, importLayerIDRequired()
	}
	layer, err := u.layerRepo.GetByID(ctx, layerID)
	if err != nil {
		return nil, err
	}
	if layer == nil {
		return nil, importLayerNotFound(layerID)
	}
	if layer.ImportStatus == enums.LayerImportStatusProcessing {
		return nil, importAlreadyInProgress(layerID)
	}
	return layer, nil
}

func (u *importUsecaseImpl) finalizeImportEnqueue(ctx context.Context, tmpPath, batchID, fileFormat string, jobReq *dto.ImportRequest) (*dto.ImportEnqueueResult, error) {
	if err := u.layerRepo.UpdateImportStatus(ctx, jobReq.LayerID, enums.LayerImportStatusProcessing, &batchID); err != nil {
		// _ = os.Remove(tmpPath)
		return nil, fmt.Errorf("update layer import status: %w", err)
	}
	jobCtx := _utils.CloneContext(ctx)
	go u.runImportJob(jobCtx, tmpPath, batchID, strings.TrimSpace(fileFormat), jobReq)
	return &dto.ImportEnqueueResult{
		Code:    0,
		BatchID: batchID,
		Status:  "processing",
		Message: "Đã nhận file; import đang chạy nền. Dùng batchId gọi GetImportBatchStatus để xem kết quả.",
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
		// _ = os.Remove(tmpPath)
		return nil, fmt.Errorf("write temp file: %w", copyErr)
	}
	if closeErr != nil {
		// _ = os.Remove(tmpPath)
		return nil, fmt.Errorf("close temp file: %w", closeErr)
	}
	if job.SourceFileName == "" {
		job.SourceFileName = sourceFileName
	}
	return u.finalizeImportEnqueue(ctx, tmpPath, batchID, fileFormat, job)
}

// ─── Core import job ───────────────────────────────────────────────────────

// pendingRegion carries a fully-built region ready for batch insert.
type pendingRegion struct {
	region *qh_domain.QHRegion
	index  int
}

// pendingExtend carries a non-polygon feature routed to qh_region_extends.
type pendingExtend struct {
	record *qh_domain.QHRegionExtend
	index  int
}

// runImportJob orchestrates the pipeline:
//
//	file reader → feature channel → N worker goroutines (normalize + label)
//	             → region channel → batcher goroutine (CreateBatch)
func (u *importUsecaseImpl) runImportJob(ctx context.Context, tmpPath, batchID, fileFormat string, req *dto.ImportRequest) {
	defer func() {
		if err := os.Remove(tmpPath); err != nil && !os.IsNotExist(err) {
			log.Printf("[import] remove temp %s: %v", tmpPath, err)
		}
	}()

	result := &dto.ImportResult{
		BatchID:         batchID,
		LabelStatistics: make(map[uint64]int32),
		Failures:        make([]dto.ImportFailure, 0),
	}

	lc := newLabelCache()

	// workers = min(importWorkers, GOMAXPROCS)
	numWorkers := importWorkers
	if mp := runtime.GOMAXPROCS(0); mp < numWorkers {
		numWorkers = mp
	}

	// featCh: parsed features from file reader
	featCh := make(chan indexedFeature, importChanBuffer)
	// regionCh: polygon regions ready to batch-insert
	regionCh := make(chan pendingRegion, importChanBuffer)
	// extendCh: non-polygon features routed to qh_region_extends
	extendCh := make(chan pendingExtend, importChanBuffer)

	var workerWg sync.WaitGroup
	var mu sync.Mutex // protects result

	// ── Worker goroutines ────────────────────────────────────────────────
	for i := 0; i < numWorkers; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			for item := range featCh {
				// Route non-polygon geometry to extend table
				if item.feat.Geometry != nil && !geometry.IsPolygonType(item.feat.Geometry) {
					ext, failure := u.buildRegionExtend(ctx, batchID, item.index, item.feat, req, lc)
					switch {
					case failure != nil:
						mu.Lock()
						result.TotalFeatures++
						result.FailedCount++
						if failure.Reason != "" {
							result.Failures = append(result.Failures, *failure)
						}
						mu.Unlock()
					case ext != nil:
						mu.Lock()
						result.TotalFeatures++
						mu.Unlock()
						extendCh <- pendingExtend{record: ext, index: item.index}
					}
					continue
				}

				r, failure := u.buildRegion(ctx, batchID, item.index, item.feat, req, lc)
				switch {
				case failure != nil:
					// explicit error (Reason may be empty when SkipInvalid swallowed it,
					// but buildRegion only returns a non-nil failure when we want to count it)
					mu.Lock()
					result.TotalFeatures++
					result.FailedCount++
					if failure.Reason != "" {
						result.Failures = append(result.Failures, *failure)
					}
					mu.Unlock()
				case r != nil:
					// good region — count it, then send OUTSIDE the lock
					// so the batcher goroutine can drain regionCh freely
					mu.Lock()
					result.TotalFeatures++
					mu.Unlock()
					regionCh <- pendingRegion{region: r, index: item.index}
					// else: silently skipped (SkipInvalid), buildRegion returned nil,nil
				}
			}
		}()
	}

	// ── Batcher goroutine (polygon regions) ─────────────────────────────
	var batcherDone sync.WaitGroup
	batcherDone.Add(1)
	go func() {
		defer batcherDone.Done()
		batch := make([]*qh_domain.QHRegion, 0, importBatchSize)

		flush := func() {
			if len(batch) == 0 {
				return
			}
			if err := u.regionRepo.CreateBatch(ctx, batch); err != nil {
				mu.Lock()
				for range batch {
					result.FailedCount++
					result.Failures = append(result.Failures, dto.ImportFailure{
						Reason: fmt.Sprintf("batch insert: %v", err),
					})
				}
				mu.Unlock()
			} else {
				mu.Lock()
				for _, r := range batch {
					result.LabelStatistics[*r.LabelID]++
					result.SuccessCount++
					result.TotalRegionsCreated++
					result.TotalAreaHa += r.AreaHa // moved here, away from worker lock
				}
				mu.Unlock()
			}
			batch = batch[:0]
		}

		for pr := range regionCh {
			batch = append(batch, pr.region)
			if len(batch) >= importBatchSize {
				flush()
			}
		}
		flush() // flush remaining
	}()

	// ── Batcher goroutine (extend records) ──────────────────────────────
	batcherDone.Add(1)
	go func() {
		defer batcherDone.Done()
		batch := make([]*qh_domain.QHRegionExtend, 0, importBatchSize)

		flush := func() {
			if len(batch) == 0 {
				return
			}
			if err := u.regionExtendRepo.CreateBatch(ctx, batch); err != nil {
				mu.Lock()
				for range batch {
					result.FailedCount++
					result.Failures = append(result.Failures, dto.ImportFailure{
						Reason: fmt.Sprintf("extend batch insert: %v", err),
					})
				}
				mu.Unlock()
			} else {
				mu.Lock()
				result.SuccessCount += len(batch)
				mu.Unlock()
			}
			batch = batch[:0]
		}

		for pe := range extendCh {
			batch = append(batch, pe.record)
			if len(batch) >= importBatchSize {
				flush()
			}
		}
		flush()
	}()

	// ── Reader (main goroutine, sequential I/O) ──────────────────────────
	readErr := processImportFileByFormat(tmpPath, fileFormat, func(feat geoJSONFeature, index int) error {
		featCh <- indexedFeature{feat: feat, index: index}
		return nil
	})
	close(featCh)

	workerWg.Wait()
	close(regionCh)
	close(extendCh)
	batcherDone.Wait()

	// ── Finalise ────────────────────────────────────────────────────────
	if readErr != nil {
		log.Printf("[import] batch=%s read error: %v", batchID, readErr)
		_ = u.layerRepo.UpdateImportStatus(ctx, req.LayerID, enums.LayerImportStatusFailed, &batchID)
		return
	}
	if err := u.labelRepo.SyncRegionCountByLayerID(ctx, req.LayerID); err != nil {
		log.Printf("[import] batch=%s sync label region_count layer=%d: %v", batchID, req.LayerID, err)
	}
	if err := u.layerRepo.UpdateImportStatus(ctx, req.LayerID, enums.LayerImportStatusDone, &batchID); err != nil {
		log.Printf("[import] batch=%s update done: %v", batchID, err)
	}
	log.Printf("[import] batch=%s layer=%d total=%d success=%d failed=%d",
		batchID, req.LayerID, result.TotalFeatures, result.SuccessCount, result.FailedCount)
}

// indexedFeature pairs a parsed GeoJSON feature with its position in the file.
type indexedFeature struct {
	feat  geoJSONFeature
	index int
}

// buildRegion normalises geometry and resolves the label for one feature.
// Returns (region, nil) on success, (nil, &failure) on error, (nil, nil) when skipped silently.
func (u *importUsecaseImpl) buildRegion(ctx context.Context, batchID string, index int, feat geoJSONFeature, req *dto.ImportRequest, lc *labelCache) (*qh_domain.QHRegion, *dto.ImportFailure) {
	if feat.Geometry == nil {
		if req.SkipInvalid {
			return nil, nil
		}
		return nil, &dto.ImportFailure{FeatureIndex: index, Reason: "empty geometry"}
	}

	geomBytes, normErr := geometry.ToMultiPolygonGeoJSONBytes(feat.Geometry)
	if normErr != nil {
		if req.SkipInvalid {
			return nil, nil
		}
		return nil, &dto.ImportFailure{FeatureIndex: index, Reason: fmt.Sprintf("geometry normalize: %v", normErr)}
	}

	geoInfo, err := geometry.ParseGeometry(string(geomBytes))
	if err != nil && req.ValidateGeometry {
		return nil, &dto.ImportFailure{FeatureIndex: index, Reason: fmt.Sprintf("invalid geometry: %v", err)}
	}
	area := 0.0
	if geoInfo != nil {
		area = geoInfo.Area
	}

	props := feat.Properties
	if props == nil {
		props = map[string]interface{}{}
	}

	labelID, err := u.resolveLabelIDCached(ctx, props, req, lc)
	if err != nil {
		if req.SkipInvalid {
			return nil, nil
		}
		return nil, &dto.ImportFailure{FeatureIndex: index, Reason: fmt.Sprintf("label resolution: %v", err)}
	}

	propsBytes, _ := json.Marshal(props)

	// colorRed := uint8(200)
	// colorGreen := uint8(200)
	// colorBlue := uint8(200)
	// if v, ok := props["red"]; ok {
	// 	if f, ok := v.(float64); ok {
	// 		colorRed = uint8(f)
	// 	}
	// }
	// if v, ok := props["green"]; ok {
	// 	if f, ok := v.(float64); ok {
	// 		colorGreen = uint8(f)
	// 	}
	// }
	// if v, ok := props["blue"]; ok {
	// 	if f, ok := v.(float64); ok {
	// 		colorBlue = uint8(f)
	// 	}
	// }

	region := &qh_domain.QHRegion{
		LayerID:            req.LayerID,
		LabelID:            &labelID,
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
		// ColorRed:           colorRed,
		// ColorGreen:         colorGreen,
		// ColorBlue:          colorBlue,
	}

	if name, ok := props["name"]; ok {
		region.Name = fmt.Sprintf("%v", name)
	}
	if dn, ok := props["display_name"]; ok {
		region.DisplayName = fmt.Sprintf("%v", dn)
	}
	if tk, ok := props["ten_khu"]; ok {
		if region.Name == "" {
			region.Name = fmt.Sprintf("%v", tk)
		}
		if region.DisplayName == "" {
			region.DisplayName = fmt.Sprintf("%v", tk)
		}
	}

	return region, nil
}

// buildRegionExtend builds a QHRegionExtend for non-polygon geometry types
// (Point, LineString, MultiLineString, etc.).
func (u *importUsecaseImpl) buildRegionExtend(ctx context.Context, batchID string, index int, feat geoJSONFeature, req *dto.ImportRequest, lc *labelCache) (*qh_domain.QHRegionExtend, *dto.ImportFailure) {
	if feat.Geometry == nil {
		if req.SkipInvalid {
			return nil, nil
		}
		return nil, &dto.ImportFailure{FeatureIndex: index, Reason: "empty geometry"}
	}

	geomBytes, err := geometry.ToRawGeoJSONBytes(feat.Geometry)
	if err != nil {
		if req.SkipInvalid {
			return nil, nil
		}
		return nil, &dto.ImportFailure{FeatureIndex: index, Reason: fmt.Sprintf("geometry serialize: %v", err)}
	}

	geomType := geometry.GetGeomType(feat.Geometry)

	props := feat.Properties
	if props == nil {
		props = map[string]interface{}{}
	}

	labelID, err := u.resolveLabelIDCached(ctx, props, req, lc)
	if err != nil {
		if req.SkipInvalid {
			return nil, nil
		}
		return nil, &dto.ImportFailure{FeatureIndex: index, Reason: fmt.Sprintf("label resolution: %v", err)}
	}

	propsBytes, _ := json.Marshal(props)

	record := &qh_domain.QHRegionExtend{
		LayerID:            req.LayerID,
		LabelID:            &labelID,
		GeomType:           geomType,
		Geometry:           qh_domain.RawExtendGeometry{Raw: geomBytes},
		OriginalProperties: propsBytes,
		SourceFile:         req.SourceFileName,
		ImportBatchID:      batchID,
		Status:             10,
		Version:            1,
		IsLatest:           true,
	}

	if name, ok := props["name"]; ok {
		record.Name = fmt.Sprintf("%v", name)
	}
	if dn, ok := props["display_name"]; ok {
		record.DisplayName = fmt.Sprintf("%v", dn)
	}
	if tk, ok := props["ten_khu"]; ok {
		if record.Name == "" {
			record.Name = fmt.Sprintf("%v", tk)
		}
		if record.DisplayName == "" {
			record.DisplayName = fmt.Sprintf("%v", tk)
		}
	}

	return record, nil
}

// ─── Label resolution (cache-aware) ───────────────────────────────────────

// resolveLabelIDCached wraps resolveLabelID with an in-memory cache so that
// repeated labels (very common in large files) hit the cache instead of DB.
func (u *importUsecaseImpl) resolveLabelIDCached(ctx context.Context, properties map[string]interface{}, req *dto.ImportRequest, lc *labelCache) (uint64, error) {
	// Fast path: explicit mapping
	if req.LabelField != "" {
		rawValue, ok := properties[req.LabelField]
		if ok {
			value := normalizeValue(rawValue)
			if labelID, ok := req.LabelMappings[value]; ok {
				return labelID, nil
			}
			return u.labelByNameCached(ctx, value, req.LayerID, lc)
		}
	}

	// Scan properties for a likely label value
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
		return u.labelByNameCached(ctx, value, req.LayerID, lc)
	}

	if req.DefaultLabelID != nil {
		return *req.DefaultLabelID, nil
	}
	return u.labelByNameCached(ctx, "chưa phân loại", req.LayerID, lc)
}

// labelByNameCached resolves nhãn theo tên (toàn cục), gắn nhãn vào layer import qua qh_label_layers,
// tạo mới nếu chưa có. Cache chỉ theo tên đã chuẩn hóa.
func (u *importUsecaseImpl) labelByNameCached(ctx context.Context, rawName string, importLayerID uint64, lc *labelCache) (uint64, error) {
	normalized := normalizeLabelName(rawName)
	if id, ok := lc.get(normalized); ok {
		return id, nil
	}

	label, err := u.labelRepo.GetByName(ctx, normalized)
	if err != nil {
		return 0, err
	}
	if label != nil {
		if err := u.labelRepo.EnsureLayerLink(ctx, label.ID, importLayerID); err != nil {
			return 0, fmt.Errorf("ensure label-layer link: %w", err)
		}
		lc.set(normalized, label.ID)
		return label.ID, nil
	}

	id, err := u.createNewLabel(ctx, normalized, importLayerID)
	if err != nil {
		return 0, err
	}
	lc.set(normalized, id)
	return id, nil
}

func (u *importUsecaseImpl) createNewLabel(ctx context.Context, name string, layerID uint64) (uint64, error) {
	normalized := normalizeLabelName(name)
	newLabel := &qh_domain.QHLabel{
		LayerID:     layerID,
		Name:        normalized,
		DisplayName: normalized,
		Color:       generateColorFromString(normalized),
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

// ─── Batch status & rollback (unchanged logic) ────────────────────────────

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
	var errMsgs []string
	if layer != nil {
		switch layer.ImportStatus {
		case enums.LayerImportStatusProcessing:
			statusStr = "processing"
		case enums.LayerImportStatusFailed:
			statusStr = "failed"
			errMsgs = []string{"import job failed"}
		}
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

// ListImportRegionErrors — danh sách bản ghi qh_region_import_error_logs theo layer (phân trang SQL).
func (u *importUsecaseImpl) ListImportRegionErrors(ctx context.Context, layerID uint64, pagable *_dto.Pagable) ([]qh_domain.QHRegionImportErrorLog, int64, error) {
	if layerID == 0 {
		return nil, 0, importLayerIDRequired()
	}
	layer, err := u.layerRepo.GetByID(ctx, layerID)
	if err != nil {
		return nil, 0, err
	}
	if layer == nil {
		return nil, 0, importLayerNotFound(layerID)
	}
	if pagable == nil {
		pagable = _dto.NewPagableFromGrpc(nil, nil, nil)
	}
	offset := pagable.GetOffset()
	limit := pagable.GetLimit()
	return u.regionRepo.ListImportErrorLogsByLayerID(ctx, layerID, offset, limit)
}

// RetryImportError — đọc log lỗi theo id, CAS pending→processing, insert từng region; cuối cùng resolved hoặc pending.
func (u *importUsecaseImpl) RetryImportError(ctx context.Context, errorID uint64) (*dto.RetryImportErrorResult, error) {
	if errorID == 0 {
		return nil, importErrorIDRequired()
	}
	row, err := u.regionRepo.GetImportErrorLogByID(ctx, errorID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, importErrorNotFound(errorID)
	}
	if row.Status == qh_domain.ImportErrorStatusProcessing {
		return nil, importRetryInProgress(errorID)
	}
	if row.Status != qh_domain.ImportErrorStatusPending {
		return nil, importRetryNotAllowed(errorID, row.Status)
	}

	locked, err := u.regionRepo.TryBeginImportErrorRetry(ctx, errorID)
	if err != nil {
		return nil, err
	}
	if !locked {
		return nil, importRetryLockConflict(errorID)
	}

	retryCommitted := false
	defer func() {
		if !retryCommitted {
			if rerr := u.regionRepo.ResetImportErrorRetryIfProcessing(ctx, errorID); rerr != nil {
				log.Printf("[retry-import] reset processing→pending errorLog=%d: %v", errorID, rerr)
			}
		}
	}()

	row, err = u.regionRepo.GetImportErrorLogByID(ctx, errorID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, importErrorNotFound(errorID)
	}

	var snapshots []qh_domain.RegionSnapshot
	if err := json.Unmarshal(row.RegionsPayload, &snapshots); err != nil {
		return nil, fmt.Errorf("decode regions payload: %w", err)
	}
	n := len(snapshots)
	if n == 0 {
		return nil, fmt.Errorf("no regions in payload")
	}
	layer, err := u.layerRepo.GetByID(ctx, row.LayerID)
	if err != nil {
		return nil, err
	}
	if layer == nil {
		return nil, importLayerNotFound(row.LayerID)
	}
	lc := newLabelCache()
	retryBatchID := generateBatchID()
	okCount, failCount := 0, 0
	remaining := make([]qh_domain.RegionSnapshot, 0)
	for i := range snapshots {
		snap := snapshots[i]
		snap.LayerID = row.LayerID
		reg, rerr := u.regionFromRetrySnapshot(ctx, &snap, retryBatchID, lc)
		if rerr != nil {
			failCount++
			remaining = append(remaining, snap)
			log.Printf("[retry-import] errorLog=%d item=%d build: %v", errorID, i, rerr)
			continue
		}
		if err := u.regionRepo.Create(ctx, reg); err != nil {
			failCount++
			remaining = append(remaining, snap)
			log.Printf("[retry-import] errorLog=%d item=%d insert: %v", errorID, i, err)
			continue
		}
		okCount++
	}
	newPayload, mErr := json.Marshal(remaining)
	if mErr != nil {
		return nil, fmt.Errorf("marshal remaining regions payload: %w", mErr)
	}
	row.RegionsPayload = newPayload
	row.BatchSize = len(remaining)
	row.RetryCount++
	if len(remaining) == 0 {
		row.Status = qh_domain.ImportErrorStatusResolved
	} else {
		row.Status = qh_domain.ImportErrorStatusPending
	}
	if err := u.regionRepo.UpdateImportErrorLog(ctx, row); err != nil {
		return nil, fmt.Errorf("update import error log: %w", err)
	}
	retryCommitted = true

	msg := fmt.Sprintf("retry done: success=%d failed=%d total=%d", okCount, failCount, n)
	return &dto.RetryImportErrorResult{
		Success:      failCount == 0,
		Message:      msg,
		TotalItems:   n,
		SuccessCount: okCount,
		FailedCount:  failCount,
	}, nil
}

func (u *importUsecaseImpl) regionFromRetrySnapshot(ctx context.Context, snap *qh_domain.RegionSnapshot, retryBatchID string, lc *labelCache) (*qh_domain.QHRegion, error) {
	if snap == nil {
		return nil, fmt.Errorf("nil snapshot")
	}
	if strings.TrimSpace(snap.Geometry) == "" {
		return nil, fmt.Errorf("empty geometry")
	}
	var geomObj map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(snap.Geometry)), &geomObj); err != nil {
		return nil, fmt.Errorf("geometry json: %w", err)
	}
	geomBytes, err := geometry.ToMultiPolygonGeoJSONBytes(geomObj)
	if err != nil {
		return nil, err
	}
	var props map[string]interface{}
	if len(snap.OriginalProperties) > 0 {
		if err := json.Unmarshal(snap.OriginalProperties, &props); err != nil {
			return nil, fmt.Errorf("originalProperties: %w", err)
		}
	} else {
		props = map[string]interface{}{}
	}
	req := &dto.ImportRequest{
		LayerID:          snap.LayerID,
		SourceFileName:   snap.SourceFile,
		ValidateGeometry: true,
		AutoFixGeometry:  false,
		SkipInvalid:      false,
	}
	labelID, err := u.resolveLabelIDCached(ctx, props, req, lc)
	if err != nil {
		return nil, err
	}
	geoInfo, err := geometry.ParseGeometry(string(geomBytes))
	if err != nil {
		return nil, fmt.Errorf("invalid geometry: %w", err)
	}
	area := geoInfo.Area
	propsBytes, _ := json.Marshal(props)

	// colorRed, colorGreen, colorBlue := uint8(200), uint8(200), uint8(200)
	// if v, ok := props["red"]; ok {
	// 	if f, ok := v.(float64); ok {
	// 		colorRed = uint8(f)
	// 	}
	// }
	// if v, ok := props["green"]; ok {
	// 	if f, ok := v.(float64); ok {
	// 		colorGreen = uint8(f)
	// 	}
	// }
	// if v, ok := props["blue"]; ok {
	// 	if f, ok := v.(float64); ok {
	// 		colorBlue = uint8(f)
	// 	}
	// }

	region := &qh_domain.QHRegion{
		LayerID:            snap.LayerID,
		LabelID:            &labelID,
		Geometry:           qh_domain.RawGeometry{Raw: geomBytes},
		OriginalProperties: propsBytes,
		AreaSqm:            area,
		AreaHa:             area / 10000,
		SourceFile:         snap.SourceFile,
		ImportBatchID:      retryBatchID,
		Status:             10,
		ProcessingStatus:   10,
		Version:            1,
		IsLatest:           true,
		// ColorRed:           colorRed,
		// ColorGreen:         colorGreen,
		// ColorBlue:          colorBlue,
	}
	region.Name = snap.Name
	region.DisplayName = snap.DisplayName
	if name, ok := props["name"]; ok && region.Name == "" {
		region.Name = fmt.Sprintf("%v", name)
	}
	if dn, ok := props["display_name"]; ok && region.DisplayName == "" {
		region.DisplayName = fmt.Sprintf("%v", dn)
	}
	if tk, ok := props["ten_khu"]; ok {
		if region.Name == "" {
			region.Name = fmt.Sprintf("%v", tk)
		}
		if region.DisplayName == "" {
			region.DisplayName = fmt.Sprintf("%v", tk)
		}
	}
	return region, nil
}

// ─── Helpers (unchanged) ──────────────────────────────────────────────────

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
	return strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", v)))
}

// normalizeLabelName returns a consistent lowercase+trimmed label name.
// Every label name stored in DB and used as cache keys goes through this — guarantees
// "Đất Rừng", "đất rừng", "  ĐẤT RỪNG  " all resolve to the same identity.
func normalizeLabelName(rawName string) string {
	s := strings.ToLower(strings.TrimSpace(rawName))
	if s == "" {
		return "chưa phân loại"
	}
	return s
}

func isCoordinateOrNumeric(value string) bool {
	if _, err := fmt.Sscanf(value, "%f", new(float64)); err == nil && len(value) < 20 {
		return true
	}
	if strings.ContainsAny(value, ", ") {
		parts := strings.FieldsFunc(value, func(r rune) bool { return r == ',' || r == ' ' })
		n := 0
		for _, p := range parts {
			if _, err := fmt.Sscanf(p, "%f", new(float64)); err == nil {
				n++
			}
		}
		if n >= 2 && n <= 4 {
			return true
		}
	}
	return false
}
