package builder

import (
	_utils "common/utils"
	"tqd/internal/dto"
	"tqd/internal/usecase/resolver/layer_resolver/types"
)

// CandidateBuilder chuyển ParcelLayerRowV2 → types.LayerCandidate
// để scoreAndRank trong engine có thể xử lý.
//
// Nguyên tắc:
//   - LandUseCode/Name/Color luôn lấy từ RegionLandUseXxx (tức qh_land_use)
//   - LandUseColor ưu tiên RegionLandUseColor (lu.color từ SQL unified)
//   - GroupCode = LabelGroupCode (đã resolve ở SQL: COALESCE(group_code, 'unknown'))
//   - KHÔNG có LayerAuthority trên DTO — bỏ field đó khỏi candidate

type CandidateBuilder struct{}

func NewCandidateBuilder() *CandidateBuilder {
	return &CandidateBuilder{}
}

func (b *CandidateBuilder) BuildFromRow(row dto.ParcelLayerRowV2) *types.LayerCandidate {
	// LandUseColor: ưu tiên RegionLandUseColor (lu.color) — đã COALESCE ở SQL
	landUseColor := row.RegionLandUseColor
	if landUseColor == "" {
		landUseColor = row.LabelColor
	}
	if landUseColor == "" {
		landUseColor = row.GroupColor
	}

	return &types.LayerCandidate{
		// Layer identity
		LayerID:          row.LayerID,
		LayerName:        row.LayerName,
		LayerDisplayName: row.LayerDisplayName,
		LayerType:        row.LayerType,
		LayerStatus:      row.LayerStatus,

		// Legal
		LegalStatus:   row.LayerLegalStatus,
		TrustValue:    row.LayerTrustValue,
		EffectiveDate: _utils.ParseStringToTime(row.LayerEffectiveDate),
		ExpiryDate:    _utils.ParseStringToTime(row.LayerExpiryDate),
		// Authority KHÔNG map vào candidate — không có field này trên LayerCandidate

		// Label
		LabelID:          row.LabelID,
		LabelName:        row.LabelName,
		LabelDisplayName: row.LabelDisplayName,
		LabelColor:       row.LabelColor,

		// Land Use — PRIMARY từ qh_land_use (lu)
		LandUseCode:    row.RegionLandUseCode, // lu.code
		LandUseName:    row.RegionLandUseName, // lu.name
		LandUseColor:   landUseColor,          // lu.color
		CanBuild:       row.LabelCanBuild,     // COALESCE(lu.can_build, luc, lg) — từ SQL
		Priority:       row.LabelPriority,     // COALESCE(lu.priority, luc, lg) — từ SQL
		BuildCondition: row.LabelBuildCondition,

		// Group (secondary, display only)
		GroupCode:  row.LabelGroupCode, // = COALESCE(group_code, 'unknown') từ SQL
		GroupName:  row.GroupName,
		GroupColor: row.GroupColor,

		// Overlap
		OverlapAreaSqm: row.OverlapAreaSqm,
		OverlapPercent: row.OverlapPct,
		ParcelAreaSqm:  row.ParcelAreaSqm,

		// Spatial
		CenterLat: row.CenterLat,
		CenterLng: row.CenterLng,
		Geometry:  row.Geometry,
	}
}

func (b *CandidateBuilder) BuildBatch(rows []dto.ParcelLayerRowV2) []*types.LayerCandidate {
	candidates := make([]*types.LayerCandidate, 0, len(rows))
	for _, row := range rows {
		candidates = append(candidates, b.BuildFromRow(row))
	}
	return candidates
}
