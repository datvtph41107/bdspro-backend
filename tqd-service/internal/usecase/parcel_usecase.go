package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"sync"
	"time"

	_dto "common/domain/dto"
	_utils "common/utils"
	"tqd/internal/config"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	qh_dto "tqd/internal/dto/qh"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"
	"tqd/internal/usecase/resolver/layer_resolver/engine/parcel_engine"
	"tqd/internal/usecase/resolver/layer_resolver/rule"
	"tqd/internal/usecase/resolver/layer_resolver/types"
)

type ParcelUsecase interface {
	GetParcelInfoByID(ctx context.Context, parcelID uint64) (*qh_domain.ParcelInfoResponse, error)
	SearchParcelsByText(ctx context.Context, text string, mapNumber, landNumber string, page, limit int) ([]*qh_domain.ParcelInfoResponse, int64, error)
	SearchPublicParcels(ctx context.Context, page, limit int) ([]*qh_domain.ParcelInfoResponse, int64, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]qh_domain.Parcel, error)
	CreateBatch(ctx context.Context, req *dto.ParcelList) (*dto.ParcelList, error)
	FindByLocation(ctx context.Context, lat, lng float64) (*qh_domain.Parcel, error)
	ResolveMapTargetByLocation(ctx context.Context, req dto.ResolveMapTargetRequestDTO) (*dto.ResolveMapTargetResponseDTO, error)
	FindParcelsByPolygon(ctx context.Context, points []repo.Point, intersect bool, maxAreaKm2 float64, page, pageSize int) ([]qh_domain.Parcel, int64, error)
	FindMapTargetsByPolygon(ctx context.Context, points []repo.Point, intersect bool, maxAreaKm2 float64, page, pageSize int, zoom *uint32, filter *dto.PolygonTargetFilter, includeAnalysis bool) (*dto.PolygonMapTargetsResponse, error)
	GetParcelInfo(ctx context.Context, parcelID uint64) (*dto.ParcelInfoResponse, error)
	GetParcelPlanning(ctx context.Context, req *dto.GetParcelPlanningRequest) (*ParcelPlanningResult, error)
	GetParcelTimeline(ctx context.Context, parcelID uint64) (*types.ParcelTimeline, error)
	GetParcelQuickInfo(ctx context.Context, parcelID uint64) ([]dto.ParcelLayerInfo, error)
	GetLayerHistory(ctx context.Context, layerID uint64) (*types.LayerLineage, error)
	GetParcelSeoSource(ctx context.Context, parcelID uint64) (*dto.ParcelSeoSource, error)
	ListParcelSeoSourcesForGenerate(ctx context.Context, limit uint32) ([]dto.ParcelSeoSource, error)
	UpdateParcelSeoID(ctx context.Context, parcelID uint64, seoID uint64) error

	GetParcelDetail(ctx context.Context, parcelID uint64) (*qh_dto.ParcelDetailResponseDTO, error)
	GetParcelLayers(ctx context.Context, parcelID uint64, pagable _dto.Pagable) ([]*dto.LayerImpact, int64, error)
	GetZoneGeometry(ctx context.Context, parcelID, layerID, zoneID uint64) (string, *ZoneGeometryProps, error)
	GetParcelLegalDocuments(ctx context.Context, parcelID uint64, docType, status *uint32, pagable _dto.Pagable) ([]*qh_domain.QHLegalDocument, int64, error)
	CheckingPolygon(ctx context.Context, points []repo.Point, intersect bool, maxAreaKm2 float64, page, pageSize int, zoom *uint32, filter *dto.PolygonTargetFilter, includeAnalysis bool) (*dto.PolygonMapTargetsResponse, error)
}

// ZoneGeometryProps — Properties cho Zone Geometry response
type ZoneGeometryProps struct {
	ZoneName         string
	LandUseCode      string
	LandUseName      string
	LandUseColor     string
	BuildStatus      enums.BuildStatus
	AlertLevel       enums.AlertLevel
	LegalDocumentIds []string
}

type parcelUsecase struct {
	parcelRepo      repo.IParcelRepo
	layerLegalRepo  repo.ILayerLegalRepo
	legalDocRepo    repo.ILegalDocumentRepo
	auditRepo       rule.AuditRepository
	layerRepo       rule.LayerRepository
	parcelEngine    *parcel_engine.ParcelEngine
	timelineEngine  *rule.TimelineEngine
	lifecycleEngine *rule.LifecycleEngine
}

func NewParcelUsecase(
	parcelRepo repo.IParcelRepo,
	layerLegalRepo repo.ILayerLegalRepo,
	legalDocRepo repo.ILegalDocumentRepo,
	auditRepo rule.AuditRepository,
	layerRepo rule.LayerRepository,
) ParcelUsecase {
	auditEngine := rule.NewAuditEngine(nil, auditRepo)
	lifecycleEngine := rule.NewLifecycleEngine(nil, layerRepo)
	timelineEngine := rule.NewTimelineEngine(nil, auditRepo, nil)

	engine := parcel_engine.NewParcelEngine(parcel_engine.DefaultEngineConfig()).
		WithAuditEngine(auditEngine).
		WithLifecycleEngine(lifecycleEngine).
		WithTimelineEngine(timelineEngine)

	return &parcelUsecase{
		parcelRepo:      parcelRepo,
		layerLegalRepo:  layerLegalRepo,
		legalDocRepo:    legalDocRepo,
		auditRepo:       auditRepo,
		layerRepo:       layerRepo,
		parcelEngine:    engine,
		timelineEngine:  timelineEngine,
		lifecycleEngine: lifecycleEngine,
	}
}

// ============================================================
// EXISTING METHODS (giữ nguyên, không thay đổi)
// ============================================================

func (u *parcelUsecase) GetParcelInfoByID(ctx context.Context, parcelID uint64) (*qh_domain.ParcelInfoResponse, error) {
	if parcelID == 0 {
		return nil, errors.New("parcel_id is required")
	}

	info, err := u.parcelRepo.GetParcelPreview(ctx, parcelID)
	if err != nil {
		return nil, fmt.Errorf("get parcel info failed: %w", err)
	}
	if info == nil || info.ID == 0 {
		return nil, fmt.Errorf("parcel %d not found", parcelID)
	}

	mapNum, _ := strconv.ParseUint(info.MapNumber, 10, 32)
	landNum, _ := strconv.ParseUint(info.LandNumber, 10, 32)

	return &qh_domain.ParcelInfoResponse{
		ID:         info.ID,
		AddressStr: info.AddressText,
		MapNumber:  uint32(mapNum),
		LandNumber: uint32(landNum),
		Lat:        info.Latitude,
		Lng:        info.Longitude,
		Directions: info.Directions,
		Shape:      info.Shape,
		ShapeId:    info.ShapeId,
		Facade:     uint32(info.Facade),
		TotalArea:  info.TotalAreaSqm,
		Geometry:   info.Geometry,
	}, nil
}

func (u *parcelUsecase) SearchParcelsByText(ctx context.Context, text string, mapNumber, landNumber string, page, limit int) ([]*qh_domain.ParcelInfoResponse, int64, error) {
	if text == "" && mapNumber == "" && landNumber == "" {
		return nil, 0, errors.New("search criteria is required")
	}

	pagable := _dto.Pagable{Page: uint32(page), Size: uint32(limit)}
	pagable.Normalize()

	parcels, total, err := u.parcelRepo.SearchByAddress(ctx, text, mapNumber, landNumber, pagable)
	if err != nil {
		return nil, 0, fmt.Errorf("search parcels failed: %w", err)
	}

	results := make([]*qh_domain.ParcelInfoResponse, 0, len(parcels))
	for _, p := range parcels {
		mapNum, _ := strconv.ParseUint(p.MapNumber, 10, 32)
		landNum, _ := strconv.ParseUint(p.LandNumber, 10, 32)

		results = append(results, &qh_domain.ParcelInfoResponse{
			ID:         p.ID,
			ParcelID:   p.ParcelID,
			AddressStr: p.AddressText,
			MapNumber:  uint32(mapNum),
			LandNumber: uint32(landNum),
			Lat:        p.Latitude,
			Lng:        p.Longitude,
			Shape:      p.Shape,
			ShapeId:    p.ShapeId,
			Facade:     uint32(p.Facade),
			TotalArea:  p.TotalAreaSqm,
			Geometry:   p.Geometry,
		})
	}

	return results, total, nil
}

func (u *parcelUsecase) SearchPublicParcels(ctx context.Context, page, limit int) ([]*qh_domain.ParcelInfoResponse, int64, error) {
	pagable := _dto.Pagable{Page: uint32(page), Size: uint32(limit)}
	pagable.Normalize()

	parcels, total, err := u.parcelRepo.SearchPublicParcels(ctx, pagable)
	if err != nil {
		return nil, 0, fmt.Errorf("search public parcels failed: %w", err)
	}

	results := make([]*qh_domain.ParcelInfoResponse, 0, len(parcels))
	for _, p := range parcels {
		mapNum, _ := strconv.ParseUint(p.MapNumber, 10, 32)
		landNum, _ := strconv.ParseUint(p.LandNumber, 10, 32)

		results = append(results, &qh_domain.ParcelInfoResponse{
			ID:         p.ID,
			ParcelID:   p.ParcelID,
			AddressStr: p.AddressText,
			MapNumber:  uint32(mapNum),
			LandNumber: uint32(landNum),
			Lat:        p.Latitude,
			Lng:        p.Longitude,
			Shape:      p.Shape,
			ShapeId:    p.ShapeId,
			Facade:     uint32(p.Facade),
			TotalArea:  p.TotalAreaSqm,
			Geometry:   p.Geometry,
		})
	}

	return results, total, nil
}

func (u *parcelUsecase) GetByIDs(ctx context.Context, ids []uint64) ([]qh_domain.Parcel, error) {
	if len(ids) == 0 {
		return []qh_domain.Parcel{}, nil
	}
	return u.parcelRepo.GetByIDs(ctx, ids)
}

func (u *parcelUsecase) CreateBatch(ctx context.Context, req *dto.ParcelList) (*dto.ParcelList, error) {
	for _, parcel := range req.Data {
		if parcel.Lat < -90 || parcel.Lat > 90 {
			return nil, fmt.Errorf("latitude must be between -90 and 90")
		}
		if parcel.Lng < -180 || parcel.Lng > 180 {
			return nil, fmt.Errorf("longitude must be between -180 and 180")
		}
	}

	if err := u.parcelRepo.CreateBatch(ctx, req.Data); err != nil {
		return nil, fmt.Errorf("failed to create parcels: %v", err)
	}
	return req, nil
}

func (u *parcelUsecase) FindByLocation(ctx context.Context, lat, lng float64) (*qh_domain.Parcel, error) {
	parcel, err := u.parcelRepo.FindByLocation(ctx, lat, lng)
	if err != nil {
		return nil, fmt.Errorf("failed to find parcel: %v", err)
	}
	if parcel == nil {
		return nil, fmt.Errorf("parcel not found")
	}
	return parcel, nil
}

// parcelLocationOverviewReader is an optional optimized repository capability
// used by the pointer-selection hot path. It resolves point-in-polygon and the
// Quick Overview projection in one query, without serializing full geometry.
// The public ParcelRepo contract remains unchanged so other implementations
// stay backward compatible during the migration.
type parcelLocationOverviewReader interface {
	ResolveParcelOverviewByLocation(ctx context.Context, lat, lng float64) (*dto.ParcelInfoResponse, error)
}

// parcelOverviewReader is the compatibility path for repositories that have
// not implemented the one-query location overview yet.
type parcelOverviewReader interface {
	GetParcelOverview(ctx context.Context, parcelID uint64) (*dto.ParcelInfoResponse, error)
}

func (u *parcelUsecase) ResolveMapTargetByLocation(
	ctx context.Context,
	req dto.ResolveMapTargetRequestDTO,
) (*dto.ResolveMapTargetResponseDTO, error) {
	if req.Latitude < -90 || req.Latitude > 90 {
		return nil, fmt.Errorf("latitude must be between -90 and 90")
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		return nil, fmt.Errorf("longitude must be between -180 and 180")
	}

	// Optimized repositories resolve containment + overview in one lightweight
	// query. This is the normal path for interactive map clicks.
	if locationOverviewRepo, ok := u.parcelRepo.(parcelLocationOverviewReader); ok {
		info, err := locationOverviewRepo.ResolveParcelOverviewByLocation(
			ctx,
			req.Latitude,
			req.Longitude,
		)
		if err != nil {
			return nil, fmt.Errorf("resolve parcel overview by location failed: %w", err)
		}
		if info != nil && info.ParcelID > 0 {
			return &dto.ResolveMapTargetResponseDTO{
				TargetType: enums.WorkspaceEntityTypeParcel.Uint32(),
				TargetID:   info.ParcelID,
				Parcel:     info,
			}, nil
		}
	} else {
		// Compatibility path. It may perform two queries and include geometry,
		// but preserves behavior for alternative repository implementations.
		parcel, err := u.parcelRepo.FindByLocation(ctx, req.Latitude, req.Longitude)
		if err != nil {
			return nil, fmt.Errorf("find parcel by location failed: %w", err)
		}

		if parcel != nil && parcel.ID > 0 {
			var info *dto.ParcelInfoResponse
			if overviewRepo, ok := u.parcelRepo.(parcelOverviewReader); ok {
				info, err = overviewRepo.GetParcelOverview(ctx, parcel.ID)
			} else {
				info, err = u.parcelRepo.GetParcelInfo(ctx, parcel.ID)
			}
			if err != nil {
				return nil, fmt.Errorf("get parcel overview failed: %w", err)
			}
			if info == nil || info.ParcelID == 0 {
				return nil, fmt.Errorf("parcel info not found")
			}
			if len(info.Geometry) == 0 && len(parcel.Geometry.Raw) > 0 {
				info.Geometry = parcel.Geometry.Raw
			}

			return &dto.ResolveMapTargetResponseDTO{
				TargetType: enums.WorkspaceEntityTypeParcel.Uint32(),
				TargetID:   info.ParcelID,
				Parcel:     info,
			}, nil
		}
	}

	region, err := u.parcelRepo.FindRegionByLocation(ctx, req.Latitude, req.Longitude)
	if err != nil {
		return nil, fmt.Errorf("find region by location failed: %w", err)
	}

	if region != nil && region.RegionID > 0 {
		return &dto.ResolveMapTargetResponseDTO{
			TargetType: enums.WorkspaceEntityTypeRegion.Uint32(),
			TargetID:   region.RegionID,
			Region:     region,
		}, nil
	}

	return nil, fmt.Errorf("map target not found")
}

func (u *parcelUsecase) FindParcelsByPolygon(
	ctx context.Context,
	points []repo.Point,
	intersect bool,
	maxAreaKm2 float64,
	page, pageSize int,
) ([]qh_domain.Parcel, int64, error) {
	limit, offset, err := validatePolygonQuery(points, maxAreaKm2, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	if intersect {
		return u.parcelRepo.FindIntersectPolygon(ctx, points, limit, offset, nil)
	}

	return u.parcelRepo.FindWithinPolygon(ctx, points, limit, offset, nil)
}

func (u *parcelUsecase) FindMapTargetsByPolygon(
	ctx context.Context,
	points []repo.Point,
	intersect bool,
	maxAreaKm2 float64,
	page, pageSize int,
	zoom *uint32,
	filter *dto.PolygonTargetFilter,
	includeAnalysis bool,
) (*dto.PolygonMapTargetsResponse, error) {
	limit, offset, err := validatePolygonQuery(points, maxAreaKm2, page, pageSize)
	if err != nil {
		return nil, err
	}

	normalizedPage, normalizedPageSize := normalizePageValues(page, pageSize)
	effectiveZoom := uint32(config.ParcelPolygonMinZ)
	if zoom != nil {
		effectiveZoom = *zoom
	}

	if filter != nil {
		filter.Zoom = &effectiveZoom
	}

	if zoom != nil && *zoom < config.ParcelPolygonMinZ {
		regionZoom := *zoom
		var (
			regions     []dto.RegionInfoResponse
			total       int64
			analysis    *dto.PolygonAnalysisSummary
			parcelCount int64
			wg          sync.WaitGroup
			mu          sync.Mutex
			queryErr    error
		)

		recordErr := func(err error) {
			if err == nil {
				return
			}
			mu.Lock()
			if queryErr == nil {
				queryErr = err
			}
			mu.Unlock()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			r, t, err := u.parcelRepo.FindRegionsByPolygon(ctx, points, intersect, regionZoom, limit, offset, filter, false)
			if err != nil {
				recordErr(err)
				return
			}
			regions = r
			total = t
		}()

		if includeAnalysis {
			wg.Add(2)
			go func() {
				defer wg.Done()
				a, _, err := u.parcelRepo.SummarizePolygonAnalysis(ctx, points, intersect, regionZoom, filter, false)
				if err != nil {
					recordErr(err)
					return
				}
				analysis = a
			}()
			go func() {
				defer wg.Done()
				count, err := u.parcelRepo.CountParcelsByPolygon(
					ctx,
					points,
					intersect,
					u.polygonFilterWithZoom(filter, regionZoom),
				)
				if err != nil {
					slog.ErrorContext(ctx, fmt.Sprintf("[FindMapTargetsByPolygon] CountParcelsByPolygon error: %v", err))
					return
				}
				parcelCount = count
			}()
		}

		wg.Wait()
		if queryErr != nil {
			return nil, queryErr
		}

		resp := &dto.PolygonMapTargetsResponse{
			ResultType: enums.WorkspaceEntityTypeRegion.Uint32(),
			Parcels:    []qh_domain.Parcel{},
			Regions:    regions,
			Total:      total,
			Z:          regionZoom,
			Page:       uint32(normalizedPage),
			PageSize:   uint32(normalizedPageSize),
		}
		if analysis != nil {
			analysis.ResultType = resp.ResultType
			analysis.Z = regionZoom
			analysis.ParcelCount = parcelCount
			resp.Analysis = analysis
		}
		return resp, nil
	}

	var (
		parcels     []qh_domain.Parcel
		total       int64
		analysis    *dto.PolygonAnalysisSummary
		parcelCount int64
		wg          sync.WaitGroup
		mu          sync.Mutex
		queryErr    error
	)

	recordErr := func(err error) {
		if err == nil {
			return
		}
		mu.Lock()
		if queryErr == nil {
			queryErr = err
		}
		mu.Unlock()
	}

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	var p []qh_domain.Parcel
	// 	var t int64
	// 	var err error
	// 	if intersect {
	// 		p, t, err = u.parcelRepo.FindIntersectPolygon(ctx, points, limit, offset, filter)
	// 	} else {
	// 		p, t, err = u.parcelRepo.FindWithinPolygon(ctx, points, limit, offset, filter)
	// 	}
	// 	if err != nil {
	// 		recordErr(err)
	// 		return
	// 	}
	// 	parcels = p
	// 	total = t
	// }()

	if includeAnalysis {
		// Two request-scoped queries run in parallel and must both be joined
		// before this method returns. Keep WaitGroup accounting equal to the
		// number of goroutines; otherwise one branch can outlive the response
		// and the late Done panics with a negative WaitGroup counter.
		wg.Add(2)
		go func() {
			defer wg.Done()
			a, _, err := u.parcelRepo.SummarizePolygonAnalysis(ctx, points, intersect, effectiveZoom, filter, false)
			if err != nil {
				recordErr(err)
				return
			}
			analysis = a
		}()
		go func() {
			defer wg.Done()
			count, err := u.parcelRepo.CountParcelsByPolygon(
				ctx,
				points,
				intersect,
				u.polygonFilterWithZoom(filter, effectiveZoom),
			)
			if err != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("[FindMapTargetsByPolygon] CountParcelsByPolygon error: %v", err))
				return
			}
			parcelCount = count
		}()
	}

	wg.Wait()
	if queryErr != nil {
		return nil, queryErr
	}

	resp := &dto.PolygonMapTargetsResponse{
		ResultType: enums.WorkspaceEntityTypeParcel.Uint32(),
		Parcels:    parcels,
		Regions:    []dto.RegionInfoResponse{},
		Total:      total,
		Page:       uint32(normalizedPage),
		PageSize:   uint32(normalizedPageSize),
	}

	if zoom != nil {
		resp.Z = *zoom
	}

	if analysis != nil {
		analysis.ResultType = resp.ResultType
		analysis.Z = effectiveZoom
		analysis.ParcelCount = parcelCount
		resp.Analysis = analysis
	}

	return resp, nil
}

func (u *parcelUsecase) CheckingPolygon(
	ctx context.Context,
	points []repo.Point,
	intersect bool,
	maxAreaKm2 float64,
	page, pageSize int,
	zoom *uint32,
	filter *dto.PolygonTargetFilter,
	includeAnalysis bool,
) (*dto.PolygonMapTargetsResponse, error) {
	_, _, err := validatePolygonQuery(points, maxAreaKm2, page, pageSize)
	var limit int = 500
	var offset int = 0
	if err != nil {
		return nil, err
	}

	normalizedPage, normalizedPageSize := normalizePageValues(page, pageSize)
	area := calculatePolygonAreaKm2(points)

	effectiveZoom := uint32(14)
	if zoom != nil && *zoom < config.ParcelPolygonMinZ {
		effectiveZoom = *zoom
	}

	startedAt := time.Now()
	var (
		regions        []dto.RegionInfoResponse
		analysis       *dto.PolygonAnalysisSummary
		parcelCount    int64
		quickLayerRows []dto.ParcelLayerInfo
		hasQuickLayers bool
		wg             sync.WaitGroup
		mu             sync.Mutex
		findRegionsMs  time.Duration
		quickInfoMs    time.Duration
		countParcelsMs time.Duration
		analysisMs     time.Duration
	)

	recordDuration := func(target *time.Duration, start time.Time) {
		mu.Lock()
		*target = time.Since(start)
		mu.Unlock()
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		start := time.Now()
		r, _, err := u.parcelRepo.FindRegionsByPolygon(ctx, points, intersect, effectiveZoom, limit, offset, filter, true)
		recordDuration(&findRegionsMs, start)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[CheckingPolygon] FindRegions error: %v", err))
			return
		}
		regions = r
	}()
	go func() {
		defer wg.Done()
		start := time.Now()
		count, err := u.parcelRepo.CountParcelsByPolygon(
			ctx,
			points,
			intersect,
			u.polygonFilterWithZoom(filter, effectiveZoom),
		)
		recordDuration(&countParcelsMs, start)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[CheckingPolygon] CountParcelsByPolygon error: %v", err))
			return
		}
		parcelCount = count
	}()

	if includeAnalysis {
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now()
			a, rows, err := u.parcelRepo.SummarizePolygonAnalysis(ctx, points, intersect, effectiveZoom, filter, true)
			recordDuration(&analysisMs, start)
			if err != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("[CheckingPolygon] Analysis error: %v", err))
				return
			}
			analysis = a
			if rows == nil {
				rows = []dto.ParcelLayerInfo{}
			}
			quickLayerRows = rows
			hasQuickLayers = true
		}()
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now()
			rows, err := u.parcelRepo.GetPolygonQuickInfo(ctx, points, intersect, effectiveZoom, filter)
			recordDuration(&quickInfoMs, start)
			if err != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("[CheckingPolygon] GetPolygonQuickInfo error: %v", err))
				return
			}
			if rows == nil {
				rows = []dto.ParcelLayerInfo{}
			}
			quickLayerRows = rows
			hasQuickLayers = true
		}()
	}

	wg.Wait()
	totalMs := time.Since(startedAt)

	resp := &dto.PolygonMapTargetsResponse{
		ResultType: 0,
		Parcels:    []qh_domain.Parcel{},
		Regions:    regions,
		Total:      parcelCount,
		Page:       uint32(normalizedPage),
		PageSize:   uint32(normalizedPageSize),
		Analysis:   analysis,
	}
	if analysis != nil {
		analysis.ParcelCount = parcelCount
	}
	if hasQuickLayers {
		resp.QuickLayerRows = quickLayerRows
	}

	if zoom != nil {
		resp.Z = *zoom
	} else {
		resp.Z = effectiveZoom
	}

	if area > 1.0 {
		resp.Warning = fmt.Sprintf("Diện tích vùng chọn lớn (%.2f ha). Kết quả thửa đất có thể không đầy đủ.", area*100)
	}
	slog.InfoContext(ctx, fmt.Sprintf("[CheckingPolygon] timing findRegions=%v quickInfo=%v countParcels=%v analysis=%v total=%v | regions=%d parcelCount=%d quickLayers=%v includeAnalysis=%v",
		findRegionsMs,
		quickInfoMs,
		countParcelsMs,
		analysisMs,
		totalMs,
		len(regions),
		parcelCount,
		hasQuickLayers,
		includeAnalysis),
	)

	return resp, nil
}

func (u *parcelUsecase) polygonFilterWithZoom(filter *dto.PolygonTargetFilter, zoom uint32) *dto.PolygonTargetFilter {
	localFilter := filter
	if localFilter == nil {
		localFilter = &dto.PolygonTargetFilter{}
	} else {
		copyFilter := *localFilter
		localFilter = &copyFilter
	}
	localFilter.Zoom = &zoom
	return localFilter
}

func (u *parcelUsecase) attachPolygonParcelCount(
	ctx context.Context,
	analysis *dto.PolygonAnalysisSummary,
	points []repo.Point,
	intersect bool,
	filter *dto.PolygonTargetFilter,
	zoom uint32,
) {
	if analysis == nil {
		return
	}

	count, err := u.parcelRepo.CountParcelsByPolygon(ctx, points, intersect, u.polygonFilterWithZoom(filter, zoom))
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[attachPolygonParcelCount] error: %v", err))
		return
	}
	analysis.ParcelCount = count
}

func validatePolygonQuery(points []repo.Point, maxAreaKm2 float64, page, pageSize int) (limit, offset int, err error) {
	if len(points) < config.MinPolygonVertices {
		return 0, 0, fmt.Errorf("polygon must have at least %d points", config.MinPolygonVertices)
	}
	if len(points) > config.MaxPolygonVertices {
		return 0, 0, fmt.Errorf("polygon has %d vertices", len(points))
	}

	if err := validatePolygon(points); err != nil {
		return 0, 0, fmt.Errorf("invalid polygon: %v", err)
	}

	area := calculatePolygonAreaKm2(points)
	allowedArea := config.MaxPolygonAreaKm2
	// if maxAreaKm2 > 0 && maxAreaKm2 < allowedArea {
	// 	allowedArea = maxAreaKm2
	// }
	if area > allowedArea {
		return 0, 0, fmt.Errorf("polygon area %.2f km² exceeds limit", area)
	}

	limit, offset = normalizePage(page, pageSize)
	return limit, offset, nil
}

func calculatePolygonAreaKm2(points []repo.Point) float64 {
	if len(points) < 3 {
		return 0
	}
	area := 0.0
	n := len(points)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += points[i].Lng * points[j].Lat
		area -= points[j].Lng * points[i].Lat
	}
	area = math.Abs(area) / 2.0
	return area * 111.0 * 111.0
}

func validatePolygon(points []repo.Point) error {
	n := len(points)
	for i := 0; i < n; i++ {
		for j := i + 2; j < n; j++ {
			if i == 0 && j == n-1 {
				continue
			}
			if i+1 == j {
				continue
			}
			p1, p2 := points[i], points[(i+1)%n]
			p3, p4 := points[j], points[(j+1)%n]
			if segmentsIntersect(p1, p2, p3, p4) {
				return fmt.Errorf("self-intersection between edges (%d,%d) and (%d,%d)", i, i+1, j, j+1)
			}
		}
	}
	return nil
}

func segmentsIntersect(a, b, c, d repo.Point) bool {
	orient := func(p, q, r repo.Point) int {
		val := (q.Lat-p.Lat)*(r.Lng-q.Lng) - (q.Lng-p.Lng)*(r.Lat-q.Lat)
		if val == 0 {
			return 0
		}
		if val > 0 {
			return 1
		}
		return 2
	}
	o1 := orient(a, b, c)
	o2 := orient(a, b, d)
	o3 := orient(c, d, a)
	o4 := orient(c, d, b)
	return o1 != o2 && o3 != o4
}

func normalizePageValues(page, pageSize int) (normalizedPage, normalizedPageSize int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = config.DefaultPageSize
	}
	if pageSize > config.MaxPageSize {
		pageSize = config.MaxPageSize
	}
	return page, pageSize
}

func normalizePage(page, pageSize int) (limit, offset int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = config.DefaultPageSize
	}
	if pageSize > config.MaxPageSize {
		pageSize = config.MaxPageSize
	}
	limit = pageSize
	offset = (page - 1) * pageSize
	return
}

func (u *parcelUsecase) GetParcelInfo(ctx context.Context, parcelID uint64) (*dto.ParcelInfoResponse, error) {
	if parcelID == 0 {
		return nil, fmt.Errorf("parcel_id is required")
	}

	info, err := u.parcelRepo.GetParcelInfo(ctx, parcelID)
	if err != nil {
		return nil, fmt.Errorf("get parcel info failed: %w", err)
	}

	if info == nil || info.ParcelID == 0 {
		return nil, fmt.Errorf("parcel %d not found", parcelID)
	}

	return info, nil
}

type parcelPlanningResolveOptions struct {
	ParcelID  uint64
	Mode      string
	Preset    string
	AsOfDate  string
	LogPrefix string
}

func (u *parcelUsecase) resolveParcelPlanning(ctx context.Context, opts parcelPlanningResolveOptions) (*parcel_engine.UnifiedResult, error) {
	if opts.ParcelID == 0 {
		return nil, fmt.Errorf("parcel_id is required")
	}
	if opts.Mode == "" {
		opts.Mode = "detail"
	}
	if opts.LogPrefix == "" {
		opts.LogPrefix = "resolveParcelPlanningUnified"
	}

	rows, err := u.parcelRepo.GetParcelLayerRowsV2(ctx, opts.ParcelID)
	if err != nil {
		return nil, fmt.Errorf("get parcel layers unified failed: %w", err)
	}

	if opts.Preset == "active_only" {
		rows = filterActiveLayers(rows)
	}
	if opts.AsOfDate != "" {
		rows = filterByAsOfDate(rows, opts.AsOfDate)
	}

	parcelInfo, infoErr := u.parcelRepo.GetParcelInfo(ctx, opts.ParcelID)
	if infoErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[%s] Failed to fetch parcel info for parcel %d: %v",
			opts.LogPrefix, opts.ParcelID, infoErr))
	}

	layerIDs := extractLayerIDsV2(rows)
	legalDocsMap, legalErr := u.layerLegalRepo.GetByLayerIDs(ctx, layerIDs)
	if legalErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[%s] Failed to fetch legal docs for parcel %d: %v",
			opts.LogPrefix, opts.ParcelID, legalErr))
		legalDocsMap = make(map[uint64][]*qh_domain.QHLayerLegal)
	}

	unifiedRows := parcel_engine.ToUnifiedRowsFromV2(rows)
	unified := u.parcelEngine.Process(unifiedRows, parcelInfo, convertLegalDocs(legalDocsMap), opts.Mode)
	if unified.ParcelID == 0 {
		unified.ParcelID = opts.ParcelID
	}

	return unified, nil
}

type ParcelPlanningResult struct {
	ParcelID   uint64
	View       string
	Mode       string
	Overview   *dto.ParcelLayersResponse
	Assessment *types.ResolveResult
	Metadata   *dto.ResponseMeta
}

func (u *parcelUsecase) GetParcelPlanning(
	ctx context.Context,
	req *dto.GetParcelPlanningRequest,
) (*ParcelPlanningResult, error) {
	return u.getParcelPlanningInternal(ctx, req)
}

func (u *parcelUsecase) getParcelPlanningInternal(
	ctx context.Context,
	req *dto.GetParcelPlanningRequest,
) (*ParcelPlanningResult, error) {
	import_time := time.Now()

	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}

	unified, err := u.resolveParcelPlanning(ctx, parcelPlanningResolveOptions{
		ParcelID:  req.ParcelID,
		Mode:      req.Mode,
		Preset:    req.Preset,
		AsOfDate:  req.AsOfDate,
		LogPrefix: "GetParcelPlanning",
	})
	if err != nil {
		return nil, err
	}

	meta := &dto.ResponseMeta{
		DataSource:     "postgis real-time (unified)",
		QueriedAt:      time.Now().UTC().Format(time.RFC3339),
		ResponseTimeMs: time.Since(import_time).Milliseconds(),
		Cached:         false,
	}

	result := &ParcelPlanningResult{
		ParcelID: unified.ParcelID,
		View:     req.View,
		Mode:     req.Mode,
		Metadata: meta,
	}

	if req.IncludeOverview {
		result.Overview = u.unifiedToParcelLayersResponse(unified)
		result.Overview.Metadata = meta
	}
	if req.IncludeAssessment {
		result.Assessment = u.unifiedToResolveResult(unified, req.Mode)
	}

	return result, nil
}

// ============================================================
// CONVERSION HELPERS (giữ nguyên)
// ============================================================

func (u *parcelUsecase) unifiedToParcelLayersResponse(ur *parcel_engine.UnifiedResult) *dto.ParcelLayersResponse {
	resp := &dto.ParcelLayersResponse{
		ParcelID:      ur.ParcelID,
		ParcelAreaSqm: ur.ParcelAreaSqm,
		ParcelInfo:    ur.ParcelInfo,
	}
	resp.PrimaryPlan = buildOverviewPrimaryPlan(ur.PrimaryPlanning)
	resp.LandUseGroups = buildOverviewLandUseGroups(ur.LandUseGroups)

	if ur.Risk != nil {
		resp.Risk = &dto.RiskAssessment{
			Level:          ur.Risk.Level,
			Score:          ur.Risk.Score,
			CanBuild:       ur.Risk.CanBuild,
			Reasons:        ur.Risk.Reasons,
			Recommendation: ur.Risk.Recommendation,
		}
	}

	resp.Layers = convertLayers(ur.Layers)
	resp.ConcludeLabel = pickConcludeLabelFromLayers(ur.Layers)

	totalOverlapRaw := 0.0
	for _, l := range ur.Layers {
		totalOverlapRaw += l.TotalArea
	}
	totalOverlapPct := 0.0
	if ur.ParcelAreaSqm > 0 {
		totalOverlapPct = (totalOverlapRaw / ur.ParcelAreaSqm) * 100
	}

	resp.Statistics = &dto.ImpactStats{
		TotalOverlapArea:     totalOverlapRaw,
		TotalOverlapPct:      totalOverlapPct,
		LayersAffected:       uint32(len(ur.Layers)),
		LabelsAffected:       countLabels(ur.Layers),
		RegionsAffected:      countRegions(ur.Layers),
		HasConflict:          ur.ParcelAreaSqm > 0 && totalOverlapRaw > ur.ParcelAreaSqm,
		SpatialDataAvailable: len(ur.Layers) > 0,
		Warning:              buildWarningFromUnified(ur),
	}

	resp.Summary = &dto.ParcelSummary{
		TotalAffectedPercentage: totalOverlapPct,
		Buildable:               ur.Risk != nil && ur.Risk.CanBuild,
	}
	if ur.PrimaryPlanning != nil {
		resp.Summary.DominantLandUse = ur.PrimaryPlanning.Name
		resp.Summary.DominantPercentage = ur.PrimaryPlanning.Ratio
		resp.Summary.TopLayerName = ur.PrimaryPlanning.LayerName
		resp.Summary.DominantLandUseCode = ur.PrimaryPlanning.Code
		resp.Summary.DominantLandUseColor = ur.PrimaryPlanning.Color
		resp.Summary.DominantLandUseGroupCode = ur.PrimaryPlanning.GroupCode
		resp.Summary.DominantLandUseGroupName = ur.PrimaryPlanning.GroupName
		resp.Summary.DominantAreaSqm = ur.PrimaryPlanning.AreaSqm
		resp.Summary.DominantCanBuild = ur.PrimaryPlanning.CanBuild
	} else if ur.ParcelInfo != nil {
		resp.Summary.DominantLandUse = ur.ParcelInfo.LandUseName
		resp.Summary.DominantLandUseCode = ur.ParcelInfo.LandUseCode
		resp.Summary.DominantLandUseColor = ur.ParcelInfo.LandUseColor
	}

	return resp
}

func buildOverviewPrimaryPlan(p *parcel_engine.PrimaryPlanningInfo) *dto.ParcelPlanInfo {
	if p == nil {
		return nil
	}
	return &dto.ParcelPlanInfo{
		LayerID:          p.LayerID,
		LayerName:        p.LayerName,
		DisplayName:      p.LayerName,
		LandUseCode:      p.Code,
		LandUseName:      p.Name,
		LandUseColor:     p.Color,
		GroupCode:        p.GroupCode,
		GroupName:        p.GroupName,
		CanBuild:         p.CanBuild,
		ImpactPercent:    p.Ratio,
		ImpactAreaSqm:    p.AreaSqm,
		LegalCanonical:   p.LegalCanonical,
		LegalConfidence:  p.LegalConfidence,
		LegalActiveState: p.LegalActiveState,
	}
}

func buildOverviewLandUseGroups(groups []*parcel_engine.LandUseGroupInfo) []*dto.PlanLandUseGroupDetail {
	if len(groups) == 0 {
		return nil
	}
	result := make([]*dto.PlanLandUseGroupDetail, 0, len(groups))
	for _, g := range groups {
		if g == nil {
			continue
		}
		if len(g.Details) == 0 {
			result = append(result, &dto.PlanLandUseGroupDetail{
				Code:           g.Code,
				Name:           g.Name,
				Color:          g.Color,
				CanBuild:       g.CanBuild,
				LandUseName:    g.Name,
				LandUseColor:   g.Color,
				OverlapPercent: g.Ratio,
				OverlapAreaSqm: g.AreaSqm,
			})
			continue
		}
		for _, d := range g.Details {
			if d == nil {
				continue
			}
			result = append(result, &dto.PlanLandUseGroupDetail{
				Code:           g.Code,
				Name:           g.Name,
				Color:          g.Color,
				CanBuild:       g.CanBuild,
				LandUseCode:    d.LandUseCode,
				LandUseName:    d.LandUseName,
				LandUseColor:   g.Color,
				OverlapPercent: d.Ratio,
				OverlapAreaSqm: d.AreaSqm,
			})
		}
	}
	return result
}

func (u *parcelUsecase) unifiedToResolveResult(ur *parcel_engine.UnifiedResult, mode string) *types.ResolveResult {
	result := &types.ResolveResult{
		HasConflict:      len(ur.Conflicts) > 0,
		Mode:             types.ResolveMode(mode),
		ProcessingTimeMs: ur.Metadata.ResponseTimeMs,
		ParcelAreaSqm:    ur.ParcelAreaSqm,
		ResolutionStatus: ur.ResolutionStatus,
	}

	if ur.CompareResult != nil {
		result.CompareResult = &types.CompareResultInfo{
			TotalLayers: ur.CompareResult.TotalLayers,
		}
		for _, c := range ur.CompareResult.Changes {
			result.CompareResult.Changes = append(result.CompareResult.Changes, &types.CompareChangeInfo{
				LayerID:    c.LayerID,
				LayerName:  c.LayerName,
				ChangeType: c.ChangeType,
				OldValue:   c.OldValue,
				NewValue:   c.NewValue,
			})
		}
	}

	if ur.HistoricalRisk != nil {
		result.HistoricalRisk = &types.HistoricalRiskInfo{
			Trend:          ur.HistoricalRisk.Trend,
			TrendDetail:    ur.HistoricalRisk.TrendDetail,
			ActiveLayers:   ur.HistoricalRisk.ActiveLayers,
			ExpiringLayers: ur.HistoricalRisk.ExpiringLayers,
			AnalyzedAt:     ur.HistoricalRisk.AnalyzedAt,
		}
	}

	if ur.Risk != nil {
		result.CanBuild = ur.Risk.CanBuild
		result.RiskScore = ur.Risk.Score
		result.RiskLevel = ur.Risk.Level
		result.RiskReasons = ur.Risk.Reasons
		result.Recommendation = ur.Risk.Recommendation
	}

	if ur.PrimaryPlanning != nil {
		result.BuildStatus = u.buildStatusFromRisk(ur.Risk, len(ur.Conflicts) > 0)
		result.Reason = u.buildReasonFromPrimary(ur.PrimaryPlanning, len(ur.Conflicts) > 0)
	}

	for _, c := range ur.Conflicts {
		result.CriticalConflicts = append(result.CriticalConflicts, &types.ConflictDetail{
			Type:            c.Type,
			Severity:        types.ConflictSeverity(c.Severity),
			Between:         c.Between,
			Description:     c.Description,
			Recommendation:  c.Recommendation,
			ActionRequired:  c.ActionRequired,
			AffectedPercent: c.AffectedPercent,
		})
	}

	result.Primary = buildRankedLayerFromPrimary(ur.PrimaryPlanning, ur.ParcelAreaSqm)
	for _, g := range ur.LandUseGroups {
		for _, d := range g.Details {
			if result.Primary != nil && result.Primary.LayerID == ur.PrimaryPlanning.LayerID {
				if result.LayerLandUseGroups == nil {
					result.LayerLandUseGroups = make(map[uint64][]*types.LandUseGroupDetail)
				}
				result.LayerLandUseGroups[ur.PrimaryPlanning.LayerID] = append(
					result.LayerLandUseGroups[ur.PrimaryPlanning.LayerID],
					&types.LandUseGroupDetail{
						Code:           g.Code,
						Name:           g.Name,
						Color:          g.Color,
						CanBuild:       g.CanBuild,
						LandUseCode:    d.LandUseCode,
						LandUseName:    d.LandUseName,
						OverlapPercent: d.Ratio,
						OverlapAreaSqm: d.AreaSqm,
					},
				)
			}
		}
	}

	if len(ur.SecondaryPlans) > 0 {
		for _, sp := range ur.SecondaryPlans {
			result.Secondary = append(result.Secondary, &types.RankedLayer{
				EvaluatedLayer: &types.EvaluatedLayer{
					LayerCandidate: &types.LayerCandidate{
						LayerID:          sp.LayerID,
						LayerDisplayName: sp.LayerName,
						LandUseCode:      sp.LandUseCode,
						LandUseName:      sp.LandUseName,
						LandUseColor:     sp.LandUseColor,
						GroupCode:        sp.GroupCode,
						OverlapPercent:   sp.OverlapPercent,
						OverlapAreaSqm:   sp.OverlapAreaSqm,
						ParcelAreaSqm:    ur.ParcelAreaSqm,
					},
				},
				Rank:  2,
				Score: 50,
			})
		}
	} else {
		for _, l := range ur.Layers {
			if ur.PrimaryPlanning != nil && l.ID != ur.PrimaryPlanning.LayerID {
				sec := buildRankedLayerFromLayerDetail(l, ur.ParcelAreaSqm)
				if sec != nil {
					result.Secondary = append(result.Secondary, sec)
				}
			}
		}
	}

	if mode == "detail" || mode == "" {
		result.Layers = convertEngineLayersToTypes(ur.Layers)
	}

	return result
}

func buildRankedLayerFromPrimary(p *parcel_engine.PrimaryPlanningInfo, parcelArea float64) *types.RankedLayer {
	if p == nil {
		return nil
	}
	return &types.RankedLayer{
		EvaluatedLayer: &types.EvaluatedLayer{
			LayerCandidate: &types.LayerCandidate{
				LayerID:          p.LayerID,
				LayerDisplayName: p.LayerName,
				LandUseCode:      p.Code,
				LandUseName:      p.Name,
				LandUseColor:     p.Color,
				GroupCode:        p.GroupCode,
				GroupName:        p.GroupName,
				CanBuild:         p.CanBuild,
				LegalStatus:      p.LegalStatus,
				OverlapPercent:   p.Ratio,
				OverlapAreaSqm:   p.AreaSqm,
				ParcelAreaSqm:    parcelArea,
			},
		},
		Rank:  1,
		Score: 100,
	}
}

func buildRankedLayerFromLayerDetail(l *parcel_engine.LayerDetail, parcelArea float64) *types.RankedLayer {
	if l == nil {
		return nil
	}
	code := ""
	name := ""
	color := ""
	if len(l.Labels) > 0 {
		code = l.Labels[0].Code
		name = l.Labels[0].Name
		color = l.Labels[0].Color
	}
	return &types.RankedLayer{
		EvaluatedLayer: &types.EvaluatedLayer{
			LayerCandidate: &types.LayerCandidate{
				LayerID:          l.ID,
				LayerName:        l.Name,
				LayerDisplayName: l.DisplayName,
				LandUseCode:      code,
				LandUseName:      name,
				LandUseColor:     color,
				OverlapPercent:   l.TotalPct,
				OverlapAreaSqm:   l.TotalArea,
				ParcelAreaSqm:    parcelArea,
			},
		},
		Rank:  2,
		Score: 50,
	}
}

func (u *parcelUsecase) buildStatusFromRisk(risk *parcel_engine.RiskInfo, hasConflict bool) string {
	cfg := u.parcelEngine.GetConfig()
	if risk == nil {
		return "unknown"
	}
	if !risk.CanBuild && risk.Level == "CRITICAL" {
		return cfg.StatusProhibited
	}
	if hasConflict {
		return cfg.StatusConditional
	}
	return cfg.StatusAllowed
}

func (u *parcelUsecase) buildReasonFromPrimary(p *parcel_engine.PrimaryPlanningInfo, hasConflict bool) string {
	if p == nil {
		return ""
	}
	reason := fmt.Sprintf("%s (%s) — %s", p.Name, p.Code, p.LayerName)
	if p.LegalCanonical != "" {
		reason += fmt.Sprintf(" [Pháp lý: %s, độ tin cậy: %s]", p.LegalCanonical, p.LegalConfidence)
	}
	if hasConflict {
		reason += " — " + u.parcelEngine.GetConfig().ConflictStatusSuffix
	}
	return reason
}

func buildWarningFromUnified(u *parcel_engine.UnifiedResult) string {
	if len(u.Warnings) > 0 {
		return u.Warnings[0]
	}
	return ""
}

func convertLayers(layers []*parcel_engine.LayerDetail) []*dto.LayerImpact {
	result := make([]*dto.LayerImpact, 0, len(layers))
	for _, l := range layers {
		impact := &dto.LayerImpact{
			ID:              l.ID,
			Name:            l.Name,
			DisplayName:     l.DisplayName,
			Type:            l.Type,
			EffectiveDate:   l.EffectiveDate,
			TotalArea:       l.TotalArea,
			TotalPct:        l.TotalPct,
			LegalStatus:     l.LegalStatus,
			LegalStatusName: l.LegalStatusName,
			LayerAvatar:     l.LayerAvatar,
		}

		if l.Authority != nil {
			impact.AuthorityIssuing = &dto.AuthorityIssuingDTO{
				ID:          l.Authority.ID,
				Name:        l.Authority.Name,
				Code:        l.Authority.Code,
				Description: l.Authority.Description,
			}
		}

		for _, doc := range l.Documents {
			impact.LegalDocs = append(impact.LegalDocs, &dto.LayerLegalDTO{
				ID:       doc.ID,
				Name:     doc.Name,
				FileUrl:  doc.FileUrl,
				FileType: doc.FileType,
			})
		}

		for _, lb := range l.Labels {
			labelImpact := &dto.LabelImpact{
				ID:          lb.ID,
				Name:        lb.Code,
				DisplayName: lb.Name,
				Color:       lb.Color,
				FontSize:    12,
				FontWeight:  "normal",
				MinZoom:     10,
				MaxZoom:     22,
				TotalArea:   lb.AreaSqm,
				TotalPct:    lb.Ratio,
			}
			for _, reg := range lb.Regions {
				labelImpact.Regions = append(labelImpact.Regions, &dto.RegionImpact{
					ID:             reg.ID,
					Name:           reg.Name,
					DisplayName:    reg.DisplayName,
					LandUseCode:    reg.LandUseCode,
					LandUseName:    reg.LandUseName,
					ColorRGB:       reg.LandUseColor,
					OverlapAreaSqm: reg.OverlapAreaSqm,
					OverlapPct:     reg.OverlapPct,
					CenterLat:      reg.CenterLat,
					CenterLng:      reg.CenterLng,
					Geometry:       reg.Geometry,
					WarnLevel:      reg.WarnLevel,
					CanBuild:       reg.CanBuild,
				})
			}
			impact.Labels = append(impact.Labels, labelImpact)
		}
		result = append(result, impact)
	}
	return result
}

func pickConcludeLabelFromLayers(layers []*parcel_engine.LayerDetail) *dto.RegionImpact {
	var best *dto.RegionImpact
	for _, l := range layers {
		for _, lb := range l.Labels {
			for _, reg := range lb.Regions {
				if best == nil || reg.OverlapAreaSqm > best.OverlapAreaSqm {
					best = &dto.RegionImpact{
						ID:             reg.ID,
						Name:           reg.Name,
						DisplayName:    reg.DisplayName,
						LandUseCode:    reg.LandUseCode,
						LandUseName:    reg.LandUseName,
						ColorRGB:       reg.LandUseColor,
						OverlapAreaSqm: reg.OverlapAreaSqm,
						OverlapPct:     reg.OverlapPct,
						CenterLat:      reg.CenterLat,
						CenterLng:      reg.CenterLng,
						Geometry:       reg.Geometry,
						WarnLevel:      reg.WarnLevel,
						CanBuild:       reg.CanBuild,
					}
				}
			}
		}
	}
	return best
}

func countLabels(layers []*parcel_engine.LayerDetail) uint32 {
	seen := make(map[uint64]bool)
	for _, l := range layers {
		for _, lb := range l.Labels {
			seen[lb.ID] = true
		}
	}
	return uint32(len(seen))
}

func countRegions(layers []*parcel_engine.LayerDetail) uint32 {
	seen := make(map[uint64]bool)
	for _, l := range layers {
		for _, lb := range l.Labels {
			for _, reg := range lb.Regions {
				seen[reg.ID] = true
			}
		}
	}
	return uint32(len(seen))
}

func convertLegalDocs(legalDocsMap map[uint64][]*qh_domain.QHLayerLegal) map[uint64][]*dto.LayerLegalDTO {
	result := make(map[uint64][]*dto.LayerLegalDTO, len(legalDocsMap))
	for layerID, legals := range legalDocsMap {
		dtos := make([]*dto.LayerLegalDTO, 0, len(legals))
		for _, l := range legals {
			dtos = append(dtos, &dto.LayerLegalDTO{
				ID:       l.ID,
				LayerID:  l.LayerID,
				Name:     l.Name,
				FileUrl:  l.FileURL,
				FileType: l.FileType,
			})
		}
		result[layerID] = dtos
	}
	return result
}

func extractLayerIDs(rows []dto.ParcelLayerRow) []uint64 {
	seen := make(map[uint64]bool)
	for _, r := range rows {
		seen[r.LayerID] = true
	}
	ids := make([]uint64, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}

func extractLayerIDsV2(rows []dto.ParcelLayerRowV2) []uint64 {
	seen := make(map[uint64]bool)
	for _, r := range rows {
		seen[r.LayerID] = true
	}
	ids := make([]uint64, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}

func filterActiveLayers(rows []dto.ParcelLayerRowV2) []dto.ParcelLayerRowV2 {
	now := time.Now()
	filtered := make([]dto.ParcelLayerRowV2, 0, len(rows))
	for _, row := range rows {
		if row.LayerExpiryDate != "" {
			expiry, err := time.Parse("2006-01-02", row.LayerExpiryDate)
			if err == nil && expiry.Before(now) {
				continue
			}
		}
		filtered = append(filtered, row)
	}
	return filtered
}

func filterByAsOfDate(rows []dto.ParcelLayerRowV2, asOfDate string) []dto.ParcelLayerRowV2 {
	targetDate, err := time.Parse("2006-01-02", asOfDate)
	if err != nil {
		return rows
	}
	filtered := make([]dto.ParcelLayerRowV2, 0, len(rows))
	for _, row := range rows {
		if row.LayerEffectiveDate == "" {
			filtered = append(filtered, row)
			continue
		}
		effectiveDate, err := time.Parse("2006-01-02", row.LayerEffectiveDate)
		if err != nil || effectiveDate.Before(targetDate) || effectiveDate.Equal(targetDate) {
			filtered = append(filtered, row)
		}
	}
	return filtered
}

func convertEngineLayersToTypes(engineLayers []*parcel_engine.LayerDetail) []*types.LayerDetail {
	if len(engineLayers) == 0 {
		return nil
	}
	result := make([]*types.LayerDetail, len(engineLayers))
	for i, el := range engineLayers {
		result[i] = convertEngineLayerToType(el)
	}
	return result
}

func convertEngineLayerToType(el *parcel_engine.LayerDetail) *types.LayerDetail {
	if el == nil {
		return nil
	}
	l := &types.LayerDetail{
		ID:              el.ID,
		Name:            el.Name,
		DisplayName:     el.DisplayName,
		Type:            el.Type,
		EffectiveDate:   el.EffectiveDate,
		TotalArea:       el.TotalArea,
		TotalPct:        el.TotalPct,
		LegalStatus:     el.LegalStatus,
		LegalStatusName: el.LegalStatusName,
		LayerAvatar:     el.LayerAvatar,
	}
	if el.Authority != nil {
		l.Authority = &types.AuthorityDetail{
			ID:          el.Authority.ID,
			Name:        el.Authority.Name,
			Code:        el.Authority.Code,
			Description: el.Authority.Description,
		}
	}
	for _, d := range el.Documents {
		l.LegalDocs = append(l.LegalDocs, &types.LegalDocDetail{
			ID:       d.ID,
			Name:     d.Name,
			FileUrl:  d.FileUrl,
			FileType: d.FileType,
		})
	}
	for _, lbl := range el.Labels {
		tl := &types.LabelDetail{
			ID:        lbl.ID,
			Name:      lbl.Name,
			Color:     lbl.Color,
			TotalArea: lbl.AreaSqm,
			TotalPct:  lbl.Ratio,
		}
		for _, r := range lbl.Regions {
			tr := &types.RegionDetail{
				ID:             r.ID,
				Name:           r.Name,
				DisplayName:    r.DisplayName,
				LandUseCode:    r.LandUseCode,
				LandUseName:    r.LandUseName,
				LandUseGroup:   r.LandUseGroup,
				LandUseColor:   r.LandUseColor,
				OverlapAreaSqm: r.OverlapAreaSqm,
				OverlapPct:     r.OverlapPct,
				Geometry:       r.Geometry,
				WarnLevel:      r.WarnLevel,
				CanBuild:       r.CanBuild,
			}
			tl.Regions = append(tl.Regions, tr)
		}
		l.Labels = append(l.Labels, tl)
	}
	return l
}

func (u *parcelUsecase) GetParcelQuickInfo(ctx context.Context, parcelID uint64) ([]dto.ParcelLayerInfo, error) {
	quickInfo, err := u.parcelRepo.GetParcelQuickInfo(ctx, parcelID)
	if err != nil {
		return nil, fmt.Errorf("get parcel quick info failed: %w", err)
	}
	return quickInfo, nil
}

func (u *parcelUsecase) GetParcelTimeline(ctx context.Context, parcelID uint64) (*types.ParcelTimeline, error) {
	if parcelID == 0 {
		return nil, fmt.Errorf("parcel_id is required")
	}

	if u.timelineEngine != nil {
		return u.timelineEngine.BuildTimeline(ctx, parcelID)
	}

	rows, err := u.parcelRepo.GetParcelLayerRowsV2(ctx, parcelID)
	if err != nil {
		return nil, fmt.Errorf("get parcel layers failed: %w", err)
	}

	events := make([]*types.TimelineEvent, 0, len(rows))
	for _, row := range rows {
		eventType := types.EventLayerEffective
		if row.LayerEffectiveDate != "" {
			eventType = types.EventLayerEffective
		}
		events = append(events, &types.TimelineEvent{
			ID:          fmt.Sprintf("evt_%d_%d", row.LayerID, row.RegionID),
			Date:        time.Now(),
			EventType:   eventType,
			LayerID:     row.LayerID,
			LayerName:   row.LayerDisplayName,
			Description: fmt.Sprintf("Layer %s — %s (%.1f%%)", row.LayerDisplayName, row.RegionLandUseName, row.OverlapPct),
		})
	}

	summary := &types.TimelineSummary{
		TotalEvents: len(events),
	}
	if len(events) > 0 {
		summary.FirstEventAt = events[0].Date
		summary.LastEventAt = events[len(events)-1].Date
	}

	return &types.ParcelTimeline{
		ParcelID: parcelID,
		Events:   events,
		Summary:  summary,
	}, nil
}

func (u *parcelUsecase) GetLayerHistory(ctx context.Context, layerID uint64) (*types.LayerLineage, error) {
	if layerID == 0 {
		return nil, fmt.Errorf("layer_id is required")
	}

	if u.lifecycleEngine != nil {
		return u.lifecycleEngine.GetLineage(ctx, layerID)
	}

	return &types.LayerLineage{
		RootLayerID: layerID,
		FamilyID:    layerID,
		CurrentID:   layerID,
		Versions: []*types.LayerVersion{
			{
				Version:   1,
				LayerID:   layerID,
				FamilyID:  layerID,
				State:     types.StateEffective,
				CreatedAt: time.Now(),
				IsLatest:  true,
			},
		},
	}, nil
}

func (u *parcelUsecase) GetParcelSeoSource(ctx context.Context, parcelID uint64) (*dto.ParcelSeoSource, error) {
	if parcelID == 0 {
		return nil, fmt.Errorf("parcel_id is required")
	}
	return u.parcelRepo.GetParcelSeoSource(ctx, parcelID)
}

func (u *parcelUsecase) ListParcelSeoSourcesForGenerate(ctx context.Context, limit uint32) ([]dto.ParcelSeoSource, error) {
	return u.parcelRepo.ListParcelSeoSourcesForGenerate(ctx, limit)
}

func (u *parcelUsecase) UpdateParcelSeoID(ctx context.Context, parcelID uint64, seoID uint64) error {
	if parcelID == 0 || seoID == 0 {
		return fmt.Errorf("parcel_id và seo_id là bắt buộc")
	}
	return u.parcelRepo.UpdateParcelSeoID(ctx, parcelID, seoID)
}

// GetParcelDetail — Lấy tổng quan thửa đất
func (u *parcelUsecase) GetParcelDetail(ctx context.Context, parcelID uint64) (*qh_dto.ParcelDetailResponseDTO, error) {
	if parcelID == 0 {
		return nil, fmt.Errorf("parcel_id is required")
	}

	parcelInfo, err := u.parcelRepo.GetParcelInfo(ctx, parcelID)
	if err != nil {
		return nil, fmt.Errorf("get parcel info failed: %w", err)
	}
	if parcelInfo == nil || parcelInfo.ParcelID == 0 {
		return nil, fmt.Errorf("parcel %d not found", parcelID)
	}

	rows, err := u.parcelRepo.GetParcelLayerRowsV2(ctx, parcelID)
	if err != nil {
		return nil, fmt.Errorf("get parcel layers failed: %w", err)
	}

	// 3. Lấy legal documents
	layerIDs := extractLayerIDsV2(rows)
	legalDocsMap, err := u.layerLegalRepo.GetByLayerIDs(ctx, layerIDs)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[GetParcelDetail] Failed to fetch legal docs: %v", err))
		legalDocsMap = make(map[uint64][]*qh_domain.QHLayerLegal)
	}

	unifiedRows := parcel_engine.ToUnifiedRowsFromV2(rows)
	unified := u.parcelEngine.Process(unifiedRows, parcelInfo, convertLegalDocs(legalDocsMap), "detail")
	if unified.ParcelID == 0 {
		unified.ParcelID = parcelID
	}

	resp := &qh_dto.ParcelDetailResponseDTO{
		Parcel: qh_dto.ParcelIdentityDTO{
			ID:               strconv.FormatUint(parcelInfo.ParcelID, 10),
			MapSheetNumber:   parcelInfo.MapNumber,
			LandParcelNumber: parcelInfo.LandNumber,
			Address:          parcelInfo.AddressText,
			TotalAreaSqm:     parcelInfo.AreaSqm,
			Centroid: &qh_dto.CentroidDTO{
				Lat: parcelInfo.Lat,
				Lng: parcelInfo.Lon,
			},
			Location: &qh_dto.LocationDTO{
				Province:     parcelInfo.Province,
				ProvinceCode: parcelInfo.ProvinceCode,
				WardCode:     parcelInfo.WardCode,
			},
		},
		CurrentUse: qh_dto.CurrentUseDTO{
			Code:            parcelInfo.LandUseCode,
			Name:            parcelInfo.LandUseName,
			Color:           parcelInfo.LandUseColor,
			BuildStatus:     u.getBuildStatusString(unified),
			LegalDocumentId: "",
		},
	}

	if unified.PrimaryPlanning != nil {
		resp.PlanningConclusion = u.buildPlanningConclusion(unified)
	}

	if unified.PrimaryPlanning != nil {
		resp.PrimaryUse = u.buildPrimaryUse(unified)
	}

	resp.Guidance = u.buildGuidance(unified)

	legalDocs := u.collectLegalDocumentsFromMap(legalDocsMap)
	resp.LegalDocuments = u.buildLegalDocumentsSummary(legalDocs)

	resp.Meta = &qh_dto.ResponseMetaDTO{
		DataVersion: "2.1.0",
		DataQuality: "VERIFIED",
		QueriedAt:   time.Now().UTC().Format(time.RFC3339),
	}
	resp.Links = &qh_dto.LinksDTO{
		Self:           fmt.Sprintf("/v2/tqd/parcels/%d", parcelID),
		Layers:         fmt.Sprintf("/v2/tqd/parcels/%d/layers", parcelID),
		LegalDocuments: fmt.Sprintf("/v2/tqd/parcels/%d/legal-documents", parcelID),
		ShareUrl:       fmt.Sprintf("/map?parcel=%d", parcelID),
	}

	return resp, nil
}

func (u *parcelUsecase) GetParcelLayers(ctx context.Context, parcelID uint64, pagable _dto.Pagable) ([]*dto.LayerImpact, int64, error) {
	if parcelID == 0 {
		return nil, 0, fmt.Errorf("parcel_id is required")
	}

	rows, err := u.parcelRepo.GetParcelLayerRowsV2(ctx, parcelID)
	if err != nil {
		return nil, 0, fmt.Errorf("get parcel layers failed: %w", err)
	}

	layerIDs := extractLayerIDsV2(rows)
	legalDocsMap, err := u.layerLegalRepo.GetByLayerIDs(ctx, layerIDs)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[GetParcelLayers] Failed to fetch legal docs: %v", err))
		legalDocsMap = make(map[uint64][]*qh_domain.QHLayerLegal)
	}

	unifiedRows := parcel_engine.ToUnifiedRowsFromV2(rows)
	unified := u.parcelEngine.Process(unifiedRows, nil, convertLegalDocs(legalDocsMap), "detail")
	if unified.ParcelID == 0 {
		unified.ParcelID = parcelID
	}

	layers := convertLayers(unified.Layers)

	offset := pagable.GetOffset()
	limit := pagable.GetLimit()

	total := int64(len(layers))
	if offset >= len(layers) {
		return []*dto.LayerImpact{}, total, nil
	}

	end := offset + limit
	if end > len(layers) {
		end = len(layers)
	}

	return layers[offset:end], total, nil
}

// GetZoneGeometry — Lấy GeoJSON của một zone cụ thể
func (u *parcelUsecase) GetZoneGeometry(ctx context.Context, parcelID, layerID, zoneID uint64) (string, *ZoneGeometryProps, error) {
	if parcelID == 0 {
		return "", nil, fmt.Errorf("parcel_id is required")
	}
	if layerID == 0 {
		return "", nil, fmt.Errorf("layer_id is required")
	}
	if zoneID == 0 {
		return "", nil, fmt.Errorf("zone_id is required")
	}

	// Lấy geometry từ repository
	geoJSON, region, err := u.parcelRepo.GetZoneGeometry(ctx, parcelID, layerID, zoneID)
	if err != nil {
		return "", nil, fmt.Errorf("get zone geometry failed: %w", err)
	}
	if geoJSON == "" {
		return "", nil, fmt.Errorf("zone not found")
	}

	// Build properties
	props := &ZoneGeometryProps{
		ZoneName:         region.Name,
		LandUseCode:      region.LandUseCode,
		LandUseName:      region.LandUseName,
		LandUseColor:     region.LandUseColor,
		BuildStatus:      u.determineBuildStatus(region.CanBuild),
		AlertLevel:       u.determineAlertLevel(region.WarnLevel),
		LegalDocumentIds: []string{},
	}

	return geoJSON, props, nil
}

func (u *parcelUsecase) GetParcelLegalDocuments(ctx context.Context, parcelID uint64, docType, status *uint32, pagable _dto.Pagable) ([]*qh_domain.QHLegalDocument, int64, error) {
	if parcelID == 0 {
		return nil, 0, fmt.Errorf("parcel_id is required")
	}

	docs, err := u.legalDocRepo.GetByParcelID(ctx, parcelID)
	if err != nil {
		return nil, 0, fmt.Errorf("get legal documents failed: %w", err)
	}

	// Filter
	filtered := make([]*qh_domain.QHLegalDocument, 0)
	for _, d := range docs {
		if docType != nil && uint32(d.DocumentType) != *docType {
			continue
		}
		if status != nil && uint32(d.Status) != *status {
			continue
		}
		filtered = append(filtered, d)
	}

	// ✅ Dùng offset/limit từ pagable (an toàn)
	offset := pagable.GetOffset()
	limit := pagable.GetLimit()

	total := int64(len(filtered))
	if offset >= len(filtered) {
		return []*qh_domain.QHLegalDocument{}, total, nil
	}

	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[offset:end], total, nil
}

// getBuildStatusString — Lấy build status từ UnifiedResult
func (u *parcelUsecase) getBuildStatusString(unified *parcel_engine.UnifiedResult) string {
	if unified.Risk != nil {
		if unified.Risk.CanBuild {
			return "ALLOWED"
		}
		return "FORBIDDEN"
	}
	return "UNKNOWN"
}

// buildPlanningConclusion — Xây dựng PlanningConclusion từ UnifiedResult
func (u *parcelUsecase) buildPlanningConclusion(unified *parcel_engine.UnifiedResult) *qh_dto.PlanningConclusionDTO {
	pc := &qh_dto.PlanningConclusionDTO{
		BuildStatus: u.getBuildStatusString(unified),
		HasConflict: len(unified.Conflicts) > 0,
	}

	if unified.Risk != nil {
		pc.Alert = &qh_dto.AlertDTO{
			Level:   unified.Risk.Level,
			Score:   unified.Risk.Score,
			Message: u.buildAlertMessage(unified),
		}
	}

	// Breakdown từ LandUseGroups
	if unified.LandUseGroups != nil {
		for idx, g := range unified.LandUseGroups {
			if len(g.Details) == 0 {
				pc.Breakdown = append(pc.Breakdown, qh_dto.PlanningBreakdownDTO{
					Rank:         idx + 1,
					LandUseCode:  g.Code,
					LandUseName:  g.Name,
					LandUseColor: g.Color,
					GroupCode:    g.Code,
					GroupName:    g.Name,
					AreaSqm:      g.AreaSqm,
					Percent:      g.Ratio,
					BuildStatus:  u.getBuildStatusFromBool(g.CanBuild),
					AlertLevel:   "NONE",
				})
				continue
			}
			for _, d := range g.Details {
				pc.Breakdown = append(pc.Breakdown, qh_dto.PlanningBreakdownDTO{
					Rank:         idx + 1,
					LandUseCode:  d.LandUseCode,
					LandUseName:  d.LandUseName,
					LandUseColor: g.Color,
					GroupCode:    g.Code,
					GroupName:    g.Name,
					AreaSqm:      d.AreaSqm,
					Percent:      d.Ratio,
					BuildStatus:  u.getBuildStatusFromBool(g.CanBuild),
					AlertLevel:   "NONE",
				})
			}
		}
	}

	return pc
}

// buildPrimaryUse — Xây dựng PrimaryUse từ UnifiedResult
func (u *parcelUsecase) buildPrimaryUse(unified *parcel_engine.UnifiedResult) *qh_dto.PrimaryUseDTO {
	if unified.PrimaryPlanning == nil {
		return nil
	}
	p := unified.PrimaryPlanning
	return &qh_dto.PrimaryUseDTO{
		Rule:        "LARGEST_AREA",
		LandUseCode: p.Code,
		LandUseName: p.Name,
		Percent:     p.Ratio,
		BuildStatus: u.getBuildStatusFromBool(p.CanBuild),
	}
}

func (u *parcelUsecase) buildGuidance(unified *parcel_engine.UnifiedResult) *qh_dto.GuidanceDTO {
	guidance := &qh_dto.GuidanceDTO{
		Warnings:        []qh_dto.WarningDTO{},
		Recommendations: []qh_dto.RecommendationDTO{},
	}

	if unified.Risk != nil {
		for _, reason := range unified.Risk.Reasons {
			guidance.Warnings = append(guidance.Warnings, qh_dto.WarningDTO{
				AlertLevel: unified.Risk.Level,
				Title:      reason,
				Action:     unified.Risk.Recommendation,
			})
		}
	}

	return guidance
}

// buildLegalDocumentsSummary — Tóm tắt legal documents
func (u *parcelUsecase) buildLegalDocumentsSummary(docs []*qh_domain.QHLegalDocument) *qh_dto.LegalDocumentsSummaryDTO {
	summary := &qh_dto.LegalDocumentsSummaryDTO{
		Items:             []qh_dto.LegalDocumentDTO{},
		ByType:            make(map[string]int),
		ByStatus:          make(map[string]int),
		HasSupersededDocs: false,
	}

	for _, d := range docs {
		// Convert domain → DTO
		item := qh_dto.LegalDocumentDTO{
			ID:            strconv.FormatUint(d.ID, 10),
			Type:          enums.DocumentTypeString[d.DocumentType],
			TypeName:      enums.DocumentTypeName[d.DocumentType],
			Name:          d.Name,
			FullName:      d.FullName,
			IssuedBy:      d.IssuedBy,
			IssuedAt:      _utils.FormatRFC3339(d.IssuedAt),
			EffectiveDate: _utils.FormatRFC3339(d.EffectiveAt),
			Status:        enums.DocumentStatusString[d.Status],
			StatusLabel:   enums.DocumentStatusLabel[d.Status],
			StatusColor:   enums.DocumentStatusColor[d.Status],
			IsPublic:      d.IsPublic,
			RequiresAuth:  d.RequiresAuth,
		}
		if d.ExpiresAt != nil {
			expiry := _utils.FormatRFC3339(d.ExpiresAt)
			item.ExpiryDate = &expiry
		}
		if d.SupersededBy != nil {
			sb := strconv.FormatUint(*d.SupersededBy, 10)
			item.SupersededBy = &sb
		}
		for _, s := range d.Supersedes {
			item.Supersedes = append(item.Supersedes, strconv.FormatUint(s, 10))
		}

		summary.Items = append(summary.Items, item)
		summary.ByType[item.Type]++
		summary.ByStatus[item.Status]++
		if d.IsSuperseded() {
			summary.HasSupersededDocs = true
		}
	}
	summary.Total = len(summary.Items)

	return summary
}

func (u *parcelUsecase) collectLegalDocumentsFromMap(legalDocsMap map[uint64][]*qh_domain.QHLayerLegal) []*qh_domain.QHLegalDocument {
	seen := make(map[uint64]bool)
	var docs []*qh_domain.QHLegalDocument

	for _, list := range legalDocsMap {
		for _, d := range list {
			doc := &qh_domain.QHLegalDocument{
				DocumentType: 0,
				Name:         d.Name,
				FileURL:      d.FileURL,
				FileType:     d.FileType,
			}
			doc.ID = d.ID

			if !seen[doc.ID] {
				seen[doc.ID] = true
				docs = append(docs, doc)
			}
		}
	}
	return docs
}

func (u *parcelUsecase) getBuildStatusFromBool(canBuild bool) string {
	if canBuild {
		return "ALLOWED"
	}
	return "FORBIDDEN"
}

func (u *parcelUsecase) determineBuildStatus(canBuild bool) enums.BuildStatus {
	if canBuild {
		return enums.BuildStatusAllowed
	}
	return enums.BuildStatusForbidden
}

func (u *parcelUsecase) determineAlertLevel(warnLevel int32) enums.AlertLevel {
	switch warnLevel {
	case 1:
		return enums.AlertLevelLow
	case 2:
		return enums.AlertLevelMedium
	case 3:
		return enums.AlertLevelHigh
	default:
		return enums.AlertLevelNone
	}
}

// buildAlertMessage — Xây dựng message alert
func (u *parcelUsecase) buildAlertMessage(unified *parcel_engine.UnifiedResult) string {
	if unified.Risk == nil {
		return "Không có cảnh báo"
	}
	if len(unified.Risk.Reasons) > 0 {
		return unified.Risk.Reasons[0]
	}
	return "Cần kiểm tra thêm thông tin quy hoạch"
}
