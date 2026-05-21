package builder

import (
	_utils "common/utils"
	"tqd/internal/dto"
	"tqd/internal/modules/resolver/layer_resolver/types"
)

type CandidateBuilder struct{}

func NewCandidateBuilder() *CandidateBuilder {
	return &CandidateBuilder{}
}

func (b *CandidateBuilder) BuildFromRow(row dto.ParcelLayerRowV2) *types.LayerCandidate {
	return &types.LayerCandidate{
		LayerID:          row.LayerID,
		LayerName:        row.LayerName,
		LayerDisplayName: row.LayerDisplayName,
		LayerType:        row.LayerType,
		LayerStatus:      row.LayerStatus,

		LegalStatus:   row.LayerLegalStatus,
		TrustValue:    row.LayerTrustValue,
		EffectiveDate: _utils.ParseStringToTime(row.LayerEffectiveDate),
		ExpiryDate:    _utils.ParseStringToTime(row.LayerExpiryDate),
		Authority:     row.LayerAuthority,

		LabelID:          row.LabelID,
		LabelName:        row.LabelName,
		LabelDisplayName: row.LabelDisplayName,
		LabelColor:       row.LabelColor,

		LandUseCode:    row.RegionLandUseCode,
		LandUseName:    row.RegionLandUseName,
		LandUseColor:   row.LabelColor,
		GroupCode:      row.LabelGroupCode,
		GroupName:      row.GroupName,
		GroupColor:     row.GroupColor,
		CanBuild:       row.LabelCanBuild,
		Priority:       row.LabelPriority,
		BuildCondition: row.LabelBuildCondition,

		OverlapAreaSqm: row.OverlapAreaSqm,
		OverlapPercent: row.OverlapPct,
		ParcelAreaSqm:  row.ParcelAreaSqm,

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
