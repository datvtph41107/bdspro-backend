package usecase

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	_dto "common/domain/dto"
	"tqd/internal/config"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
	"tqd/internal/modules/resolver/layer_resolver"
	"tqd/internal/modules/resolver/layer_resolver/builder"
	"tqd/internal/modules/resolver/layer_resolver/types"
)

type ParcelUsecase interface {
	// Preview & Search
	GetParcelInfoByID(ctx context.Context, parcelID uint64) (*qh_domain.ParcelInfoResponse, error)
	SearchParcelsByText(ctx context.Context, text string, page, limit int) ([]*qh_domain.ParcelInfoResponse, int64, error)

	// Basic
	GetByIDs(ctx context.Context, ids []uint64) ([]qh_domain.Parcel, error)
	CreateBatch(ctx context.Context, req *dto.ParcelList) (*dto.ParcelList, error)
	FindByLocation(ctx context.Context, lat, lng float64) (*qh_domain.Parcel, error)
	FindParcelsByPolygon(ctx context.Context, points []repo.Point, intersect bool, maxAreaKm2 float64, page, pageSize int) ([]qh_domain.Parcel, int64, error)

	// Detail VIP
	GetParcelInfo(ctx context.Context, parcelID uint64) (*dto.ParcelInfoResponse, error)
	GetParcelLayers(ctx context.Context, parcelID uint64) (*dto.ParcelLayersResponse, error)
	GetParcelLayersV2(ctx context.Context, req *dto.GetParcelLayersV2Request) (*types.ResolveResult, error)
}

type parcelUsecase struct {
	parcelRepo       repo.IParcelRepo
	layerLegalRepo   repo.ILayerLegalRepo
	layerResolver    layer_resolver.ILayerResolver
	candidateBuilder *builder.CandidateBuilder
}

func NewParcelUsecase(
	parcelRepo repo.IParcelRepo,
	layerLegalRepo repo.ILayerLegalRepo,
	layerResolver layer_resolver.ILayerResolver,

) ParcelUsecase {
	return &parcelUsecase{
		parcelRepo:       parcelRepo,
		layerLegalRepo:   layerLegalRepo,
		layerResolver:    layerResolver,
		candidateBuilder: builder.NewCandidateBuilder(),
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

func (u *parcelUsecase) SearchParcelsByText(ctx context.Context, text string, page, limit int) ([]*qh_domain.ParcelInfoResponse, int64, error) {
	if text == "" {
		return nil, 0, errors.New("search text is required")
	}

	pagable := _dto.Pagable{Page: uint32(page), Size: uint32(limit)}
	pagable.Normalize()

	parcels, total, err := u.parcelRepo.SearchByAddress(ctx, text, pagable)
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

func (u *parcelUsecase) parseDirection(dir string) uint32 {
	directions := map[string]uint32{
		"Bắc":      1,
		"Đông Bắc": 2,
		"Đông":     3,
		"Đông Nam": 4,
		"Nam":      5,
		"Tây Nam":  6,
		"Tây":      7,
		"Tây Bắc":  8,
	}
	if val, ok := directions[dir]; ok {
		return val
	}
	return 0
}

func (u *parcelUsecase) parseShape(shape string) uint32 {
	shapes := map[string]uint32{
		"Vuông vức":     1,
		"Chữ nhật":      2,
		"Tóp hậu":       3,
		"Nở hậu":        4,
		"Hình thang":    5,
		"Hình tam giác": 6,
	}
	if val, ok := shapes[shape]; ok {
		return val
	}
	return 0
}

func (u *parcelUsecase) getLandUseDisplayName(code string) string {
	names := map[string]string{
		"ODT": "Đất ở đô thị",
		"ONT": "Đất ở nông thôn",
		"DGT": "Đất giao thông",
		"DCX": "Đất cây xanh",
		"DSH": "Đất sông hồ",
		"DHT": "Đất hạ tầng kỹ thuật",
		"DCC": "Đất công cộng",
		"DLX": "Đất làng xóm",
		"TMD": "Đất thương mại dịch vụ",
		"LNP": "Đất lâm nghiệp",
		"DCN": "Đất công nghiệp",
		"DDT": "Đất di tích",
		"DGD": "Đất giáo dục",
		"DHH": "Đất hỗn hợp",
		"BDX": "Bãi đỗ xe",
		"DAN": "Đất an ninh quốc phòng",
		"DNT": "Đất nghĩa trang",
		"DDL": "Đất du lịch",
		"DDR": "Đất dự trữ",
	}
	if name, ok := names[code]; ok {
		return name
	}
	return code
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

func (u *parcelUsecase) GetParcelLayersV2(ctx context.Context, req *dto.GetParcelLayersV2Request) (*types.ResolveResult, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	rows, err := u.parcelRepo.GetParcelLayerRowsV2(ctx, req.ParcelID)
	if err != nil {
		return nil, fmt.Errorf("get parcel layers failed: %w", err)
	}

	if len(rows) == 0 {
		return &types.ResolveResult{
			CanBuild: false,
			Reason:   "Không có dữ liệu quy hoạch",
		}, nil
	}

	candidates := u.candidateBuilder.BuildBatch(rows)
	if len(candidates) == 0 {
		return &types.ResolveResult{
			CanBuild: false,
			Reason:   "Không có dữ liệu quy hoạch",
		}, nil
	}

	result := u.layerResolver.Resolve(ctx, candidates, req.Mode)

	result.LayerLandUseGroups = buildLayerLandUseGroups(candidates)

	return result, nil
}

func buildLayerLandUseGroups(candidates []*types.LayerCandidate) map[uint64][]*types.LandUseGroupDetail {
	groupMap := make(map[uint64]map[string]*types.LandUseGroupDetail)

	for _, c := range candidates {
		if c.GroupCode == "" || c.GroupCode == "unknown" {
			continue
		}
		if _, ok := groupMap[c.LayerID]; !ok {
			groupMap[c.LayerID] = make(map[string]*types.LandUseGroupDetail)
		}

		key := c.GroupCode + ":" + c.LandUseCode
		existing, exists := groupMap[c.LayerID][key]
		if exists {
			existing.OverlapAreaSqm += c.OverlapAreaSqm
			existing.OverlapPercent += c.OverlapPercent
		} else {
			groupMap[c.LayerID][key] = &types.LandUseGroupDetail{
				Code:           c.GroupCode,
				Name:           c.GroupName,
				Color:          c.GroupColor,
				CanBuild:       c.CanBuild,
				Priority:       c.Priority,
				LandUseCode:    c.LandUseCode,
				LandUseName:    c.LandUseName,
				LandUseColor:   c.LandUseColor,
				OverlapPercent: c.OverlapPercent,
				OverlapAreaSqm: c.OverlapAreaSqm,
			}
		}
	}

	result := make(map[uint64][]*types.LandUseGroupDetail, len(groupMap))
	for layerID, groups := range groupMap {
		list := make([]*types.LandUseGroupDetail, 0, len(groups))
		for _, g := range groups {
			list = append(list, g)
		}
		result[layerID] = list
	}
	return result
}

func (u *parcelUsecase) GetParcelLayers(ctx context.Context, parcelID uint64) (*dto.ParcelLayersResponse, error) {
	start := time.Now()

	if parcelID == 0 {
		return nil, fmt.Errorf("parcel_id is required")
	}

	rows, err := u.parcelRepo.GetParcelLayerRows(ctx, parcelID)
	if err != nil {
		return nil, fmt.Errorf("get parcel layers failed: %w", err)
	}

	parcelArea := 0.0
	if len(rows) > 0 {
		parcelArea = rows[0].ParcelAreaSqm
	} else {
		info, infoErr := u.parcelRepo.GetParcelInfo(ctx, parcelID)
		if infoErr != nil {
			return nil, fmt.Errorf("get parcel info for empty layers failed: %w", infoErr)
		}
		if info == nil || info.ParcelID == 0 {
			return nil, fmt.Errorf("parcel %d not found", parcelID)
		}
		parcelArea = info.AreaSqm
	}

	// THÊM MỚI: Lấy legal docs cho tất cả layer IDs
	layerIDs := extractLayerIDs(rows)
	legalDocsMap, err := u.layerLegalRepo.GetByLayerIDs(ctx, layerIDs)
	if err != nil {
		// Log error nhưng không fail request
		legalDocsMap = make(map[uint64][]*qh_domain.QHLayerLegal)
	}

	resp := buildParcelLayersResponse(parcelID, parcelArea, rows, legalDocsMap)
	resp.Metadata = &dto.ResponseMeta{
		DataSource:     "postgis real-time",
		QueriedAt:      time.Now().UTC().Format(time.RFC3339),
		ResponseTimeMs: time.Since(start).Milliseconds(),
		Cached:         false,
	}

	return resp, nil
}

// Helper function
func extractLayerIDs(rows []dto.ParcelLayerRow) []uint64 {
	if len(rows) == 0 {
		return []uint64{}
	}

	layerIDMap := make(map[uint64]bool)
	for _, row := range rows {
		layerIDMap[row.LayerID] = true
	}

	layerIDs := make([]uint64, 0, len(layerIDMap))
	for id := range layerIDMap {
		layerIDs = append(layerIDs, id)
	}

	return layerIDs
}

func buildParcelLayersResponse(parcelID uint64,
	parcelArea float64,
	rows []dto.ParcelLayerRow,
	legalDocsMap map[uint64][]*qh_domain.QHLayerLegal,
) *dto.ParcelLayersResponse {
	layerMap := make(map[uint64]*dto.LayerImpact, 8)
	labelSeen := make(map[string]struct{}, 32)
	regionSeen := make(map[uint64]struct{}, 64)

	var totalOverlapRaw float64

	for _, r := range rows {
		layer, ok := layerMap[r.LayerID]
		if !ok {
			var legalDocs []*dto.LayerLegalDTO
			if legals, exists := legalDocsMap[r.LayerID]; exists {
				legalDocs = make([]*dto.LayerLegalDTO, 0, len(legals))
				for _, legal := range legals {
					legalDocs = append(legalDocs, &dto.LayerLegalDTO{
						ID:       legal.ID,
						LayerID:  legal.LayerID,
						Name:     legal.Name,
						FileUrl:  legal.FileURL,
						FileType: legal.FileType,
					})
				}
			}

			var authority *dto.AuthorityIssuingDTO
			if r.AuthorityIssuringID > 0 {
				authority = &dto.AuthorityIssuingDTO{
					ID:          r.AuthorityIssuringID,
					Name:        r.AuthorityIssuringName,
					Code:        r.AuthorityIssuringCode,
					Description: r.AuthorityIssuringDesc,
				}
			}

			layer = &dto.LayerImpact{
				ID:               r.LayerID,
				Name:             r.LayerName,
				DisplayName:      r.LayerDisplayName,
				Type:             r.LayerType,
				EffectiveDate:    r.LayerEffectiveDate,
				LegalDoc:         r.LayerLegalDoc,
				LegalDocs:        legalDocs,
				AuthorityIssuing: authority,
				Labels:           []*dto.LabelImpact{},
			}
			layerMap[r.LayerID] = layer
		}

		labelKey := buildLabelKey(r.LayerID, r.RegionLandUseCode)

		var label *dto.LabelImpact
		for _, existing := range layer.Labels {
			if buildLabelKey(r.LayerID, existing.Name) == labelKey {
				label = existing
				break
			}
		}

		if label == nil {
			label = &dto.LabelImpact{
				ID:          stableLabelID(r.LayerID, r.RegionLandUseCode),
				Name:        r.RegionLandUseCode,
				DisplayName: r.RegionLandUseName,
				Color:       r.RegionLandUseColor,
				FontSize:    12,
				FontWeight:  "normal",
				MinZoom:     10,
				MaxZoom:     22,
				Regions:     []*dto.RegionImpact{},
			}
			layer.Labels = append(layer.Labels, label)
			labelSeen[labelKey] = struct{}{}
		}

		region := &dto.RegionImpact{
			ID:             r.RegionID,
			Name:           r.RegionName,
			DisplayName:    r.RegionDisplayName,
			LandUseCode:    r.RegionLandUseCode,
			LandUseName:    r.RegionLandUseName,
			LandUseGroup:   r.RegionLandUseGroup,
			ColorRGB:       regionRGB(r),
			OverlapAreaSqm: r.OverlapAreaSqm,
			OverlapPct:     r.OverlapPct,
			CenterLat:      r.CenterLat,
			CenterLng:      r.CenterLng,
			Geometry:       r.Geometry,
		}

		label.Regions = append(label.Regions, region)

		label.TotalArea += r.OverlapAreaSqm
		layer.TotalArea += r.OverlapAreaSqm
		totalOverlapRaw += r.OverlapAreaSqm

		regionSeen[r.RegionID] = struct{}{}
	}

	layers := make([]*dto.LayerImpact, 0, len(layerMap))
	for _, layer := range layerMap {
		if parcelArea > 0 {
			layer.TotalPct = minFloat((layer.TotalArea/parcelArea)*100, 100)
		}
		layer.Severity = computeLayerSeverity(layer)
		sortLabelRegions(layer, parcelArea)
		layers = append(layers, layer)
	}

	sort.Slice(layers, func(i, j int) bool {
		if layers[i].TotalArea == layers[j].TotalArea {
			return layers[i].ID < layers[j].ID
		}
		return layers[i].TotalArea > layers[j].TotalArea
	})

	// ... rest of the existing code for stats, risk, summary remains the same ...
	totalOverlapArea := totalOverlapRaw
	totalOverlapPct := 0.0
	if parcelArea > 0 {
		totalOverlapPct = minFloat((totalOverlapRaw/parcelArea)*100, 100)
	}

	layersAffected := uint32(len(layers))
	labelsAffected := uint32(len(labelSeen))
	regionsAffected := uint32(len(regionSeen))

	hasConflict := parcelArea > 0 && totalOverlapRaw > parcelArea
	stats := &dto.ImpactStats{
		TotalOverlapArea:     totalOverlapArea,
		TotalOverlapPct:      totalOverlapPct,
		LayersAffected:       layersAffected,
		LabelsAffected:       labelsAffected,
		RegionsAffected:      regionsAffected,
		HasConflict:          hasConflict,
		SpatialDataAvailable: len(rows) > 0,
		Warning:              buildWarning(totalOverlapPct, hasConflict, layersAffected),
	}

	risk := calculateRisk(rows, parcelArea, stats)
	summary := buildSummary(layers, stats, risk)

	return &dto.ParcelLayersResponse{
		ParcelID:      parcelID,
		ParcelAreaSqm: parcelArea,
		Summary:       summary,
		Layers:        layers,
		ConcludeLabel: pickConcludeLabel(layers),
		Statistics:    stats,
		Risk:          risk,
	}
}

// pickConcludeLabel chọn region có diện tích giao (overlap) với thửa lớn nhất.
func pickConcludeLabel(layers []*dto.LayerImpact) *dto.RegionImpact {
	var best *dto.RegionImpact
	for _, layer := range layers {
		if layer == nil {
			continue
		}
		for _, label := range layer.Labels {
			if label == nil {
				continue
			}
			for _, reg := range label.Regions {
				if reg == nil {
					continue
				}
				if best == nil ||
					reg.OverlapAreaSqm > best.OverlapAreaSqm ||
					(reg.OverlapAreaSqm == best.OverlapAreaSqm && reg.ID < best.ID) {
					best = reg
				}
			}
		}
	}
	return best
}

func buildLabelKey(layerID uint64, landUseCode string) string {
	return fmt.Sprintf("%d:%s", layerID, strings.TrimSpace(landUseCode))
}

func stableLabelID(layerID uint64, landUseCode string) uint64 {
	key := buildLabelKey(layerID, landUseCode)
	return uint64(crc32.ChecksumIEEE([]byte(key)))
}

func sortLabelRegions(layer *dto.LayerImpact, parcelArea float64) {
	for _, label := range layer.Labels {
		if parcelArea > 0 {
			label.TotalPct = minFloat((label.TotalArea/parcelArea)*100, 100)
		}
		sort.Slice(label.Regions, func(i, j int) bool {
			if label.Regions[i].OverlapAreaSqm == label.Regions[j].OverlapAreaSqm {
				return label.Regions[i].ID < label.Regions[j].ID
			}
			return label.Regions[i].OverlapAreaSqm > label.Regions[j].OverlapAreaSqm
		})
	}

	sort.Slice(layer.Labels, func(i, j int) bool {
		if layer.Labels[i].TotalArea == layer.Labels[j].TotalArea {
			return layer.Labels[i].ID < layer.Labels[j].ID
		}
		return layer.Labels[i].TotalArea > layer.Labels[j].TotalArea
	})
}

func buildSummary(layers []*dto.LayerImpact, stats *dto.ImpactStats, risk *dto.RiskAssessment) *dto.ParcelSummary {
	summary := &dto.ParcelSummary{
		DominantLandUse:         "",
		DominantPercentage:      0,
		TotalAffectedPercentage: 0,
		TopLayerName:            "",
		Buildable:               true,
	}

	if stats != nil {
		summary.TotalAffectedPercentage = stats.TotalOverlapPct
	}

	var maxLabelPct float64
	var maxLayerPct float64

	for _, layer := range layers {
		if layer.TotalPct > maxLayerPct {
			maxLayerPct = layer.TotalPct
			summary.TopLayerName = layer.DisplayName
		}
		for _, label := range layer.Labels {
			if label.TotalPct > maxLabelPct {
				maxLabelPct = label.TotalPct
				summary.DominantLandUse = label.DisplayName
				summary.DominantPercentage = label.TotalPct
			}
		}
	}

	if risk != nil {
		summary.Buildable = risk.CanBuild
	}

	return summary
}

func buildWarning(totalPct float64, hasConflict bool, layersAffected uint32) string {
	switch {
	case totalPct <= 0:
		return ""
	case hasConflict:
		return "Phát hiện chồng lấn quy hoạch vượt quá diện tích thửa đất"
	case layersAffected > 1:
		return "Thửa đất đang chịu ảnh hưởng từ nhiều layer quy hoạch"
	default:
		return "Thửa đất có ảnh hưởng quy hoạch"
	}
}

func calculateRisk(rows []dto.ParcelLayerRow, parcelArea float64, stats *dto.ImpactStats) *dto.RiskAssessment {
	risk := &dto.RiskAssessment{
		Level:          "LOW",
		Score:          0,
		CanBuild:       true,
		Reasons:        []string{},
		Recommendation: "Rủi ro thấp, có thể xem xét tiếp các điều kiện khác.",
	}

	if parcelArea <= 0 || len(rows) == 0 {
		risk.Reasons = append(risk.Reasons, "Không phát hiện ảnh hưởng quy hoạch")
		return risk
	}

	score := 0.0
	reasons := make([]string, 0, len(rows))

	for _, row := range rows {
		pct := (row.OverlapAreaSqm / parcelArea) * 100

		switch row.RegionLandUseCode {
		case "DGT", "DHT":
			score += pct * 2.0
			reasons = append(reasons, fmt.Sprintf("%.1f%% diện tích thuộc %s - không được phép xây dựng", pct, row.RegionLandUseName))
		case "DCX", "DSH":
			score += pct * 1.0
			reasons = append(reasons, fmt.Sprintf("%.1f%% diện tích thuộc %s - hạn chế xây dựng", pct, row.RegionLandUseName))
		case "ODT", "ONT":
			score += pct * 0.35
			reasons = append(reasons, fmt.Sprintf("%.1f%% diện tích thuộc %s - phù hợp mục đích ở", pct, row.RegionLandUseName))
		default:
			score += pct * 0.8
			reasons = append(reasons, fmt.Sprintf("%.1f%% diện tích thuộc %s - cần kiểm tra thêm", pct, row.RegionLandUseName))
		}
	}

	if stats != nil && stats.LayersAffected > 1 {
		score += float64(stats.LayersAffected) * 5
	}

	switch {
	case score >= 85:
		risk.Level = "CRITICAL"
		risk.Score = 95
		risk.CanBuild = false
		risk.Recommendation = "Rủi ro rất cao, không nên giao dịch trước khi xác minh quy hoạch."
	case score >= 60:
		risk.Level = "HIGH"
		risk.Score = 85
		risk.CanBuild = false
		risk.Recommendation = "Rủi ro cao, cần kiểm tra kỹ quy hoạch trước khi quyết định."
	case score >= 30:
		risk.Level = "MEDIUM"
		risk.Score = 65
		risk.CanBuild = true
		risk.Recommendation = "Có rủi ro quy hoạch, nên kiểm tra thêm hồ sơ pháp lý."
	default:
		risk.Level = "LOW"
		risk.Score = uint32(score)
		risk.CanBuild = true
		risk.Recommendation = "Rủi ro thấp, có thể tiếp tục xem xét."
	}

	if len(reasons) == 0 {
		reasons = append(reasons, "Không phát hiện ảnh hưởng quy hoạch đáng kể")
	}

	risk.Reasons = reasons
	return risk
}

func computeLayerSeverity(layer *dto.LayerImpact) string {
	if layer == nil {
		return "LOW"
	}

	hasCritical := false
	hasHigh := false
	hasMedium := false

	for _, label := range layer.Labels {
		for _, region := range label.Regions {
			switch region.LandUseCode {
			case "DGT", "DHT":
				hasCritical = true
			case "DCX", "DSH":
				hasHigh = true
			case "ODT", "ONT":
				hasMedium = true
			}
		}
	}

	switch {
	case hasCritical:
		return "CRITICAL"
	case hasHigh:
		return "HIGH"
	case hasMedium && layer.TotalPct > 0:
		return "MEDIUM"
	default:
		return "LOW"
	}
}

func regionRGB(r dto.ParcelLayerRow) *dto.RGB {
	if r.RegionLandUseColor != "" {
		if rgb, ok := hexToRGB(r.RegionLandUseColor); ok {
			return rgb
		}
	}

	return &dto.RGB{
		Red:   r.RegionColorRed,
		Green: r.RegionColorGreen,
		Blue:  r.RegionColorBlue,
	}
}

func hexToRGB(hex string) (*dto.RGB, bool) {
	hex = strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(hex) != 6 {
		return nil, false
	}

	v, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return nil, false
	}

	return &dto.RGB{
		Red:   uint8((v >> 16) & 0xFF),
		Green: uint8((v >> 8) & 0xFF),
		Blue:  uint8(v & 0xFF),
	}, true
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
