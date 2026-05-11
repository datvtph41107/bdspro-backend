// infra/mapper/parcel_mapper.go
package mapper

import (
	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"strconv"
	"strings"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
)

type ParcelMapper struct{}

func NewParcelMapper() *ParcelMapper {
	return &ParcelMapper{}
}

// =====================================================
// EXISTING METHODS
// =====================================================

func (m *ParcelMapper) ProtoToDomains(req *tqdpb.ParcelListDTO) []qh_domain.Parcel {
	if req == nil || len(req.Data) == 0 {
		return []qh_domain.Parcel{}
	}

	parcels := make([]qh_domain.Parcel, 0, len(req.Data))
	for _, p := range req.Data {
		parcels = append(parcels, qh_domain.Parcel{
			Lat: p.Latitude,
			Lng: p.Longitude,
			Geometry: qh_domain.RawGeometry{
				Raw: p.Geometry,
			},
		})
	}
	return parcels
}

func (m *ParcelMapper) DomainToProto(p *qh_domain.Parcel) *tqdpb.Parcel {
	if p == nil {
		return nil
	}

	v := &tqdpb.Parcel{
		Id:        p.ID,
		Latitude:  p.Lat,
		Longitude: p.Lng,
	}

	if p.Geometry.Raw != nil {
		v.Geometry = []byte(p.Geometry.Raw)
	}

	return v
}

func (m *ParcelMapper) ToProtoFromDomainList(parcels *dto.ParcelList) []*tqdpb.Parcel {
	if parcels == nil || len(parcels.Data) == 0 {
		return []*tqdpb.Parcel{}
	}

	protoParcels := make([]*tqdpb.Parcel, 0, len(parcels.Data))
	for _, p := range parcels.Data {
		protoParcels = append(protoParcels, m.DomainToProto(&p))
	}
	return protoParcels
}

func (m *ParcelMapper) ToProtoParcelInfoDetail(resp *dto.ParcelInfoResponse) *tqdpb.ParcelInfoResponse {
	if resp == nil {
		return nil
	}

	return &tqdpb.ParcelInfoResponse{
		ParcelId:     resp.ParcelID,
		Lat:          resp.Lat,
		Lon:          resp.Lon,
		AreaSqm:      resp.AreaSqm,
		Address:      resp.AddressText,
		Geometry:     resp.Geometry,
		MapNumber:    resp.MapNumber,
		LandNumber:   resp.LandNumber,
		PropertyCode: resp.PropertyCode,
		PropertyUuid: resp.PropertyUUID,
		Shape:        resp.ShapeType,
		Direction:    resp.Direction,
		LandUseCode:  resp.LandUseCode,
		LandUseName:  resp.LandUseName,
		LandUseColor: resp.LandUseColor,
		Province:     resp.Province,
		ProvinceCode: resp.ProvinceCode,
		WardCode:     resp.WardCode,
	}
}

// SỬA: Thêm logic xử lý legalDocStr
func (m *ParcelMapper) ToProtoParcelLayers(resp *dto.ParcelLayersResponse) *tqdpb.ParcelLayersResponse {
	if resp == nil {
		return nil
	}

	layers := make([]*tqdpb.Layer, 0, len(resp.Layers))
	for _, layer := range resp.Layers {
		if layer == nil {
			continue
		}

		// Tạo legalDocStr từ LegalDocs
		var legalDocStr string
		if len(layer.LegalDocs) > 0 {
			urls := make([]string, 0, len(layer.LegalDocs))
			for _, legal := range layer.LegalDocs {
				if legal.FileUrl != "" {
					urls = append(urls, legal.FileUrl)
				}
			}
			legalDocStr = strings.Join(urls, ",")
		} else {
			legalDocStr = layer.LegalDoc // Fallback to old field
		}

		labels := make([]*tqdpb.Label, 0, len(layer.Labels))
		for _, label := range layer.Labels {
			if label == nil {
				continue
			}

			regions := make([]*tqdpb.Region, 0, len(label.Regions))
			for _, region := range label.Regions {
				if region == nil {
					continue
				}

				var color *tqdpb.RGBColor
				if region.ColorRGB != nil {
					color = &tqdpb.RGBColor{
						Red:   uint32(region.ColorRGB.Red),
						Green: uint32(region.ColorRGB.Green),
						Blue:  uint32(region.ColorRGB.Blue),
					}
				} else {
					color = &tqdpb.RGBColor{}
				}

				regions = append(regions, &tqdpb.Region{
					Id:             region.ID,
					Name:           region.Name,
					DisplayName:    region.DisplayName,
					LandUseCode:    region.LandUseCode,
					LandUseName:    region.LandUseName,
					LandUseGroup:   region.LandUseGroup,
					Color:          color,
					OverlapArea:    region.OverlapAreaSqm,
					OverlapPercent: region.OverlapPct,
					CenterLat:      region.CenterLat,
					CenterLng:      region.CenterLng,
					Geometry:       region.Geometry,
				})
			}

			labels = append(labels, &tqdpb.Label{
				Id:           label.ID,
				Name:         label.Name,
				DisplayName:  label.DisplayName,
				Color:        label.Color,
				FontSize:     label.FontSize,
				FontWeight:   label.FontWeight,
				MinZoom:      label.MinZoom,
				MaxZoom:      label.MaxZoom,
				Regions:      regions,
				TotalArea:    label.TotalArea,
				TotalPercent: label.TotalPct,
			})
		}

		layers = append(layers, &tqdpb.Layer{
			Id:            layer.ID,
			Name:          layer.Name,
			DisplayName:   layer.DisplayName,
			Type:          layer.Type,
			Severity:      layer.Severity,
			EffectiveDate: layer.EffectiveDate,
			LegalDoc:      legalDocStr, // SỬA: sử dụng biến đã khai báo
			Labels:        labels,
			TotalArea:     layer.TotalArea,
			TotalPercent:  layer.TotalPct,
		})
	}

	summary := &tqdpb.ParcelSummary{}
	if resp.Summary != nil {
		summary = &tqdpb.ParcelSummary{
			DominantLandUse:         resp.Summary.DominantLandUse,
			DominantPercentage:      resp.Summary.DominantPercentage,
			TotalAffectedPercentage: resp.Summary.TotalAffectedPercentage,
			TopLayerName:            resp.Summary.TopLayerName,
			Buildable:               resp.Summary.Buildable,
		}
	}

	stats := &dto.ImpactStats{}
	if resp.Statistics != nil {
		stats = resp.Statistics
	}

	risk := &dto.RiskAssessment{}
	if resp.Risk != nil {
		risk = resp.Risk
	}

	meta := &dto.ResponseMeta{}
	if resp.Metadata != nil {
		meta = resp.Metadata
	}

	return &tqdpb.ParcelLayersResponse{
		ParcelId:             resp.ParcelID,
		ParcelAreaSqm:        resp.ParcelAreaSqm,
		Summary:              summary,
		Layers:               layers,
		TotalOverlapArea:     stats.TotalOverlapArea,
		TotalOverlapPercent:  stats.TotalOverlapPct,
		LayersAffected:       stats.LayersAffected,
		LabelsAffected:       stats.LabelsAffected,
		RegionsAffected:      stats.RegionsAffected,
		HasConflict:          stats.HasConflict,
		Warning:              stats.Warning,
		RiskLevel:            risk.Level,
		RiskScore:            risk.Score,
		CanBuild:             risk.CanBuild,
		RiskReasons:          risk.Reasons,
		Recommendation:       risk.Recommendation,
		DataSource:           meta.DataSource,
		QueriedAt:            meta.QueriedAt,
		ResponseTimeMs:       meta.ResponseTimeMs,
		SpatialDataAvailable: stats.SpatialDataAvailable,
	}
}

// =====================================================
// PARCEL INFO TO PROTO (CHO PREVIEW & SEARCH)
// =====================================================

// ParcelInfoToProto - Chuyển QHParcelInfo -> ParcelInfo proto
func (m *ParcelMapper) ParcelInfoToProto(p *qh_domain.QHParcelInfo) *tqdpb.ParcelInfo {
	if p == nil {
		return nil
	}

	// Convert map_number và land_number từ string sang uint32
	mapNum, _ := strconv.ParseUint(p.MapNumber, 10, 32)
	landNum, _ := strconv.ParseUint(p.LandNumber, 10, 32)

	// Xử lý LandTypeId
	var landTypeId uint64
	if p.LandTypeId != nil {
		landTypeId = *p.LandTypeId
	}

	result := &tqdpb.ParcelInfo{
		Id:         p.ID,
		AddressStr: p.AddressText,
		MapNumber:  uint32(mapNum),
		LandNumber: uint32(landNum),
		Lat:        p.Latitude,
		Lng:        p.Longitude,
		Facade:     uint32(p.Facade),
		TotalArea:  p.TotalAreaSqm,
		Geometry:   p.Geometry,
		LandType: &sharepb.ItemV3Proto{
			Id:   landTypeId,
			Name: m.getLandTypeName(landTypeId),
		},
	}

	if p.Shape != nil {
		result.Shape = &sharepb.ItemV3Proto{
			Id:   p.Shape.ID,
			Name: p.Shape.Name,
		}
	}

	return result
}

// ToProtoParcelInfo - Chuyển ParcelInfoResponse -> ParcelInfo proto
func (m *ParcelMapper) ToProtoParcelInfo(p *qh_domain.ParcelInfoResponse) *tqdpb.ParcelInfo {
	if p == nil {
		return nil
	}

	result := &tqdpb.ParcelInfo{
		Id:         p.ID,
		ParcelId:   p.ParcelID,
		AddressStr: p.AddressStr,
		MapNumber:  p.MapNumber,
		LandNumber: p.LandNumber,
		Lat:        p.Lat,
		Lng:        p.Lng,
		Facade:     p.Facade,
		TotalArea:  p.TotalArea,
		Geometry:   p.Geometry,
		LandType: &sharepb.ItemV3Proto{
			Id:   p.LandTypeId,
			Name: m.getLandTypeName(p.LandTypeId),
		},
	}

	if p.Shape != nil {
		result.Shape = &sharepb.ItemV3Proto{
			Id:   p.Shape.ID,
			Name: p.Shape.Name,
		}
	}

	if len(p.Directions) > 0 {
		result.Directions = make([]*sharepb.ItemV3Proto, 0, len(p.Directions))
		for _, d := range p.Directions {
			result.Directions = append(result.Directions, &sharepb.ItemV3Proto{
				Id:   d.ID,
				Name: d.Name,
			})
		}
	}

	return result
}

// =====================================================
// PROTO BUILDERS
// =====================================================

func (m *ParcelMapper) buildLayersProto(layers []*dto.LayerImpact) []*tqdpb.Layer {
	if len(layers) == 0 {
		return nil
	}

	result := make([]*tqdpb.Layer, 0, len(layers))
	for _, l := range layers {
		// Tạo legalDocStr từ LegalDocs
		var legalDocStr string
		if len(l.LegalDocs) > 0 {
			urls := make([]string, 0, len(l.LegalDocs))
			for _, legal := range l.LegalDocs {
				if legal.FileUrl != "" {
					urls = append(urls, legal.FileUrl)
				}
			}
			legalDocStr = strings.Join(urls, ",")
		} else {
			legalDocStr = l.LegalDoc
		}

		result = append(result, &tqdpb.Layer{
			Id:            l.ID,
			Name:          l.Name,
			DisplayName:   l.DisplayName,
			EffectiveDate: l.EffectiveDate,
			LegalDoc:      legalDocStr,
			Labels:        m.buildLabelsProto(l.Labels),
			TotalArea:     l.TotalArea,
			TotalPercent:  l.TotalPct,
		})
	}
	return result
}

func (m *ParcelMapper) buildLabelsProto(labels []*dto.LabelImpact) []*tqdpb.Label {
	if len(labels) == 0 {
		return nil
	}

	result := make([]*tqdpb.Label, 0, len(labels))
	for _, lb := range labels {
		result = append(result, &tqdpb.Label{
			Id:           lb.ID,
			Name:         lb.Name,
			DisplayName:  lb.DisplayName,
			Color:        lb.Color,
			Regions:      m.buildRegionsProto(lb.Regions),
			TotalArea:    lb.TotalArea,
			TotalPercent: lb.TotalPct,
		})
	}
	return result
}

func (m *ParcelMapper) buildRegionsProto(regions []*dto.RegionImpact) []*tqdpb.Region {
	if len(regions) == 0 {
		return nil
	}

	result := make([]*tqdpb.Region, 0, len(regions))
	for _, r := range regions {
		result = append(result, &tqdpb.Region{
			Id:          r.ID,
			Name:        r.Name,
			LandUseCode: r.LandUseCode,
			LandUseName: r.LandUseName,
			Color: &tqdpb.RGBColor{
				Red:   uint32(r.ColorRGB.Red),
				Green: uint32(r.ColorRGB.Green),
				Blue:  uint32(r.ColorRGB.Blue),
			},
			OverlapArea:    r.OverlapAreaSqm,
			OverlapPercent: r.OverlapPct,
			Geometry:       r.Geometry,
		})
	}
	return result
}

// =====================================================
// HELPER METHODS - NAMING
// =====================================================

func (m *ParcelMapper) getDirectionName(direction uint32) string {
	names := map[uint32]string{
		1: "Bắc", 2: "Đông Bắc", 3: "Đông", 4: "Đông Nam",
		5: "Nam", 6: "Tây Nam", 7: "Tây", 8: "Tây Bắc",
	}
	if name, ok := names[direction]; ok {
		return name
	}
	return "Không xác định"
}

func (m *ParcelMapper) getShapeName(shape uint32) string {
	names := map[uint32]string{
		1: "Vuông vức", 2: "Chữ nhật", 3: "Tóp hậu",
		4: "Nở hậu", 5: "Hình thang", 6: "Hình tam giác",
	}
	if name, ok := names[shape]; ok {
		return name
	}
	return "Hình dạng khác"
}

func (m *ParcelMapper) getLandTypeName(landTypeId uint64) string {
	names := map[uint64]string{
		1: "Đất ở đô thị", 2: "Đất ở nông thôn", 3: "Đất nông nghiệp",
		4: "Đất lâm nghiệp", 5: "Đất thương mại", 6: "Đất dịch vụ",
		7: "Đất công nghiệp", 8: "Đất cây xanh", 9: "Đất giao thông",
		10: "Đất hạ tầng", 11: "Đất công cộng", 12: "Đất hỗn hợp",
	}
	if name, ok := names[landTypeId]; ok {
		return name
	}
	return "Đất khác"
}

// =====================================================
// HELPER METHODS - PARSING
// =====================================================

func (m *ParcelMapper) parseDirectionString(dir string) uint32 {
	directions := map[string]uint32{
		"Bắc": 1, "Đông Bắc": 2, "Đông": 3, "Đông Nam": 4,
		"Nam": 5, "Tây Nam": 6, "Tây": 7, "Tây Bắc": 8,
	}
	if val, ok := directions[dir]; ok {
		return val
	}
	return 0
}

func (m *ParcelMapper) parseShapeString(shape string) uint32 {
	shapes := map[string]uint32{
		"Vuông vức": 1, "Chữ nhật": 2, "Tóp hậu": 3,
		"Nở hậu": 4, "Hình thang": 5, "Hình tam giác": 6,
	}
	if val, ok := shapes[shape]; ok {
		return val
	}
	return 0
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
