package mapper

import (
	"encoding/json"

	seo_strategy "crm/internal/domain/seo_strategy"
	crmpb "pb/types/crm"
)

type SeoQueryEvidenceMapper struct{}

func NewSeoQueryEvidenceMapper() *SeoQueryEvidenceMapper {
	return &SeoQueryEvidenceMapper{}
}

func (m *SeoQueryEvidenceMapper) ToPb(
	item *seo_strategy.SeoQueryEvidence,
) *crmpb.SeoQueryEvidenceResponse {
	if item == nil {
		return nil
	}

	metadataJSON := "{}"
	if len(item.Metadata) > 0 && json.Valid(item.Metadata) {
		metadataJSON = string(item.Metadata)
	}

	response := &crmpb.SeoQueryEvidenceResponse{
		Id:              item.ID,
		ObjectiveId:     item.ObjectiveID,
		UserNeedId:      item.UserNeedID,
		RawQuery:        item.RawQuery,
		NormalizedQuery: item.NormalizedQuery,
		QueryHash:       item.QueryHash,
		SourceSystem:    item.SourceSystem,
		SourceKind:      string(item.SourceKind),
		ObservationType: string(item.ObservationType),
		SourceRecordKey: item.SourceRecordKey,
		SourceReference: item.SourceReference,
		Locale:          item.Locale,
		CountryCode:     item.CountryCode,
		DeviceCategory:  string(item.DeviceCategory),
		Confidence:      string(item.Confidence),
		Status:          string(item.Status),
		MetadataJson:    metadataJSON,
		CreatedAt:       formatStrategyTime(item.CreatedAt),
		UpdatedAt:       formatStrategyTime(item.UpdatedAt),
	}

	if item.ObservedFrom != nil {
		value := formatStrategyTime(*item.ObservedFrom)
		response.ObservedFrom = &value
	}
	if item.ObservedTo != nil {
		value := formatStrategyTime(*item.ObservedTo)
		response.ObservedTo = &value
	}
	if item.ValidUntil != nil {
		value := formatStrategyTime(*item.ValidUntil)
		response.ValidUntil = &value
	}
	if item.Impressions != nil {
		value := *item.Impressions
		response.Impressions = &value
	}
	if item.Clicks != nil {
		value := *item.Clicks
		response.Clicks = &value
	}
	if item.CTR != nil {
		value := *item.CTR
		response.Ctr = &value
	}
	if item.AvgPosition != nil {
		value := *item.AvgPosition
		response.AvgPosition = &value
	}
	if item.SearchVolume != nil {
		value := *item.SearchVolume
		response.SearchVolume = &value
	}

	return response
}
