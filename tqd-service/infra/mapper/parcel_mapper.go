package mapper

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	_utils "common/utils"
	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	qh_dto "tqd/internal/dto/qh"
	"tqd/internal/enums"
	"tqd/internal/usecase/resolver/layer_resolver/types"
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

func (m *ParcelMapper) ToProtoParcelDetail(data *qh_dto.ParcelDetailResponseDTO) *tqdpb.ParcelDetailResponse {
	if data == nil {
		return nil
	}

	resp := &tqdpb.ParcelDetailResponse{
		Parcel: &tqdpb.ParcelIdentity{
			Id:               data.Parcel.ID,
			MapSheetNumber:   data.Parcel.MapSheetNumber,
			LandParcelNumber: data.Parcel.LandParcelNumber,
			Address:          data.Parcel.Address,
			TotalAreaSqm:     data.Parcel.TotalAreaSqm,
		},
		CurrentUse: &tqdpb.CurrentUse{
			Code:            data.CurrentUse.Code,
			Name:            data.CurrentUse.Name,
			Color:           data.CurrentUse.Color,
			BuildStatus:     m.parseBuildStatus(data.CurrentUse.BuildStatus),
			LegalDocumentId: data.CurrentUse.LegalDocumentId,
		},
	}

	// Centroid
	if data.Parcel.Centroid != nil {
		resp.Parcel.Centroid = &tqdpb.Centroid{
			Lat: data.Parcel.Centroid.Lat,
			Lng: data.Parcel.Centroid.Lng,
		}
	}

	// Ownership
	if data.Parcel.Ownership != nil {
		resp.Parcel.Ownership = &tqdpb.Ownership{
			CertificateNumber: data.Parcel.Ownership.CertificateNumber,
			IssuedAt:          data.Parcel.Ownership.IssuedAt,
			IssuedBy:          data.Parcel.Ownership.IssuedBy,
			HolderName:        data.Parcel.Ownership.HolderName,
			LegalDocumentId:   data.Parcel.Ownership.LegalDocumentId,
		}
	}

	// Location
	if data.Parcel.Location != nil {
		resp.Parcel.Location = &tqdpb.Location{
			Province:     data.Parcel.Location.Province,
			ProvinceCode: data.Parcel.Location.ProvinceCode,
			WardCode:     data.Parcel.Location.WardCode,
		}
	}

	// ✅ Planning Conclusion — đã là pointer
	if data.PlanningConclusion != nil {
		pc := data.PlanningConclusion
		resp.PlanningConclusion = &tqdpb.PlanningConclusion{
			BuildStatus: m.parseBuildStatus(pc.BuildStatus),
			HasConflict: pc.HasConflict,
		}

		if pc.Alert != nil {
			resp.PlanningConclusion.Alert = &tqdpb.Alert{
				Level:   m.parseAlertLevel(pc.Alert.Level),
				Score:   pc.Alert.Score,
				Message: pc.Alert.Message,
			}
		}

		if pc.LegalBasis != nil {
			resp.PlanningConclusion.LegalBasis = &tqdpb.LegalBasis{
				Summary:          pc.LegalBasis.Summary,
				LegalDocumentIds: pc.LegalBasis.LegalDocumentIds,
			}
		}

		for _, b := range pc.Breakdown {
			resp.PlanningConclusion.Breakdown = append(resp.PlanningConclusion.Breakdown, &tqdpb.PlanningBreakdown{
				Rank:             uint32(b.Rank),
				LandUseCode:      b.LandUseCode,
				LandUseName:      b.LandUseName,
				LandUseColor:     b.LandUseColor,
				GroupCode:        b.GroupCode,
				GroupName:        b.GroupName,
				AreaSqm:          b.AreaSqm,
				Percent:          b.Percent,
				BuildStatus:      m.parseBuildStatus(b.BuildStatus),
				AlertLevel:       m.parseAlertLevel(b.AlertLevel),
				LegalDocumentIds: b.LegalDocumentIds,
			})
		}

		if pc.Statistics != nil {
			resp.PlanningConclusion.Statistics = &tqdpb.PlanningStatistics{
				TotalLayersAffected: uint32(pc.Statistics.TotalLayersAffected),
				TotalZonesAffected:  uint32(pc.Statistics.TotalZonesAffected),
				TotalOverlapAreaSqm: pc.Statistics.TotalOverlapAreaSqm,
				TotalOverlapPercent: pc.Statistics.TotalOverlapPercent,
			}
		}
	}

	if data.PrimaryUse != nil {
		resp.PrimaryUse = &tqdpb.PrimaryUse{
			Rule:        data.PrimaryUse.Rule,
			LandUseCode: data.PrimaryUse.LandUseCode,
			LandUseName: data.PrimaryUse.LandUseName,
			Percent:     data.PrimaryUse.Percent,
			BuildStatus: m.parseBuildStatus(data.PrimaryUse.BuildStatus),
		}
	}

	if data.Guidance != nil {
		resp.Guidance = &tqdpb.Guidance{}

		for _, w := range data.Guidance.Warnings {
			warning := &tqdpb.Warning{
				AlertLevel:      m.parseAlertLevel(w.AlertLevel),
				LandUseCode:     w.LandUseCode,
				AffectedAreaSqm: w.AffectedAreaSqm,
				AffectedPercent: w.AffectedPercent,
				Title:           w.Title,
				Action:          w.Action,
			}
			if w.LegalBasis != nil {
				warning.LegalBasis = &tqdpb.LegalBasis{
					Summary:          w.LegalBasis.Summary,
					LegalDocumentIds: w.LegalBasis.LegalDocumentIds,
				}
			}
			resp.Guidance.Warnings = append(resp.Guidance.Warnings, warning)
		}

		for _, r := range data.Guidance.Recommendations {
			resp.Guidance.Recommendations = append(resp.Guidance.Recommendations, &tqdpb.Recommendation{
				Type:             r.Type,
				Message:          r.Message,
				LegalDocumentIds: r.LegalDocumentIds,
			})
		}
	}

	if data.LegalDocuments != nil {
		resp.LegalDocuments = &tqdpb.LegalDocumentsSummary{
			Items: m.ToProtoLegalDocumentItemsFromDTO(data.LegalDocuments.Items),
			Summary: &tqdpb.LegalDocumentSummary{
				Total:             uint32(data.LegalDocuments.Total),
				ByType:            m.convertMapStringUint32(data.LegalDocuments.ByType),
				ByStatus:          m.convertMapStringUint32(data.LegalDocuments.ByStatus),
				HasSupersededDocs: data.LegalDocuments.HasSupersededDocs,
			},
		}
	}

	// Meta
	if data.Meta != nil {
		resp.Meta = &tqdpb.ResponseMetadata{
			DataVersion:    data.Meta.DataVersion,
			DataQuality:    data.Meta.DataQuality,
			QueriedAt:      data.Meta.QueriedAt,
			ResponseTimeMs: uint32(data.Meta.ResponseTimeMs),
			PartialData:    data.Meta.PartialData,
		}
	}

	// Links
	if data.Links != nil {
		resp.Links = &tqdpb.Links{
			Self: data.Links.Self,
		}
		if data.Links.Layers != "" {
			resp.Links.Layers = &data.Links.Layers
		}
		if data.Links.LegalDocuments != "" {
			resp.Links.LegalDocuments = &data.Links.LegalDocuments
		}
		if data.Links.ShareUrl != "" {
			resp.Links.ShareUrl = &data.Links.ShareUrl
		}
		if data.Links.Next != "" {
			resp.Links.Next = &data.Links.Next
		}
		if data.Links.Prev != "" {
			resp.Links.Prev = &data.Links.Prev
		}
	}

	return resp
}

func (m *ParcelMapper) ToProtoLegalDocumentSummary(docs []*qh_domain.QHLegalDocument) *tqdpb.LegalDocumentSummary {
	if len(docs) == 0 {
		return &tqdpb.LegalDocumentSummary{
			Total:             0,
			ByType:            make(map[string]uint32),
			ByStatus:          make(map[string]uint32),
			HasSupersededDocs: false,
		}
	}

	summary := &tqdpb.LegalDocumentSummary{
		Total:             uint32(len(docs)),
		ByType:            make(map[string]uint32),
		ByStatus:          make(map[string]uint32),
		HasSupersededDocs: false,
	}

	for _, d := range docs {
		if d == nil {
			continue
		}

		// ByType — key là string để frontend dễ hiển thị
		typeKey := enums.DocumentTypeString[d.DocumentType]
		if typeKey == "" {
			typeKey = "OTHER"
		}
		summary.ByType[typeKey]++

		// ByStatus — key là string
		statusKey := enums.DocumentStatusString[d.Status]
		if statusKey == "" {
			statusKey = "UNKNOWN"
		}
		summary.ByStatus[statusKey]++

		// Check superseded
		if d.IsSuperseded() {
			summary.HasSupersededDocs = true
		}
	}

	return summary
}

// ToProtoLegalDocumentSummaryFromDTO maps []qh_dto.LegalDocumentDTO → *tqdpb.LegalDocumentSummary
func (m *ParcelMapper) ToProtoLegalDocumentSummaryFromDTO(docs []qh_dto.LegalDocumentDTO) *tqdpb.LegalDocumentSummary {
	if len(docs) == 0 {
		return &tqdpb.LegalDocumentSummary{
			Total:             0,
			ByType:            make(map[string]uint32),
			ByStatus:          make(map[string]uint32),
			HasSupersededDocs: false,
		}
	}

	summary := &tqdpb.LegalDocumentSummary{
		Total:             uint32(len(docs)),
		ByType:            make(map[string]uint32),
		ByStatus:          make(map[string]uint32),
		HasSupersededDocs: false,
	}

	for _, d := range docs {
		// ByType
		summary.ByType[d.Type]++

		// ByStatus
		summary.ByStatus[d.Status]++

		// Check superseded
		if d.SupersededBy != nil && *d.SupersededBy != "" {
			summary.HasSupersededDocs = true
		}
	}

	return summary
}

func (m *ParcelMapper) ToProtoLayerDetailItems(layers []*dto.LayerImpact, parcelID uint64) []*tqdpb.LayerDetailItem {
	if len(layers) == 0 {
		return []*tqdpb.LayerDetailItem{}
	}

	result := make([]*tqdpb.LayerDetailItem, 0, len(layers))
	for _, l := range layers {
		if l == nil {
			continue
		}

		item := &tqdpb.LayerDetailItem{
			LayerId:          l.ID,
			LayerCode:        l.Name,
			LayerName:        l.Name,
			LayerDisplayName: l.DisplayName,
			LayerType:        "PLANNING",
			LayerDescription: "",
			LegalStatus: &tqdpb.LayerLegalStatus{
				Code:            m.getLegalStatusCode(l.LegalStatus),
				Label:           l.LegalStatusName,
				EffectiveDate:   l.EffectiveDate,
				IssuedDate:      l.EffectiveDate,
				LegalDocumentId: "",
			},
			LegalDocumentIds: []string{},
			ZoneCount:        uint32(len(l.Labels)),
			TotalArea:        l.TotalArea,
			TotalPercent:     l.TotalPct,
		}

		if l.AuthorityIssuing != nil {
			item.IssuingAuthority = &tqdpb.IssuingAuthority{
				Id:   l.AuthorityIssuing.ID,
				Name: l.AuthorityIssuing.Name,
				Code: l.AuthorityIssuing.Code,
			}
		}

		for _, label := range l.Labels {
			for _, region := range label.Regions {
				zone := &tqdpb.ZoneItem{
					ZoneId:           region.ID,
					ZoneName:         region.Name,
					ZoneDisplayName:  region.DisplayName,
					LandUseCode:      region.LandUseCode,
					LandUseName:      region.LandUseName,
					LandUseColor:     region.ColorRGB,
					AreaSqm:          region.OverlapAreaSqm,
					Percent:          region.OverlapPct,
					BuildStatus:      m.parseBuildStatusFromBool(region.CanBuild),
					AlertLevel:       m.parseAlertLevelFromWarnLevel(region.WarnLevel),
					LegalDocumentIds: []string{},
					Centroid: &tqdpb.Centroid{
						Lat: region.CenterLat,
						Lng: region.CenterLng,
					},
					Links: &tqdpb.ZoneLinks{
						Geometry: fmt.Sprintf("/v2/tqd/parcels/%d/layers/%d/zones/%d/geometry", parcelID, l.ID, region.ID),
					},
				}
				item.Zones = append(item.Zones, zone)
			}
		}

		result = append(result, item)
	}

	return result
}

// ============================================================
// PHASE 3: LEGAL DOCUMENTS — ToProtoLegalDocumentItems
// ============================================================

// ToProtoLegalDocumentItems maps []qh_domain.QHLegalDocument → []*tqdpb.LegalDocumentItem
// ENUM trả về uint32, không convert sang string
func (m *ParcelMapper) ToProtoLegalDocumentItems(docs []*qh_domain.QHLegalDocument) []*tqdpb.LegalDocumentItem {
	if len(docs) == 0 {
		return []*tqdpb.LegalDocumentItem{}
	}

	result := make([]*tqdpb.LegalDocumentItem, 0, len(docs))
	for _, d := range docs {
		if d == nil {
			continue
		}

		item := &tqdpb.LegalDocumentItem{
			Id:               strconv.FormatUint(d.ID, 10),
			Type:             uint32(d.DocumentType), // uint32 trực tiếp
			TypeName:         enums.DocumentTypeName[d.DocumentType],
			TypeIcon:         m.getTypeIcon(d.DocumentType),
			Name:             d.Name,
			FullName:         d.FullName,
			IssuedBy:         d.IssuedBy,
			IssuedAt:         _utils.FormatRFC3339(d.IssuedAt),
			EffectiveDate:    _utils.FormatRFC3339(d.EffectiveAt),
			Status:           uint32(d.Status), // uint32 trực tiếp
			StatusLabel:      enums.DocumentStatusLabel[d.Status],
			StatusColor:      enums.DocumentStatusColor[d.Status],
			RelevantArticles: d.RelevantArticles,
			IsPublic:         d.IsPublic,
			RequiresAuth:     d.RequiresAuth,
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

		if d.FileURL != "" || d.DownloadURL != "" {
			item.File = &tqdpb.LegalFile{
				Type:        d.FileType,
				SizeKb:      uint32(d.FileSizeKb),
				Pages:       uint32(d.PageCount),
				Url:         stringPtr(d.FileURL),
				PreviewUrl:  stringPtr(d.PreviewURL),
				DownloadUrl: stringPtr(d.DownloadURL),
			}
		}

		if len(d.AffectedLayerIDs) > 0 || len(d.AffectedZoneIDs) > 0 {
			item.AffectsEntities = &tqdpb.AffectsEntities{
				LayerIds: d.AffectedLayerIDs,
				ZoneIds:  d.AffectedZoneIDs,
			}
		}

		result = append(result, item)
	}

	return result
}

// ToProtoLegalDocumentItemsFromDTO maps []qh_dto.LegalDocumentDTO → []*tqdpb.LegalDocumentItem
func (m *ParcelMapper) ToProtoLegalDocumentItemsFromDTO(docs []qh_dto.LegalDocumentDTO) []*tqdpb.LegalDocumentItem {
	if len(docs) == 0 {
		return []*tqdpb.LegalDocumentItem{}
	}

	result := make([]*tqdpb.LegalDocumentItem, 0, len(docs))
	for _, d := range docs {
		item := &tqdpb.LegalDocumentItem{
			Id:               d.ID,
			Type:             m.parseDocumentType(d.Type),
			TypeName:         d.TypeName,
			TypeIcon:         d.TypeIcon,
			Name:             d.Name,
			FullName:         d.FullName,
			IssuedBy:         d.IssuedBy,
			IssuedAt:         d.IssuedAt,
			EffectiveDate:    d.EffectiveDate,
			Status:           m.parseDocumentStatus(d.Status),
			StatusLabel:      d.StatusLabel,
			StatusColor:      d.StatusColor,
			RelevantArticles: d.RelevantArticles,
			IsPublic:         d.IsPublic,
			RequiresAuth:     d.RequiresAuth,
		}

		if d.ExpiryDate != nil {
			item.ExpiryDate = d.ExpiryDate
		}
		if d.SupersededBy != nil {
			item.SupersededBy = d.SupersededBy
		}
		item.Supersedes = d.Supersedes

		if d.File != nil {
			item.File = &tqdpb.LegalFile{
				Type:        d.File.Type,
				SizeKb:      uint32(d.File.SizeKb),
				Pages:       uint32(d.File.Pages),
				Url:         d.File.URL,
				PreviewUrl:  d.File.PreviewURL,
				DownloadUrl: d.File.DownloadURL,
			}
		}

		if d.AffectsEntities != nil {
			item.AffectsEntities = &tqdpb.AffectsEntities{
				LayerIds: d.AffectsEntities.LayerIDs,
				ZoneIds:  d.AffectsEntities.ZoneIDs,
			}
		}

		result = append(result, item)
	}

	return result
}

// ============================================================
// HELPER FUNCTIONS — parse string sang uint32
// ============================================================

func (m *ParcelMapper) parseBuildStatus(status string) uint32 {
	switch status {
	case "ALLOWED":
		return uint32(enums.BuildStatusAllowed)
	case "FORBIDDEN":
		return uint32(enums.BuildStatusForbidden)
	case "RESTRICTED":
		return uint32(enums.BuildStatusRestricted)
	default:
		return uint32(enums.BuildStatusUnknown)
	}
}

func (m *ParcelMapper) parseBuildStatusFromBool(canBuild bool) uint32 {
	if canBuild {
		return uint32(enums.BuildStatusAllowed)
	}
	return uint32(enums.BuildStatusForbidden)
}

func (m *ParcelMapper) parseAlertLevel(level string) uint32 {
	switch level {
	case "LOW":
		return uint32(enums.AlertLevelLow)
	case "MEDIUM":
		return uint32(enums.AlertLevelMedium)
	case "HIGH":
		return uint32(enums.AlertLevelHigh)
	default:
		return uint32(enums.AlertLevelNone)
	}
}

func (m *ParcelMapper) parseAlertLevelFromWarnLevel(warnLevel int32) uint32 {
	switch warnLevel {
	case 1:
		return uint32(enums.AlertLevelLow)
	case 2:
		return uint32(enums.AlertLevelMedium)
	case 3:
		return uint32(enums.AlertLevelHigh)
	default:
		return uint32(enums.AlertLevelNone)
	}
}

func (m *ParcelMapper) parseDocumentType(docType string) uint32 {
	switch docType {
	case "PLANNING_DECISION":
		return uint32(enums.DocTypePlanningDecision)
	case "LAND_CERTIFICATE":
		return uint32(enums.DocTypeLandCertificate)
	case "LAW":
		return uint32(enums.DocTypeLaw)
	case "DECREE":
		return uint32(enums.DocTypeDecree)
	case "REGULATION":
		return uint32(enums.DocTypeRegulation)
	case "CIRCULAR":
		return uint32(enums.DocTypeCircular)
	default:
		return uint32(enums.DocTypeOther)
	}
}

func (m *ParcelMapper) parseDocumentStatus(status string) uint32 {
	switch status {
	case "ACTIVE":
		return uint32(enums.DocStatusActive)
	case "SUPERSEDED":
		return uint32(enums.DocStatusSuperseded)
	case "EXPIRED":
		return uint32(enums.DocStatusExpired)
	case "DRAFT":
		return uint32(enums.DocStatusDraft)
	default:
		return uint32(enums.DocStatusActive)
	}
}

func (m *ParcelMapper) getLegalStatusCode(status uint32) string {
	switch status {
	case 10:
		return "OFFICIAL"
	case 20:
		return "ACTIVE"
	case 30:
		return "EXPIRED"
	case 40:
		return "ADJUST"
	default:
		return "UNKNOWN"
	}
}

func (m *ParcelMapper) getTypeIcon(docType enums.DocumentType) string {
	switch docType {
	case enums.DocTypePlanningDecision:
		return "file-certificate"
	case enums.DocTypeLandCertificate:
		return "certificate"
	case enums.DocTypeLaw:
		return "gavel"
	case enums.DocTypeDecree:
		return "scroll"
	case enums.DocTypeRegulation:
		return "book"
	case enums.DocTypeCircular:
		return "mail"
	default:
		return "file"
	}
}

func (m *ParcelMapper) convertMapStringUint32(mp map[string]int) map[string]uint32 {
	if mp == nil {
		return nil
	}
	result := make(map[string]uint32, len(mp))
	for k, v := range mp {
		result[k] = uint32(v)
	}
	return result
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// ============================================================
// EXISTING METHODS — GIỮ NGUYÊN (không thay đổi)
// ============================================================

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

	if result.HistoricalRisk != nil {
		resp.HistoricalRisk = &tqdpb.HistoricalRisk{
			Trend:          result.HistoricalRisk.Trend,
			TrendDetail:    result.HistoricalRisk.TrendDetail,
			ActiveLayers:   int32(result.HistoricalRisk.ActiveLayers),
			ExpiringLayers: int32(result.HistoricalRisk.ExpiringLayers),
			AnalyzedAt:     result.HistoricalRisk.AnalyzedAt,
		}
	}

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
		RegionId:         v.RegionID,
		LayerId:          v.LayerID,
		LabelId:          v.LabelID,
		LandUseId:        v.LandUseID,
		LegendId:         v.LegendID,
		Name:             v.Name,
		DisplayName:      v.DisplayName,
		Description:      v.Description,
		LayerName:        v.LayerName,
		LayerDisplayName: v.LayerDisplayName,
		LandUseCode:      v.LandUseCode,
		LandUseName:      v.LandUseName,
		LandUseColor:     v.LandUseColor,
		CanBuild:         v.CanBuild,
		LabelName:        v.LabelName,
		LabelColor:       v.LabelColor,
		LegendColor:      v.LegendColor,
		LegendType:       v.LegendType,
		GeometryType:     v.GeometryType,
		LegalDoc:         v.LegalDoc,
		PlanningName:     v.PlanningName,
		AreaSqm:          v.AreaSqm,
		AreaHa:           v.AreaHa,
		PerimeterM:       v.PerimeterM,
		CenterLat:        v.CenterLat,
		CenterLon:        v.CenterLon,
		Bounds: &tqdpb.RegionBounds{
			MinLon: v.MinLon,
			MinLat: v.MinLat,
			MaxLon: v.MaxLon,
			MaxLat: v.MaxLat,
		},
		GeoJson:      v.GeoJSON,
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

func (m *ParcelMapper) PolygonMapTargetsToProto(resp *dto.PolygonMapTargetsResponse) *tqdpb.FindParcelsByPolygonResponse {
	if resp == nil {
		return &tqdpb.FindParcelsByPolygonResponse{
			Parcels:    []*tqdpb.Parcel{},
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
		hasParcels := len(resp.Parcels) > 0
		hasRegions := len(resp.Regions) > 0
		switch {
		case hasParcels && hasRegions:
			// Mixed: giữ 0 để client đọc cả data (parcels) và regions
			resultType = 0
		case hasRegions:
			resultType = WorkspaceEntityTypeRegion
		default:
			resultType = WorkspaceEntityTypeParcel
		}
	}
	var layers []*tqdpb.LayerQuickInfo
	var quickWarnings []string
	if resp.QuickLayerRows != nil {
		quickLayers := m.BuildLayerQuickInfoFromRows(resp.QuickLayerRows)
		layers = quickLayers.Layers
		quickWarnings = quickLayers.Warnings
	}
	return &tqdpb.FindParcelsByPolygonResponse{
		Parcels:           m.DomainListToProto(resp.Parcels),
		Regions:           m.RegionInfoListToProto(resp.Regions),
		Total:             int32(resp.Total),
		Page:              int32(resp.Page),
		PageSize:          int32(resp.PageSize),
		AdjustedGeoJson:   resp.AdjustedGeoJSON,
		AdjustmentReason:  resp.AdjustmentReason,
		OriginalAreaKm2:   resp.OriginalAreaKm2,
		AdjustedAreaKm2:   resp.AdjustedAreaKm2,
		MaxAllowedAreaKm2: resp.MaxAllowedAreaKm2,
		Warning:           resp.Warning,
		ResultType:        resultType,
		Z:                 resp.Z,
		Analysis:          polygonAnalysisSummaryToProto(resp.Analysis),
		Layers:            layers,
		QuickWarnings:     quickWarnings,
	}
}

func (m *ParcelMapper) BuildLayerQuickInfoFromRows(rows []dto.ParcelLayerInfo) *tqdpb.LayerQuickInfoResponse {
	layerMap := make(map[uint64]*tqdpb.LayerQuickInfo)
	layers := make([]*tqdpb.LayerQuickInfo, 0)
	warnings := make([]string, 0)

	for _, q := range rows {
		layer, exists := layerMap[q.LayerID]
		if !exists {
			layer = &tqdpb.LayerQuickInfo{
				ParcelId:         q.ParcelID,
				LayerId:          q.LayerID,
				LayerName:        q.LayerName,
				LayerDisplayName: q.LayerDisplayName,
				LayerType:        q.LayerType,
				LegalStatus:      uint32(q.LayerLegalStatus),
				LayerOrder:       q.LayerOrder,
				LegalStatusName:  enums.LegalStatusMap[enums.LegalStatus(q.LayerLegalStatus)],
				LayerAvatar:      q.LayerAvatar,
				RelationType:     q.RelationType,
				UpdatedAt:        _utils.FormatTimeToString(&q.LayerUpdatedAt),
				LandUses:         make([]*tqdpb.LandUseQuickInfo, 0),

				PlanningProjectId: q.PlanningProjectID,
			}

			layerMap[q.LayerID] = layer
			layers = append(layers, layer)
		}

		layer.LandUses = append(layer.LandUses, &tqdpb.LandUseQuickInfo{
			LandUseId:        q.LandUseID,
			LandUseName:      q.LandUseName,
			LandUseColor:     q.LandUseColor,
			WarnLevel:        uint32(q.WarnLevel),
			IntersectAreaSqm: q.IntersectAreaSqm,
			Percent:          q.Percent,
			CanBuild:         q.CanBuild,
		})

		if q.WarnLevel > 0 {
			warnings = append(
				warnings,
				fmt.Sprintf(
					"%.2f%% diện tích thuộc '%s' trong '%s'.",
					q.Percent,
					q.LandUseName,
					q.LayerName,
				),
			)
		}
	}

	sort.SliceStable(layers, func(i, j int) bool {
		return layers[i].LayerOrder > layers[j].LayerOrder
	})

	// if len(warnings) == 0 {
	// 	warnings = append(warnings, "Không có cảnh báo nào.")
	// }

	return &tqdpb.LayerQuickInfoResponse{
		Layers:   layers,
		Warnings: warnings,
	}
}

func regionImpactToProto(region *dto.RegionImpact) *tqdpb.Region {
	if region == nil {
		return nil
	}
	return &tqdpb.Region{
		Id:             region.ID,
		Name:           region.Name,
		DisplayName:    region.DisplayName,
		LandUseCode:    region.LandUseCode,
		LandUseName:    region.LandUseName,
		LandUseGroup:   region.LandUseGroup,
		Color:          region.ColorRGB,
		OverlapArea:    region.OverlapAreaSqm,
		OverlapPercent: region.OverlapPct,
		CenterLat:      region.CenterLat,
		CenterLng:      region.CenterLng,
		Geometry:       region.Geometry,
		WarnLevel:      region.WarnLevel,
		CanBuild:       region.CanBuild,
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

func (m *ParcelMapper) ToProtoParcelOverview(data *dto.ParcelLayersResponse) *tqdpb.ParcelOverviewResponse {
	resp := data
	layers := make([]*tqdpb.Layer, 0, len(data.Layers))
	for _, layer := range data.Layers {
		if layer == nil {
			continue
		}
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
			legalDocStr = layer.LegalDoc
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
			LegalStatus:      layer.LegalStatus,
			LegalStatusName:  layer.LegalStatusName,
			LayerAvatar:      layer.LayerAvatar,
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

func (m *ParcelMapper) ParcelInfoToProto(p *qh_domain.QHParcelInfo) *tqdpb.ParcelInfo {
	if p == nil {
		return nil
	}
	mapNum, _ := strconv.ParseUint(p.MapNumber, 10, 32)
	landNum, _ := strconv.ParseUint(p.LandNumber, 10, 32)
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

func (m *ParcelMapper) buildLayersProto(layers []*dto.LayerImpact) []*tqdpb.Layer {
	if len(layers) == 0 {
		return nil
	}
	result := make([]*tqdpb.Layer, 0, len(layers))
	for _, l := range layers {
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
			Id:             r.ID,
			Name:           r.Name,
			DisplayName:    r.DisplayName,
			LandUseCode:    r.LandUseCode,
			LandUseName:    r.LandUseName,
			LandUseGroup:   r.LandUseGroup,
			Color:          r.ColorRGB,
			OverlapArea:    r.OverlapAreaSqm,
			OverlapPercent: r.OverlapPct,
			CenterLat:      r.CenterLat,
			CenterLng:      r.CenterLng,
			Geometry:       r.Geometry,
			WarnLevel:      r.WarnLevel,
			CanBuild:       r.CanBuild,
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

func mapTypeLayerToProto(l *types.LayerDetail) *tqdpb.Layer {
	if l == nil {
		return nil
	}
	pl := &tqdpb.Layer{
		Id:              l.ID,
		Name:            l.Name,
		DisplayName:     l.DisplayName,
		Type:            l.Type,
		EffectiveDate:   l.EffectiveDate,
		TotalArea:       l.TotalArea,
		TotalPercent:    l.TotalPct,
		LegalStatus:     l.LegalStatus,
		LegalStatusName: l.LegalStatusName,
		LayerAvatar:     l.LayerAvatar,
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
				Color:          r.LandUseColor,
				OverlapArea:    r.OverlapAreaSqm,
				OverlapPercent: r.OverlapPct,
				Geometry:       r.Geometry,
				WarnLevel:      r.WarnLevel,
				CanBuild:       r.CanBuild,
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
				Color:          r.LandUseColor,
				OverlapAreaSqm: r.OverlapAreaSqm,
				OverlapPct:     r.OverlapPct,
				Geometry:       r.Geometry,
			})
		}
		pl.Labels = append(pl.Labels, plb)
	}
	return pl
}
