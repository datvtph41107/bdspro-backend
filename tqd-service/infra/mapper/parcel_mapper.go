package mapper

import (
	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"strconv"
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

func (m *ParcelMapper) ToProtoParcelDetail(resp *qh_domain.ParcelDetailResponse) *tqdpb.ParcelDetailResponse {
	if resp == nil {
		return nil
	}

	// Build layers
	layers := make([]*tqdpb.Layer, 0, len(resp.Layers))
	for _, l := range resp.Layers {
		labels := make([]*tqdpb.Label, 0, len(l.Labels))
		for _, lb := range l.Labels {
			regions := make([]*tqdpb.Region, 0, len(lb.Regions))
			for _, r := range lb.Regions {
				regions = append(regions, &tqdpb.Region{
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
			labels = append(labels, &tqdpb.Label{
				Id:           lb.ID,
				Name:         lb.Name,
				DisplayName:  lb.DisplayName,
				Color:        lb.Color,
				Regions:      regions,
				TotalArea:    lb.TotalArea,
				TotalPercent: lb.TotalPct,
			})
		}
		layers = append(layers, &tqdpb.Layer{
			Id:            l.ID,
			Name:          l.Name,
			DisplayName:   l.DisplayName,
			EffectiveDate: l.EffectiveDate,
			LegalDoc:      l.LegalDoc,
			Labels:        labels,
			TotalArea:     l.TotalArea,
			TotalPercent:  l.TotalPct,
		})
	}

	// TRẢ VỀ TRỰC TIẾP, KHÔNG QUA NESTED MESSAGES
	return &tqdpb.ParcelDetailResponse{
		// Parcel cơ bản
		ParcelId: resp.ParcelID,
		Lat:      resp.Lat,
		Lon:      resp.Lon,
		AreaSqm:  resp.TotalAreaSqm,
		Address:  resp.AddressText,
		Geometry: resp.Geometry,

		// Thông tin thửa đất
		MapNumber:    resp.MapNumber,
		LandNumber:   resp.LandNumber,
		PropertyCode: resp.PropertyCode,
		PropertyUuid: resp.PropertyUUID,
		Shape:        resp.ShapeType,
		Direction:    resp.Direction,

		// Loại đất hiện tại
		LandUseCode:  resp.CurrentLandUse.Code,
		LandUseName:  resp.CurrentLandUse.Name,
		LandUseColor: resp.CurrentLandUse.Color,

		// Hành chính
		Province:     resp.Administrative.Province,
		ProvinceCode: resp.Administrative.ProvinceCode,
		WardCode:     resp.Administrative.WardCode,

		// Layers
		Layers: layers,

		// Thống kê
		TotalOverlapArea:    resp.Statistics.TotalOverlapArea,
		TotalOverlapPercent: resp.Statistics.TotalOverlapPct,
		LayersAffected:      resp.Statistics.LayersAffected,
		LabelsAffected:      resp.Statistics.LabelsAffected,
		RegionsAffected:     resp.Statistics.RegionsAffected,
		HasConflict:         resp.Statistics.HasConflict,
		Warning:             resp.Statistics.Warning,

		// Rủi ro
		RiskLevel:      resp.Risk.Level,
		RiskScore:      resp.Risk.Score,
		CanBuild:       resp.Risk.CanBuild,
		RiskReasons:    resp.Risk.Reasons,
		Recommendation: resp.Risk.Recommendation,

		// Meta
		DataSource:     resp.Metadata.DataSource,
		QueriedAt:      resp.Metadata.QueriedAt,
		ResponseTimeMs: resp.Metadata.ResponseTimeMs,
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

	// Parse direction từ string sang uint32
	directionVal := m.parseDirectionString(p.Direction)
	shapeVal := m.parseShapeString(p.ShapeType)

	// Xử lý LandTypeId
	var landTypeId uint64
	if p.LandTypeId != nil {
		landTypeId = *p.LandTypeId
	}

	// Xử lý ShapeId (convert uint32 sang uint64)
	shapeId := uint64(p.ShapeId)

	return &tqdpb.ParcelInfo{
		Id:         p.ID,
		AddressStr: p.AddressText,
		MapNumber:  uint32(mapNum),
		LandNumber: uint32(landNum),
		Lat:        p.Latitude,
		Lng:        p.Longitude,
		Facade:     uint32(p.Facade),
		TotalArea:  p.TotalAreaSqm,
		Geometry:   p.Geometry,
		Direction: &sharepb.ItemV3Proto{
			Id:   uint64(directionVal),
			Name: m.getDirectionName(directionVal),
		},
		Shape: &sharepb.ItemV3Proto{
			Id:   shapeId,
			Name: m.getShapeName(shapeVal),
		},
		LandType: &sharepb.ItemV3Proto{
			Id:   landTypeId,
			Name: m.getLandTypeName(landTypeId),
		},
	}
}

// ToProtoParcelInfo - Chuyển ParcelInfoResponse -> ParcelInfo proto
func (m *ParcelMapper) ToProtoParcelInfo(p *qh_domain.ParcelInfoResponse) *tqdpb.ParcelInfo {
	if p == nil {
		return nil
	}

	return &tqdpb.ParcelInfo{
		Id:         p.ID,
		AddressStr: p.AddressStr,
		MapNumber:  p.MapNumber,
		LandNumber: p.LandNumber,
		Lat:        p.Lat,
		Lng:        p.Lng,
		Facade:     p.Facade,
		TotalArea:  p.TotalArea,
		Geometry:   p.Geometry,
		Direction: &sharepb.ItemV3Proto{
			Id:   uint64(p.Direction),
			Name: m.getDirectionName(p.Direction),
		},
		Shape: &sharepb.ItemV3Proto{
			Id:   uint64(p.Shape), // SỬA: convert uint32 -> uint64
			Name: m.getShapeName(p.Shape),
		},
		LandType: &sharepb.ItemV3Proto{
			Id:   p.LandTypeId,
			Name: m.getLandTypeName(p.LandTypeId),
		},
	}
}

// =====================================================
// PROTO BUILDERS
// =====================================================

func (m *ParcelMapper) buildLayersProto(layers []*qh_domain.LayerImpact) []*tqdpb.Layer {
	if len(layers) == 0 {
		return nil
	}

	result := make([]*tqdpb.Layer, 0, len(layers))
	for _, l := range layers {
		result = append(result, &tqdpb.Layer{
			Id:            l.ID,
			Name:          l.Name,
			DisplayName:   l.DisplayName,
			EffectiveDate: l.EffectiveDate,
			LegalDoc:      l.LegalDoc,
			Labels:        m.buildLabelsProto(l.Labels),
			TotalArea:     l.TotalArea,
			TotalPercent:  l.TotalPct,
		})
	}
	return result
}

func (m *ParcelMapper) buildLabelsProto(labels []*qh_domain.LabelImpact) []*tqdpb.Label {
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

func (m *ParcelMapper) buildRegionsProto(regions []*qh_domain.RegionImpact) []*tqdpb.Region {
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
