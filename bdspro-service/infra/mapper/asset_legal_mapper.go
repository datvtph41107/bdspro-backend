package mapper

import (
	"bdspro/internal/domain"
	_models "common/models"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
)

type AssetLegalMapper struct{}

func NewAssetLegalMapper() *AssetLegalMapper {
	return &AssetLegalMapper{}
}

func (m *AssetLegalMapper) AssetLegalToPb(entity *domain.AssetLegal) *bdspropb.AssetLegalDTO {
	return &bdspropb.AssetLegalDTO{
		Id:                  entity.ID,
		DocumentName:        entity.DocumentName,
		DocumentUrl:         entity.DocumentURL,
		IssuedDate:          _utils.FormatTimeToString(entity.IssuedDate),
		ExpiryDate:          _utils.FormatTimeToString(entity.ExpiryDate),
		Description:         entity.Description,
		AssetId:             entity.AssetID,
		RelatedSplitMergeId: entity.RelatedSplitMergeID,
	}
}

func (m *AssetLegalMapper) PbToListAssetLegal(entities []*bdspropb.LegalItem) []domain.AssetLegal {
	domainLegalItems := make([]domain.AssetLegal, len(entities))
	for i, legalItem := range entities {
		domainLegalItems[i] = m.PbToDomain(legalItem)
	}
	return domainLegalItems
}

func (m *AssetLegalMapper) PbToDomain(entity *bdspropb.LegalItem) domain.AssetLegal {
	return domain.AssetLegal{
		BaseEntity: _models.BaseEntity{
			ID: entity.Id,
		},
		DocumentName: entity.DocumentName,
		DocumentURL:  entity.DocumentUrl,
		DocumentType: entity.DocumentType,
		IssuedDate:   _utils.ParseStringToTime(entity.IssuedDate),
		ExpiryDate:   _utils.ParseStringToTime(entity.ExpiryDate),
		Description:  entity.Description,
		AssetID:      entity.AssetId,
	}
}
