// infra/mapper/parcel_mapper.go
package mapper

import (
	"fmt"
	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"strconv"
	"strings"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/modules/resolver/layer_resolver/types"
)

type ParcelMapper struct {
	LandTypeNames         map[uint64]string
	ConflictSummarySuffix string
}

const (
	WorkspaceEntityTypeParcel uint32 = 1
	WorkspaceEntityTypeRegion uint32 = 3
)

func NewParcelMapper() *ParcelMapper {
	landTypes := make(map[uint64]string, len(types.DefaultLandTypeNames))
	for k, v := range types.DefaultLandTypeNames {
		landTypes[k] = v
	}
	return &ParcelMapper{
		LandTypeNames:         landTypes,
		ConflictSummarySuffix: types.DefaultConflictSummarySuffix,
	}
}

func (m *ParcelMapper) ResolveResultToProtoAssessment(result *types.ResolveResult, parcelID uint64, mode string) *tqdpb.ParcelAssessmentResponse {
	resp := &tqdpb.ParcelAssessmentResponse{
		ParcelId:         parcelID,
		HasConflict:      result.HasConflict,
		RiskScore:        result.RiskScore,
		RiskLevel:        result.RiskLevel,
		RiskReasons:      result.RiskReasons,
		ResolutionStatus: result.ResolutionStatus,
		Metadata: &tqdpb.ResponseMetadataInfo{
			ResponseTimeMs: result.ProcessingTimeMs,
		},
	}

	// InfoConflicts
	for _, c := range result.InfoConflicts {
		resp.InfoConflicts = append(resp.InfoConflicts, &tqdpb.ConflictDetail{
			Type:            c.Type,
			Severity:        string(c.Severity),
			Description:     c.Description,
			Recommendation:  c.Recommendation,
			ActionRequired:  c.ActionRequired,
			AffectedPercent: c.AffectedPercent,
		})
	}

	// CompareResult
	if result.CompareResult != nil {
		resp.CompareResult = &tqdpb.CompareResult{
			TotalLayers: int32(result.CompareResult.TotalLayers),
		}
		for _, c := range result.CompareResult.Changes {
			resp.CompareResult.Changes = append(resp.CompareResult.Changes, &tqdpb.CompareChange{
				LayerId:    c.LayerID,
				LayerName:  c.LayerName,
				ChangeType: c.ChangeType,
				OldValue:   c.OldValue,
				NewValue:   c.NewValue,
			})
		}
	}

	// HistoricalRisk
	if result.HistoricalRisk != nil {
		resp.HistoricalRisk = &tqdpb.HistoricalRisk{
			Trend:          result.HistoricalRisk.Trend,
			TrendDetail:    result.HistoricalRisk.TrendDetail,
			ActiveLayers:   int32(result.HistoricalRisk.ActiveLayers),
			ExpiringLayers: int32(result.HistoricalRisk.ExpiringLayers),
			AnalyzedAt:     result.HistoricalRisk.AnalyzedAt,
		}
	}

	// Assessment: chỉ build cho detail mode (popup/compare không cần assessment chi tiết)
	if mode == "detail" || mode == "" {
		resp.Assessment = &tqdpb.AssessmentResponse{
			Buildability: &tqdpb.BuildabilityResponse{
				CanBuild:       result.CanBuild,
				Status:         result.BuildStatus,
				Reason:         result.Reason,
				Recommendation: result.Recommendation,
			},
			Summary: m.buildAssessmentSummary(result),
			Verdict: result.RiskLevel,
		}
	}

	if result.Primary != nil {
		resp.CurrentPlan = &tqdpb.PlanDetailInfo{
			LayerId:       result.Primary.LayerID,
			LayerName:     result.Primary.LayerName,
			DisplayName:   result.Primary.LayerDisplayName,
			LandUseCode:   result.Primary.LandUseCode,
			LandUseName:   result.Primary.LandUseName,
			LandUseColor:  result.Primary.LandUseColor,
			GroupCode:     result.Primary.GroupCode,
			GroupName:     result.Primary.GroupName,
			CanBuild:      result.Primary.CanBuild,
			ImpactPercent: result.Primary.OverlapPercent,
			ImpactAreaSqm: result.Primary.OverlapAreaSqm,
			LandUseGroups: mapLandUseGroups(result.LayerLandUseGroups, result.Primary.LayerID),
		}
		resp.ParcelAreaSqm = result.Primary.ParcelAreaSqm
	}

	// OtherPlans: bỏ qua cho popup mode
	if mode != "popup" {
		for _, sec := range result.Secondary {
			resp.OtherPlans = append(resp.OtherPlans, &tqdpb.PlanDetailInfo{
				LayerId:       sec.LayerID,
				LayerName:     sec.LayerName,
				DisplayName:   sec.LayerDisplayName,
				LandUseCode:   sec.LandUseCode,
				LandUseName:   sec.LandUseName,
				LandUseColor:  sec.LandUseColor,
				GroupCode:     sec.GroupCode,
				GroupName:     sec.GroupName,
				CanBuild:      sec.CanBuild,
				ImpactPercent: sec.OverlapPercent,
				ImpactAreaSqm: sec.OverlapAreaSqm,
				LandUseGroups: mapLandUseGroups(result.LayerLandUseGroups, sec.LayerID),
			})
		}
	}

	// Conflicts: bỏ qua cho popup mode
	if mode != "popup" {
		for _, c := range result.CriticalConflicts {
			resp.Conflicts = append(resp.Conflicts, &tqdpb.ConflictDetail{
				Type:            c.Type,
				Severity:        string(c.Severity),
				Description:     c.Description,
				Recommendation:  c.Recommendation,
				ActionRequired:  c.ActionRequired,
				AffectedPercent: c.AffectedPercent,
			})
		}
	}

	// Layers: chỉ map cho detail mode
	if mode == "detail" || mode == "" {
		for _, l := range result.Layers {
			resp.Layers = append(resp.Layers, mapTypeLayerToLayerDetail(l))
		}
	}

	return resp
}

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

	if p.SeoID != nil {
		v.SeoId = *p.SeoID
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
		SeoId:        resp.SeoID,
		IsSeo:        resp.IsSeo,
	}
}

func (m *ParcelMapper) RegionInfoToProto(v *dto.RegionInfoResponse) *tqdpb.RegionInfoResponse {
	if v == nil {
		return nil
	}

	return &tqdpb.RegionInfoResponse{
		RegionId:  v.RegionID,
		LayerId:   v.LayerID,
		LabelId:   v.LabelID,
		LandUseId: v.LandUseID,
		LegendId:  v.LegendID,

		Name:        v.Name,
		DisplayName: v.DisplayName,
		Description: v.Description,

		LayerName:        v.LayerName,
		LayerDisplayName: v.LayerDisplayName,

		LandUseCode:  v.LandUseCode,
		LandUseName:  v.LandUseName,
		LandUseColor: v.LandUseColor,
		CanBuild:     v.CanBuild,

		LabelName:  v.LabelName,
		LabelColor: v.LabelColor,

		LegendColor:  v.LegendColor,
		LegendType:   v.LegendType,
		GeometryType: v.GeometryType,

		LegalDoc:     v.LegalDoc,
		PlanningName: v.PlanningName,

		AreaSqm:    v.AreaSqm,
		AreaHa:     v.AreaHa,
		PerimeterM: v.PerimeterM,

		CenterLat: v.CenterLat,
		CenterLon: v.CenterLon,

		Bounds: &tqdpb.RegionBounds{
			MinLon: v.MinLon,
			MinLat: v.MinLat,
			MaxLon: v.MaxLon,
			MaxLat: v.MaxLat,
		},

		GeoJson: v.GeoJSON,

		Province:     v.Province,
		ProvinceCode: v.ProvinceCode,
		WardCode:     v.WardCode,
	}
}

func (m *ParcelMapper) RegionInfoListToProto(items []dto.RegionInfoResponse) []*tqdpb.RegionInfoResponse {
	if len(items) == 0 {
		return []*tqdpb.RegionInfoResponse{}
	}

	result := make([]*tqdpb.RegionInfoResponse, 0, len(items))

	for i := range items {
		result = append(result, m.RegionInfoToProto(&items[i]))
	}

	return result
}

func (m *ParcelMapper) DomainListToProto(parcels []qh_domain.Parcel) []*tqdpb.Parcel {
	if len(parcels) == 0 {
		return []*tqdpb.Parcel{}
	}

	result := make([]*tqdpb.Parcel, 0, len(parcels))

	for i := range parcels {
		result = append(result, m.DomainToProto(&parcels[i]))
	}

	return result
}

func polygonLandUseStatsToProto(items []dto.PolygonLandUseStat) []*tqdpb.PolygonLandUseStat {
	if len(items) == 0 {
		return []*tqdpb.PolygonLandUseStat{}
	}

	result := make([]*tqdpb.PolygonLandUseStat, 0, len(items))

	for i := range items {
		item := items[i]
		result = append(result, &tqdpb.PolygonLandUseStat{
			LandUseId:    item.LandUseID,
			LandUseCode:  item.LandUseCode,
			LandUseName:  item.LandUseName,
			LandUseColor: item.LandUseColor,
			CanBuild:     item.CanBuild,
			RegionCount:  item.RegionCount,
			AreaSqm:      item.AreaSqm,
			AreaPct:      item.AreaPct,
		})
	}

	return result
}

func polygonLayerStatsToProto(items []dto.PolygonLayerStat) []*tqdpb.PolygonLayerStat {
	if len(items) == 0 {
		return []*tqdpb.PolygonLayerStat{}
	}

	result := make([]*tqdpb.PolygonLayerStat, 0, len(items))

	for i := range items {
		item := items[i]
		result = append(result, &tqdpb.PolygonLayerStat{
			LayerId:          item.LayerID,
			LayerName:        item.LayerName,
			LayerDisplayName: item.LayerDisplayName,
			RegionCount:      item.RegionCount,
			AreaSqm:          item.AreaSqm,
			AreaPct:          item.AreaPct,
		})
	}

	return result
}

func polygonBuildabilityStatsToProto(items []dto.PolygonBuildabilityStat) []*tqdpb.PolygonBuildabilityStat {
	if len(items) == 0 {
		return []*tqdpb.PolygonBuildabilityStat{}
	}

	result := make([]*tqdpb.PolygonBuildabilityStat, 0, len(items))

	for i := range items {
		item := items[i]
		result = append(result, &tqdpb.PolygonBuildabilityStat{
			CanBuild:    item.CanBuild,
			RegionCount: item.RegionCount,
			AreaSqm:     item.AreaSqm,
			AreaPct:     item.AreaPct,
		})
	}

	return result
}

func polygonAnalysisSummaryToProto(summary *dto.PolygonAnalysisSummary) *tqdpb.PolygonAnalysisSummary {
	if summary == nil {
		return nil
	}

	return &tqdpb.PolygonAnalysisSummary{
		ResultType:          summary.ResultType,
		Z:                   summary.Z,
		QueryAreaSqm:        summary.QueryAreaSqm,
		TotalMatchedAreaSqm: summary.TotalMatchedAreaSqm,
		RegionCount:         summary.RegionCount,
		ParcelCount:         summary.ParcelCount,
		LandUseStats:        polygonLandUseStatsToProto(summary.LandUseStats),
		LayerStats:          polygonLayerStatsToProto(summary.LayerStats),
		BuildabilityStats:   polygonBuildabilityStatsToProto(summary.BuildabilityStats),
	}
}

func (m *ParcelMapper) PolygonMapTargetsToProto(
	resp *dto.PolygonMapTargetsResponse,
) *tqdpb.FindParcelsByPolygonResponse {
	if resp == nil {
		return &tqdpb.FindParcelsByPolygonResponse{
			Data:       []*tqdpb.Parcel{},
			Regions:    []*tqdpb.RegionInfoResponse{},
			Total:      0,
			Page:       1,
			PageSize:   0,
			ResultType: WorkspaceEntityTypeParcel,
			Z:          0,
		}
	}

	resultType := resp.ResultType
	if resultType == 0 {
		/**
		 * Fallback an toàn:
		 * - Có region thì coi là region result.
		 * - Ngược lại coi là parcel result để giữ legacy.
		 */
		if len(resp.Regions) > 0 {
			resultType = WorkspaceEntityTypeRegion
		} else {
			resultType = WorkspaceEntityTypeParcel
		}
	}

	return &tqdpb.FindParcelsByPolygonResponse{
		Data:    m.DomainListToProto(resp.Parcels),
		Regions: m.RegionInfoListToProto(resp.Regions),

		Total:    int32(resp.Total),
		Page:     int32(resp.Page),
		PageSize: int32(resp.PageSize),

		AdjustedGeoJson:   resp.AdjustedGeoJSON,
		AdjustmentReason:  resp.AdjustmentReason,
		OriginalAreaKm2:   resp.OriginalAreaKm2,
		AdjustedAreaKm2:   resp.AdjustedAreaKm2,
		MaxAllowedAreaKm2: resp.MaxAllowedAreaKm2,
		Warning:           resp.Warning,

		ResultType: resultType,
		Z:          resp.Z,
		Analysis:   polygonAnalysisSummaryToProto(resp.Analysis),
	}
}

func regionImpactToProto(region *dto.RegionImpact) *tqdpb.Region {
	if region == nil {
		return nil
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
	return &tqdpb.Region{
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
	}
}

func parcelPlanInfoToProto(plan *dto.ParcelPlanInfo, groups []*dto.PlanLandUseGroupDetail) *tqdpb.PlanDetailInfo {
	if plan == nil {
		return nil
	}

	return &tqdpb.PlanDetailInfo{
		LayerId:          plan.LayerID,
		LayerName:        plan.LayerName,
		DisplayName:      plan.DisplayName,
		LandUseCode:      plan.LandUseCode,
		LandUseName:      plan.LandUseName,
		LandUseColor:     plan.LandUseColor,
		GroupCode:        plan.GroupCode,
		GroupName:        plan.GroupName,
		CanBuild:         plan.CanBuild,
		ImpactPercent:    plan.ImpactPercent,
		ImpactAreaSqm:    plan.ImpactAreaSqm,
		LandUseGroups:    overviewLandUseGroupsToProto(groups),
		LegalCanonical:   plan.LegalCanonical,
		LegalConfidence:  plan.LegalConfidence,
		LegalActiveState: plan.LegalActiveState,
	}
}

func overviewLandUseGroupsToProto(groups []*dto.PlanLandUseGroupDetail) []*tqdpb.PlanLandUseGroupDetail {
	if len(groups) == 0 {
		return nil
	}

	result := make([]*tqdpb.PlanLandUseGroupDetail, 0, len(groups))
	for _, g := range groups {
		if g == nil {
			continue
		}
		result = append(result, &tqdpb.PlanLandUseGroupDetail{
			Code:           g.Code,
			Name:           g.Name,
			Color:          g.Color,
			CanBuild:       g.CanBuild,
			Priority:       int32(g.Priority),
			LandUseCode:    g.LandUseCode,
			LandUseName:    g.LandUseName,
			LandUseColor:   g.LandUseColor,
			OverlapPercent: g.OverlapPercent,
			OverlapAreaSqm: g.OverlapAreaSqm,
		})
	}

	return result
}

// ToProtoParcelOverview - Map DTO → Proto OverviewResponse
func (m *ParcelMapper) ToProtoParcelOverview(data *dto.ParcelLayersResponse) *tqdpb.ParcelOverviewResponse {
	resp := data

	layers := make([]*tqdpb.Layer, 0, len(data.Layers))
	for _, layer := range data.Layers {
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
				regions = append(regions, regionImpactToProto(region))
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

		var protoLegalDocs []*tqdpb.LayerLegalDocInfo
		for _, ld := range layer.LegalDocs {
			protoLegalDocs = append(protoLegalDocs, &tqdpb.LayerLegalDocInfo{
				Id:       ld.ID,
				Name:     ld.Name,
				FileUrl:  ld.FileUrl,
				FileType: ld.FileType,
			})
		}

		var protoAuthority *tqdpb.AuthorityIssuingInfo
		if layer.AuthorityIssuing != nil {
			protoAuthority = &tqdpb.AuthorityIssuingInfo{
				Id:          layer.AuthorityIssuing.ID,
				Name:        layer.AuthorityIssuing.Name,
				Code:        layer.AuthorityIssuing.Code,
				Description: layer.AuthorityIssuing.Description,
			}
		}

		layers = append(layers, &tqdpb.Layer{
			Id:               layer.ID,
			Name:             layer.Name,
			DisplayName:      layer.DisplayName,
			Type:             layer.Type,
			Severity:         layer.Severity,
			EffectiveDate:    layer.EffectiveDate,
			LegalDoc:         legalDocStr,
			Labels:           labels,
			TotalArea:        layer.TotalArea,
			TotalPercent:     layer.TotalPct,
			LegalDocs:        protoLegalDocs,
			AuthorityIssuing: protoAuthority,
		})
	}

	summary := &tqdpb.ParcelSummary{}
	if resp.Summary != nil {
		summary = &tqdpb.ParcelSummary{
			DominantLandUse:          resp.Summary.DominantLandUse,
			DominantPercentage:       resp.Summary.DominantPercentage,
			TotalAffectedPercentage:  resp.Summary.TotalAffectedPercentage,
			TopLayerName:             resp.Summary.TopLayerName,
			Buildable:                resp.Summary.Buildable,
			DominantLandUseCode:      resp.Summary.DominantLandUseCode,
			DominantLandUseColor:     resp.Summary.DominantLandUseColor,
			DominantLandUseGroupCode: resp.Summary.DominantLandUseGroupCode,
			DominantLandUseGroupName: resp.Summary.DominantLandUseGroupName,
			DominantAreaSqm:          resp.Summary.DominantAreaSqm,
			DominantCanBuild:         resp.Summary.DominantCanBuild,
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

	return &tqdpb.ParcelOverviewResponse{
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
		ConcludeLabel:        regionImpactToProto(resp.ConcludeLabel),
		Info:                 m.ToProtoParcelInfoDetail(resp.ParcelInfo),
		PrimaryPlan:          parcelPlanInfoToProto(resp.PrimaryPlan, resp.LandUseGroups),
		LandUseGroups:        overviewLandUseGroupsToProto(resp.LandUseGroups),
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
		IsSeo:      p.IsSeo,
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
		IsSeo:      p.IsSeo,
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

func (m *ParcelMapper) getLandTypeName(landTypeId uint64) string {
	if name, ok := m.LandTypeNames[landTypeId]; ok {
		return name
	}
	return "Đất khác"
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func mapLandUseGroups(groupsMap map[uint64][]*types.LandUseGroupDetail, layerID uint64) []*tqdpb.PlanLandUseGroupDetail {
	if groupsMap == nil {
		return nil
	}
	groups, ok := groupsMap[layerID]
	if !ok || len(groups) == 0 {
		return nil
	}
	result := make([]*tqdpb.PlanLandUseGroupDetail, 0, len(groups))
	for _, g := range groups {
		result = append(result, &tqdpb.PlanLandUseGroupDetail{
			Code:           g.Code,
			Name:           g.Name,
			Color:          g.Color,
			CanBuild:       g.CanBuild,
			Priority:       int32(g.Priority),
			LandUseCode:    g.LandUseCode,
			LandUseName:    g.LandUseName,
			LandUseColor:   g.LandUseColor,
			OverlapPercent: g.OverlapPercent,
			OverlapAreaSqm: g.OverlapAreaSqm,
		})
	}
	return result
}

func (m *ParcelMapper) buildAssessmentSummary(result *types.ResolveResult) string {
	if result == nil {
		return ""
	}
	if result.Primary == nil {
		return result.Reason
	}
	summary := result.Primary.LandUseName
	if result.Primary.OverlapPercent > 0 {
		summary += fmt.Sprintf(" (%.1f%%)", result.Primary.OverlapPercent)
	}
	if result.HasConflict {
		summary += m.ConflictSummarySuffix
	}
	return summary
}

// ─── Types Layer → Proto Layer conversion ──────────────────────────────────

func mapTypeLayerToProto(l *types.LayerDetail) *tqdpb.Layer {
	if l == nil {
		return nil
	}
	pl := &tqdpb.Layer{
		Id:            l.ID,
		Name:          l.Name,
		DisplayName:   l.DisplayName,
		Type:          l.Type,
		EffectiveDate: l.EffectiveDate,
		TotalArea:     l.TotalArea,
		TotalPercent:  l.TotalPct,
	}
	if l.Authority != nil {
		pl.AuthorityIssuing = &tqdpb.AuthorityIssuingInfo{
			Id:          l.Authority.ID,
			Name:        l.Authority.Name,
			Code:        l.Authority.Code,
			Description: l.Authority.Description,
		}
	}
	for _, d := range l.LegalDocs {
		pl.LegalDocs = append(pl.LegalDocs, &tqdpb.LayerLegalDocInfo{
			Id:       d.ID,
			Name:     d.Name,
			FileUrl:  d.FileUrl,
			FileType: d.FileType,
		})
	}
	for _, lbl := range l.Labels {
		plb := &tqdpb.Label{
			Id:           lbl.ID,
			Name:         lbl.Name,
			DisplayName:  lbl.DisplayName,
			Color:        lbl.Color,
			TotalArea:    lbl.TotalArea,
			TotalPercent: lbl.TotalPct,
		}
		for _, r := range lbl.Regions {
			plb.Regions = append(plb.Regions, &tqdpb.Region{
				Id:             r.ID,
				Name:           r.Name,
				DisplayName:    r.DisplayName,
				LandUseCode:    r.LandUseCode,
				LandUseName:    r.LandUseName,
				OverlapArea:    r.OverlapAreaSqm,
				OverlapPercent: r.OverlapPct,
				Geometry:       r.Geometry,
				Color: &tqdpb.RGBColor{
					Red:   uint32(r.ColorRed),
					Green: uint32(r.ColorGreen),
					Blue:  uint32(r.ColorBlue),
				},
			})
		}
		pl.Labels = append(pl.Labels, plb)
	}
	return pl
}

func mapTypeLayerToLayerDetail(l *types.LayerDetail) *tqdpb.LayerDetail {
	if l == nil {
		return nil
	}
	pl := &tqdpb.LayerDetail{
		Id:            l.ID,
		Name:          l.Name,
		DisplayName:   l.DisplayName,
		Type:          l.Type,
		EffectiveDate: l.EffectiveDate,
		TotalArea:     l.TotalArea,
		TotalPct:      l.TotalPct,
	}
	if l.Authority != nil {
		pl.Authority = &tqdpb.AuthorityDetail{
			Id:          l.Authority.ID,
			Name:        l.Authority.Name,
			Code:        l.Authority.Code,
			Description: l.Authority.Description,
		}
	}
	for _, d := range l.LegalDocs {
		pl.LegalDocs = append(pl.LegalDocs, &tqdpb.LegalDocDetail{
			Id:       d.ID,
			Name:     d.Name,
			FileUrl:  d.FileUrl,
			FileType: d.FileType,
		})
	}
	for _, lbl := range l.Labels {
		plb := &tqdpb.LabelDetail{
			Id:          lbl.ID,
			Name:        lbl.Name,
			DisplayName: lbl.DisplayName,
			Color:       lbl.Color,
			TotalArea:   lbl.TotalArea,
			TotalPct:    lbl.TotalPct,
		}
		for _, r := range lbl.Regions {
			plb.Regions = append(plb.Regions, &tqdpb.RegionDetail{
				Id:             r.ID,
				Name:           r.Name,
				DisplayName:    r.DisplayName,
				LandUseCode:    r.LandUseCode,
				LandUseName:    r.LandUseName,
				ColorRed:       uint32(r.ColorRed),
				ColorGreen:     uint32(r.ColorGreen),
				ColorBlue:      uint32(r.ColorBlue),
				OverlapAreaSqm: r.OverlapAreaSqm,
				OverlapPct:     r.OverlapPct,
				Geometry:       r.Geometry,
			})
		}
		pl.Labels = append(pl.Labels, plb)
	}
	return pl
}
