package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	_dto "common/domain/dto"
	"tqd/internal/config"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
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
	GetParcelDetail(ctx context.Context, parcelID uint64) (*qh_domain.ParcelDetailResponse, error)
}

type parcelUsecase struct {
	parcelRepo      repo.IParcelRepo
	landUseDictRepo repo.LandUseDictRepository
}

func NewParcelUsecase(parcelRepo repo.IParcelRepo, landUseDictRepo repo.LandUseDictRepository) ParcelUsecase {
	return &parcelUsecase{
		parcelRepo:      parcelRepo,
		landUseDictRepo: landUseDictRepo,
	}
}

// =====================================================
// Preview & Search
// =====================================================

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
		Direction:  u.parseDirection(info.Direction),
		Shape:      u.parseShape(info.ShapeType),
		ShapeId:    info.ShapeId,
		Facade:     uint32(info.Facade),
		TotalArea:  info.TotalAreaSqm,
		LandTypeId: func() uint64 {
			if info.LandTypeId != nil {
				return *info.LandTypeId
			}
			return 0
		}(),
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
			AddressStr: p.AddressText,
			MapNumber:  uint32(mapNum),
			LandNumber: uint32(landNum),
			Lat:        p.Latitude,
			Lng:        p.Longitude,
			Direction:  u.parseDirection(p.Direction),
			Shape:      u.parseShape(p.ShapeType),
			ShapeId:    p.ShapeId,
			Facade:     uint32(p.Facade),
			TotalArea:  p.TotalAreaSqm,
			Geometry:   p.Geometry,
		})
	}

	return results, total, nil
}

// =====================================================
// Helper Methods
// =====================================================

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

func (u *parcelUsecase) getLandUseName(lu *qh_domain.LandUseDict) string {
	if lu == nil {
		return "Chưa phân loại"
	}
	return lu.Name
}

func (u *parcelUsecase) getLandUseGroup(lu *qh_domain.LandUseDict) string {
	if lu == nil || lu.GroupName == "" {
		return "Đất khác"
	}
	return lu.GroupName
}

func (u *parcelUsecase) getLandUseColor(lu *qh_domain.LandUseDict) string {
	if lu == nil || lu.Color == "" {
		return "#CCCCCC"
	}
	return lu.Color
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

// =====================================================
// Basic Operations
// =====================================================

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

// =====================================================
// Detail VIP
// =====================================================

func (u *parcelUsecase) GetParcelDetail(ctx context.Context, parcelID uint64) (*qh_domain.ParcelDetailResponse, error) {
	startTime := time.Now()

	if parcelID == 0 {
		return nil, errors.New("parcel_id is required")
	}

	baseInfo, err := u.parcelRepo.GetParcelDetail(ctx, parcelID)
	if err != nil {
		return nil, fmt.Errorf("get parcel detail failed: %w", err)
	}
	if baseInfo == nil || baseInfo.ID == 0 {
		return nil, fmt.Errorf("parcel %d not found", parcelID)
	}

	spatialImpacts, _ := u.parcelRepo.GetSpatialImpacts(ctx, parcelID)
	landUse, _ := u.landUseDictRepo.GetByCode(ctx, baseInfo.ActualLandUseCode)

	response := &qh_domain.ParcelDetailResponse{
		ParcelID:     baseInfo.ID,
		Lat:          baseInfo.Lat,
		Lon:          baseInfo.Lon,
		TotalAreaSqm: baseInfo.TotalAreaSqm,
		AddressText:  baseInfo.AddressText,
		Geometry:     baseInfo.Geometry,
		MapNumber:    baseInfo.MapNumber,
		LandNumber:   baseInfo.LandNumber,
		PropertyCode: baseInfo.PropertyCode,
		PropertyUUID: baseInfo.PropertyUUID,
		ShapeType:    baseInfo.ShapeType,
		Direction:    baseInfo.Direction,
		CurrentLandUse: &qh_domain.LandUseBrief{
			Code:  baseInfo.ActualLandUseCode,
			Name:  u.getLandUseName(landUse),
			Group: u.getLandUseGroup(landUse),
			Color: u.getLandUseColor(landUse),
		},
		Administrative: &qh_domain.AdminBrief{
			Province:     baseInfo.ProvinceName,
			ProvinceCode: baseInfo.ProvinceCode,
			WardCode:     baseInfo.WardCode,
			Address:      baseInfo.AddressText,
		},
	}

	if len(spatialImpacts) > 0 {
		response.Layers = u.buildLayers(spatialImpacts)
		response.Statistics = u.calculateStats(response.Layers, baseInfo.TotalAreaSqm)
		response.Risk = u.assessRisk(response.Layers, baseInfo.ActualLandUseCode, baseInfo.TotalAreaSqm)
		response.Metadata = &qh_domain.ResponseMeta{
			DataSource:     "qh_regions (spatial)",
			QueriedAt:      time.Now().Format(time.RFC3339),
			ResponseTimeMs: time.Since(startTime).Milliseconds(),
			Cached:         false,
		}
	} else {
		response.Statistics = &qh_domain.ImpactStats{SpatialDataAvailable: false}
		response.Risk = &qh_domain.RiskAssessment{
			Level:          "LOW",
			Score:          0,
			CanBuild:       true,
			Reasons:        []string{"✅ Không phát hiện ảnh hưởng quy hoạch"},
			Recommendation: "✅ Rủi ro thấp, có thể xem xét giao dịch.",
		}
		response.Metadata = &qh_domain.ResponseMeta{
			DataSource:     "qh_regions (spatial)",
			QueriedAt:      time.Now().Format(time.RFC3339),
			ResponseTimeMs: time.Since(startTime).Milliseconds(),
			Cached:         false,
		}
	}

	return response, nil
}

// =====================================================
// Build methods for Detail VIP
// =====================================================

func (u *parcelUsecase) buildLayers(impacts []*repo.SpatialImpactRow) []*qh_domain.LayerImpact {
	layerMap := make(map[uint64]*qh_domain.LayerImpact)
	for _, imp := range impacts {
		layer, ok := layerMap[imp.LayerID]
		if !ok {
			layer = &qh_domain.LayerImpact{
				ID:            imp.LayerID,
				Name:          imp.LayerName,
				DisplayName:   imp.LayerDisplayName,
				Type:          imp.LayerType,
				EffectiveDate: imp.EffectiveDate,
				LegalDoc:      imp.LegalDoc,
				Labels:        []*qh_domain.LabelImpact{},
			}
			layerMap[imp.LayerID] = layer
		}

		var label *qh_domain.LabelImpact
		for _, lbl := range layer.Labels {
			if lbl.ID == imp.LabelID {
				label = lbl
				break
			}
		}
		if label == nil {
			label = &qh_domain.LabelImpact{
				ID:          imp.LabelID,
				Name:        imp.LabelName,
				DisplayName: imp.LabelDisplayName,
				Color:       imp.LabelColor,
				FontSize:    imp.LabelFontSize,
				FontWeight:  imp.LabelFontWeight,
				MinZoom:     imp.LabelMinZoom,
				MaxZoom:     imp.LabelMaxZoom,
				Regions:     []*qh_domain.RegionImpact{},
			}
			layer.Labels = append(layer.Labels, label)
		}

		region := &qh_domain.RegionImpact{
			ID:           imp.RegionID,
			Name:         imp.RegionName,
			DisplayName:  imp.RegionDisplayName,
			LandUseCode:  imp.LandUseCode,
			LandUseName:  imp.LandUseName,
			LandUseGroup: imp.LandUseGroup,
			ColorRGB: &qh_domain.RGB{
				Red:   imp.ColorRed,
				Green: imp.ColorGreen,
				Blue:  imp.ColorBlue,
			},
			OverlapAreaSqm: imp.OverlapArea,
			OverlapPct:     imp.OverlapPct,
			Geometry:       imp.Geometry,
		}
		label.Regions = append(label.Regions, region)
		label.TotalArea += imp.OverlapArea
		label.TotalPct += imp.OverlapPct
	}

	layers := make([]*qh_domain.LayerImpact, 0, len(layerMap))
	for _, layer := range layerMap {
		for _, label := range layer.Labels {
			layer.TotalArea += label.TotalArea
			layer.TotalPct += label.TotalPct
			sort.Slice(label.Regions, func(i, j int) bool {
				return label.Regions[i].OverlapPct > label.Regions[j].OverlapPct
			})
		}
		layers = append(layers, layer)
	}
	sort.Slice(layers, func(i, j int) bool {
		return layers[i].TotalArea > layers[j].TotalArea
	})
	return layers
}

func (u *parcelUsecase) calculateStats(layers []*qh_domain.LayerImpact, totalArea float64) *qh_domain.ImpactStats {
	stats := &qh_domain.ImpactStats{SpatialDataAvailable: len(layers) > 0}
	if len(layers) == 0 {
		return stats
	}
	for _, layer := range layers {
		stats.TotalOverlapArea += layer.TotalArea
		stats.LayersAffected++
		for _, label := range layer.Labels {
			stats.LabelsAffected++
			stats.RegionsAffected += uint32(len(label.Regions))
		}
	}
	if totalArea > 0 {
		stats.TotalOverlapPct = (stats.TotalOverlapArea / totalArea) * 100
	}
	stats.HasConflict = stats.TotalOverlapPct > 100
	if stats.HasConflict {
		stats.Warning = "⚠️ PHÁT HIỆN CHỒNG LẤN QUY HOẠCH! Tổng diện tích ảnh hưởng vượt quá diện tích thửa đất."
	}
	return stats
}

func (u *parcelUsecase) assessRisk(layers []*qh_domain.LayerImpact, landUseCode string, totalArea float64) *qh_domain.RiskAssessment {
	risk := &qh_domain.RiskAssessment{
		Level:          "LOW",
		Score:          0,
		CanBuild:       true,
		Reasons:        []string{},
		Recommendation: "✅ Rủi ro thấp, có thể xem xét giao dịch.",
	}
	if len(layers) == 0 {
		risk.Reasons = append(risk.Reasons, "✅ Không phát hiện ảnh hưởng quy hoạch")
		return risk
	}

	score := 0.0
	landUseMap := make(map[string]float64)
	for _, layer := range layers {
		for _, label := range layer.Labels {
			for _, region := range label.Regions {
				landUseMap[region.LandUseCode] += region.OverlapAreaSqm
			}
		}
	}

	for code, area := range landUseMap {
		pct := (area / totalArea) * 100
		switch code {
		case "DGT", "DHT":
			score += pct * 2
			risk.Reasons = append(risk.Reasons, fmt.Sprintf("⚠️ %.1f%% diện tích là %s - không được phép xây dựng", pct, u.getLandUseDisplayName(code)))
		case "DCX", "DSH":
			score += pct
			risk.Reasons = append(risk.Reasons, fmt.Sprintf("⚠️ %.1f%% diện tích là %s - hạn chế xây dựng", pct, u.getLandUseDisplayName(code)))
		case "ODT", "ONT":
			score += pct * 0.5
			risk.Reasons = append(risk.Reasons, fmt.Sprintf("✅ %.1f%% diện tích là %s - được phép xây dựng", pct, u.getLandUseDisplayName(code)))
		default:
			score += pct * 0.8
			risk.Reasons = append(risk.Reasons, fmt.Sprintf("⚠️ %.1f%% diện tích là %s - cần kiểm tra thêm", pct, u.getLandUseDisplayName(code)))
		}
	}

	if score > 150 {
		risk.Level = "CRITICAL"
		risk.Score = 95
		risk.CanBuild = false
		risk.Recommendation = "⚠️ RỦI RO CAO: KHÔNG NÊN GIAO DỊCH. Liên hệ cơ quan quy hoạch."
	} else if score > 80 {
		risk.Level = "HIGH"
		risk.Score = 85
		risk.CanBuild = false
		risk.Recommendation = "⚠️ RỦI RO: Cần kiểm tra kỹ quy hoạch trước khi quyết định."
	} else if score > 30 {
		risk.Level = "MEDIUM"
		risk.Score = 65
		risk.CanBuild = true
		risk.Recommendation = "⚠️ Có rủi ro quy hoạch, cần tìm hiểu thêm."
	}
	return risk
}

// =====================================================
// Helper functions
// =====================================================

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
