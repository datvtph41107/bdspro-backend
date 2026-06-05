package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	_dto "common/domain/dto"
	"tqd/internal/config"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
	"tqd/internal/modules/resolver/layer_resolver/engine/parcel_engine"
	"tqd/internal/modules/resolver/layer_resolver/rule"
	"tqd/internal/modules/resolver/layer_resolver/types"
)

type ParcelUsecase interface {
	// Preview & Search
	GetParcelInfoByID(ctx context.Context, parcelID uint64) (*qh_domain.ParcelInfoResponse, error)
	SearchParcelsByText(ctx context.Context, text string, mapNumber, landNumber string, page, limit int) ([]*qh_domain.ParcelInfoResponse, int64, error)

	// Basic
	GetByIDs(ctx context.Context, ids []uint64) ([]qh_domain.Parcel, error)
	CreateBatch(ctx context.Context, req *dto.ParcelList) (*dto.ParcelList, error)
	FindByLocation(ctx context.Context, lat, lng float64) (*qh_domain.Parcel, error)
	FindParcelsByPolygon(ctx context.Context, points []repo.Point, intersect bool, maxAreaKm2 float64, page, pageSize int) ([]qh_domain.Parcel, int64, error)

	// Parcel metadata (Free)
	GetParcelInfo(ctx context.Context, parcelID uint64) (*dto.ParcelInfoResponse, error)

	// Overview - Tổng quan quy hoạch (Free): layers, regions, geometry, basic risk
	GetParcelOverview(ctx context.Context, parcelID uint64) (*dto.ParcelLayersResponse, error)

	// Assessment - Đánh giá chuyên sâu (VIP): conflicts, landUseGroups, assessment, 3 modes
	GetParcelAssessment(ctx context.Context, req *dto.GetParcelAssessmentRequest) (*types.ResolveResult, error)

	// Timeline - Lịch sử biến động quy hoạch của thửa đất
	GetParcelTimeline(ctx context.Context, parcelID uint64) (*types.ParcelTimeline, error)

	// LayerHistory - Lịch sử phiên bản của 1 layer
	GetLayerHistory(ctx context.Context, layerID uint64) (*types.LayerLineage, error)
}

type parcelUsecase struct {
	parcelRepo      repo.IParcelRepo
	layerLegalRepo  repo.ILayerLegalRepo
	auditRepo       rule.AuditRepository
	layerRepo       rule.LayerRepository
	parcelEngine    *parcel_engine.ParcelEngine
	timelineEngine  *rule.TimelineEngine
	lifecycleEngine *rule.LifecycleEngine
}

func NewParcelUsecase(
	parcelRepo repo.IParcelRepo,
	layerLegalRepo repo.ILayerLegalRepo,
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
		auditRepo:       auditRepo,
		layerRepo:       layerRepo,
		parcelEngine:    engine,
		timelineEngine:  timelineEngine,
		lifecycleEngine: lifecycleEngine,
	}
}

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

	// Convert map_number and land_number to uint32
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
		// LandTypeId: func() uint64 {
		// 	if info.LandTypeId != nil {
		// 		return *info.LandTypeId
		// 	}
		// 	return 0
		// }(),
		Geometry: info.Geometry,
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
			// Direction:  p.Direction,
			Shape:     p.Shape,
			ShapeId:   p.ShapeId,
			Facade:    uint32(p.Facade),
			TotalArea: p.TotalAreaSqm,
			Geometry:  p.Geometry,
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

func (u *parcelUsecase) FindParcelsByPolygon(
	ctx context.Context,
	points []repo.Point,
	intersect bool,
	maxAreaKm2 float64,
	page, pageSize int,
) ([]qh_domain.Parcel, int64, error) {
	if len(points) < config.MinPolygonVertices {
		return nil, 0, fmt.Errorf("polygon must have at least %d points", config.MinPolygonVertices)
	}
	if len(points) > config.MaxPolygonVertices {
		return nil, 0, fmt.Errorf("polygon has %d vertices", len(points))
	}

	if err := validatePolygon(points); err != nil {
		return nil, 0, fmt.Errorf("invalid polygon: %v", err)
	}

	area := calculatePolygonAreaKm2(points)
	allowedArea := config.MaxPolygonAreaKm2
	if maxAreaKm2 > 0 && maxAreaKm2 < allowedArea {
		allowedArea = maxAreaKm2
	}
	if area > allowedArea {
		return nil, 0, fmt.Errorf("polygon area %.2f km² exceeds limit", area)
	}

	limit, offset := normalizePage(page, pageSize)
	var parcels []qh_domain.Parcel
	var total int64
	var err error

	if intersect {
		parcels, total, err = u.parcelRepo.FindIntersectPolygon(ctx, points, limit, offset)
	} else {
		parcels, total, err = u.parcelRepo.FindWithinPolygon(ctx, points, limit, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	return parcels, total, nil
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

// GetParcelAssessment - Đánh giá chuyên sâu (VIP): phân tích conflicts, landUseGroups, assessment với 3 mode
func (u *parcelUsecase) GetParcelAssessment(ctx context.Context, req *dto.GetParcelAssessmentRequest) (*types.ResolveResult, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	rows, err := u.parcelRepo.GetParcelLayerRowsV2(ctx, req.ParcelID)
	if err != nil {
		return nil, fmt.Errorf("get parcel layers failed: %w", err)
	}

	if req.Preset == "active_only" {
		rows = filterActiveLayers(rows)
	}
	if req.AsOfDate != "" {
		rows = filterByAsOfDate(rows, req.AsOfDate)
	}

	unifiedRows := parcel_engine.ToUnifiedRowsFromV2(rows)

	parcelInfo, _ := u.parcelRepo.GetParcelInfo(ctx, req.ParcelID)

	layerIDs := extractLayerIDsV2(rows)
	legalDocsMap, _ := u.layerLegalRepo.GetByLayerIDs(ctx, layerIDs)
	legalDocsDTO := convertLegalDocs(legalDocsMap)

	unified := u.parcelEngine.Process(unifiedRows, parcelInfo, legalDocsDTO, req.Mode)

	return u.unifiedToResolveResult(unified, req.Mode), nil
}

// GetParcelOverview - Tổng quan quy hoạch (Free): layers, regions, geometry, basic risk, conclude_label
func (u *parcelUsecase) GetParcelOverview(ctx context.Context, parcelID uint64) (*dto.ParcelLayersResponse, error) {
	start := time.Now()

	if parcelID == 0 {
		return nil, fmt.Errorf("parcel_id is required")
	}

	rows, err := u.parcelRepo.GetParcelLayerRows(ctx, parcelID)
	if err != nil {
		return nil, fmt.Errorf("get parcel layers failed: %w", err)
	}

	layerIDs := extractLayerIDs(rows)
	legalDocsMap, err := u.layerLegalRepo.GetByLayerIDs(ctx, layerIDs)
	if err != nil {
		log.Printf("[GetParcelOverview] Failed to fetch legal docs for parcel %d: %v", parcelID, err)
		legalDocsMap = make(map[uint64][]*qh_domain.QHLayerLegal)
	}
	legalDocsDTO := convertLegalDocs(legalDocsMap)

	parcelInfo, infoErr := u.parcelRepo.GetParcelInfo(ctx, parcelID)
	if infoErr != nil {
		log.Printf("[GetParcelOverview] Failed to fetch parcel info for parcel %d: %v", parcelID, infoErr)
	}

	unifiedRows := parcel_engine.ToUnifiedRowsFromV1(rows)
	unified := u.parcelEngine.Process(unifiedRows, parcelInfo, legalDocsDTO, "detail")

	resp := u.unifiedToParcelLayersResponse(unified)
	resp.Metadata = &dto.ResponseMeta{
		DataSource:     "postgis real-time",
		QueriedAt:      time.Now().UTC().Format(time.RFC3339),
		ResponseTimeMs: time.Since(start).Milliseconds(),
		Cached:         false,
	}

	return resp, nil
}

// ============================================
// CONVERSION HELPERS: UnifiedResult -> Legacy types
// ============================================

func (u *parcelUsecase) unifiedToParcelLayersResponse(ur *parcel_engine.UnifiedResult) *dto.ParcelLayersResponse {
	resp := &dto.ParcelLayersResponse{
		ParcelID:      ur.ParcelID,
		ParcelAreaSqm: ur.ParcelAreaSqm,
		ParcelInfo:    ur.ParcelInfo,
	}

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
	}

	return resp
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

	// Secondary plans: ưu tiên từ SecondaryPlans (compare/detail mode), fallback từ Layers (V1 detail)
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

	// Layers: convert engine layers → types layers (chỉ cho detail mode)
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
			ID:            l.ID,
			Name:          l.Name,
			DisplayName:   l.DisplayName,
			Type:          l.Type,
			EffectiveDate: l.EffectiveDate,
			TotalArea:     l.TotalArea,
			TotalPct:      l.TotalPct,
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
					ColorRGB:       reg.Color,
					OverlapAreaSqm: reg.OverlapAreaSqm,
					OverlapPct:     reg.OverlapPct,
					CenterLat:      reg.CenterLat,
					CenterLng:      reg.CenterLng,
					Geometry:       reg.Geometry,
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
						ColorRGB:       reg.Color,
						OverlapAreaSqm: reg.OverlapAreaSqm,
						OverlapPct:     reg.OverlapPct,
						CenterLat:      reg.CenterLat,
						CenterLng:      reg.CenterLng,
						Geometry:       reg.Geometry,
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

// ─── Engine → Types layer conversion ───────────────────────────────────────

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
		ID:            el.ID,
		Name:          el.Name,
		DisplayName:   el.DisplayName,
		Type:          el.Type,
		EffectiveDate: el.EffectiveDate,
		TotalArea:     el.TotalArea,
		TotalPct:      el.TotalPct,
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
				OverlapAreaSqm: r.OverlapAreaSqm,
				OverlapPct:     r.OverlapPct,
				Geometry:       r.Geometry,
			}
			if r.Color != nil {
				tr.ColorRed = r.Color.Red
				tr.ColorGreen = r.Color.Green
				tr.ColorBlue = r.Color.Blue
			}
			tl.Regions = append(tl.Regions, tr)
		}
		l.Labels = append(l.Labels, tl)
	}
	return l
}

// GetParcelTimeline - Lịch sử biến động quy hoạch của thửa đất
func (u *parcelUsecase) GetParcelTimeline(ctx context.Context, parcelID uint64) (*types.ParcelTimeline, error) {
	if parcelID == 0 {
		return nil, fmt.Errorf("parcel_id is required")
	}

	if u.timelineEngine != nil {
		return u.timelineEngine.BuildTimeline(ctx, parcelID)
	}

	// Fallback: build from parcel layer rows
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

// GetLayerHistory - Lịch sử phiên bản của 1 layer
func (u *parcelUsecase) GetLayerHistory(ctx context.Context, layerID uint64) (*types.LayerLineage, error) {
	if layerID == 0 {
		return nil, fmt.Errorf("layer_id is required")
	}

	if u.lifecycleEngine != nil {
		return u.lifecycleEngine.GetLineage(ctx, layerID)
	}

	// Fallback: return basic lineage
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
